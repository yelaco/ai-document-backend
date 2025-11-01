use actix_web::web;

use crate::api::http::handlers::user_handlers;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(web::scope("/api").service(web::scope("/users").service(user_handlers::get_user)));
}
