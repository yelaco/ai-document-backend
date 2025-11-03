use actix_web::web;
use actix_web_httpauth::middleware::HttpAuthentication;

use crate::presentation::http::controllers::{auth_controller, user_controller};

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api")
            .service(
                web::scope("/auth")
                    .service(web::resource("/login").route(web::post().to(auth_controller::login)))
                    .service(
                        web::resource("/register").route(web::post().to(auth_controller::register)),
                    ),
            )
            .service(
                web::resource("/me")
                    .wrap(HttpAuthentication::bearer(
                        crate::core::auth::jwt_auth_validator,
                    ))
                    .route(web::get().to(user_controller::get_me)),
            )
            .service(
                web::scope("/users")
                    .wrap(HttpAuthentication::bearer(
                        crate::core::auth::jwt_auth_validator,
                    ))
                    .service(user_controller::get_user),
            ),
    );
}
