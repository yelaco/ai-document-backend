use actix_web::middleware::from_fn;
use actix_web::{App, HttpServer, web};
use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tracing_actix_web::TracingLogger;

use ai_document_backend::application::user::UserService;
use ai_document_backend::infrastructure::persistence::user::PostgresUserRepository;
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

    let user_repo = Arc::new(PostgresUserRepository::new(pool));

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
                .app_data(web::Data::new(
                    ai_document_backend::application::auth::AuthService::new(
                        user_repo.clone(),
                        settings.clone(),
                    ),
                ))
                .app_data(web::Data::new(UserService::new(user_repo.clone())))
                .configure(configure_http)
        }
    })
    .bind(("0.0.0.0", settings.port))?
    .run()
    .await
}
