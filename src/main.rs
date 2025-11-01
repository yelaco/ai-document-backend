mod api;
mod application;
mod config;
mod domain;
mod infrastructure;
mod interfaces;

use actix_web::{App, HttpServer, web};
use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tracing_actix_web::TracingLogger;

use crate::api::http::routes::configure as configure_http;
use crate::application::services::UserService;
use crate::infrastructure::persistence::user::PostgresUserRepository;

fn init_tracing() {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    init_tracing();
    let settings = crate::config::load().expect("Failed to load configuration");

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
                .app_data(actix_web::web::Data::new(settings.clone()))
                .app_data(web::Data::new(UserService::new(user_repo.clone())))
                .configure(configure_http)
                .wrap(TracingLogger::default())
        }
    })
    .bind(("0.0.0.0", settings.port))?
    .run()
    .await
}
