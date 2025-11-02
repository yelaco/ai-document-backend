use actix_web::web;

use crate::presentation::http::controllers::user_controller;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api").service(web::scope("/users").service(user_controller::get_user)),
    );
}
