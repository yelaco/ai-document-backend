use actix_web::middleware::from_fn;
use actix_web::{App, HttpServer, web};
use ai_document_backend::application::auth::AuthService;
use ai_document_backend::application::document::DocumentService;
use ai_document_backend::application::embedding::EmbeddingService;
use ai_document_backend::infrastructure::persistence::document::PostgresDocumentRepository;
use qdrant_client::Qdrant;
use qdrant_client::qdrant::{CreateCollectionBuilder, Distance, VectorParamsBuilder};
use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tokio::sync::Mutex;
use tracing_actix_web::TracingLogger;

use ai_document_backend::application::user::UserService;
use ai_document_backend::infrastructure::persistence::{
    refresh_token::PostgresRefreshTokenRepository, user::PostgresUserRepository,
};
use ai_document_backend::presentation::http::middleware::attach_request_context;
use ai_document_backend::presentation::http::routes::configure as configure_http;

fn init_tracing() {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    init_tracing();
    let settings = ai_document_backend::config::load().expect("Failed to load configuration");

    let pool = PgPoolOptions::new()
        .max_connections(settings.database.max_connections)
        .connect(&settings.database.url)
        .await
        .expect("Failed to create pool");

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("Failed to run migrations");

    let user_repo = Arc::new(PostgresUserRepository::new(pool.clone()));
    let refresh_token_repo = Arc::new(PostgresRefreshTokenRepository::new(pool.clone()));
    let document_repo = Arc::new(PostgresDocumentRepository::new(pool.clone()));

    let qdrant_client = Arc::new({
        let qdrant_url = format!(
            "http://{}:{}",
            settings.qdrant.host, settings.qdrant.grpc_port
        );
        let client = Qdrant::from_url(&qdrant_url)
            .build()
            .expect("Failed to create Qdrant client");

        client
            .health_check()
            .await
            .expect("Failed to connect to Qdrant");

        let _ = client
            .create_collection(
                CreateCollectionBuilder::new("documents")
                    .vectors_config(VectorParamsBuilder::new(768, Distance::Cosine)),
            )
            .await
            .inspect_err(|e| {
                tracing::warn!("Failed to create Qdrant collection: {}", e);
            });

        client
    });

    HttpServer::new({
        let settings = settings.clone();

        move || {
            App::new()
                .wrap(TracingLogger::default())
                .wrap(
                    actix_cors::Cors::default()
                        .allow_any_origin()
                        .allow_any_method()
                        .allow_any_header()
                        .max_age(3600),
                )
                .wrap(from_fn(attach_request_context))
                .app_data(actix_web::web::Data::new(settings.clone()))
                .app_data(web::Data::new(AuthService::new(
                    user_repo.clone(),
                    refresh_token_repo.clone(),
                    settings.clone(),
                )))
                .app_data(web::Data::new(UserService::new(user_repo.clone())))
                .app_data(web::Data::new(DocumentService::new(document_repo.clone())))
                .app_data(web::Data::new(Arc::new(Mutex::new(EmbeddingService::new(
                    qdrant_client.clone(),
                )))))
                .configure(configure_http)
        }
    })
    .bind(("0.0.0.0", settings.app_port))?
    .run()
    .await
}
