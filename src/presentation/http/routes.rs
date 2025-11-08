use actix_web::{middleware, web};
use actix_web_httpauth::middleware::HttpAuthentication;

use crate::presentation::http::controllers::{
    auth_controller, document_controller, user_controller,
};

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api")
            .wrap(middleware::NormalizePath::new(
                middleware::TrailingSlash::Trim,
            ))
            // Auth routes
            .service(
                web::scope("/auth")
                    .service(
                        web::resource("/register").route(web::post().to(auth_controller::register)),
                    )
                    .service(web::resource("/login").route(web::post().to(auth_controller::login)))
                    .service(
                        web::resource("/logout")
                            .wrap(HttpAuthentication::bearer(
                                crate::core::auth::jwt_auth_validator,
                            ))
                            .route(web::post().to(auth_controller::logout)),
                    )
                    .service(
                        web::resource("/refresh")
                            .wrap(HttpAuthentication::bearer(
                                crate::core::auth::jwt_auth_validator,
                            ))
                            .route(web::post().to(auth_controller::refresh)),
                    ),
            )
            // Me routes
            .service(
                web::resource("/me")
                    .wrap(HttpAuthentication::bearer(
                        crate::core::auth::jwt_auth_validator,
                    ))
                    .route(web::get().to(user_controller::get_me)),
            )
            // User routes
            .service(
                web::scope("/users")
                    .wrap(HttpAuthentication::bearer(
                        crate::core::auth::jwt_auth_validator,
                    ))
                    .service(
                        web::resource("/{id}")
                            .route(web::get())
                            .to(user_controller::get_user),
                    ),
            )
            // Document routes
            .service(
                web::scope("/documents")
                    .wrap(HttpAuthentication::bearer(
                        crate::core::auth::jwt_auth_validator,
                    ))
                    .service(
                        web::resource("")
                            .route(web::get().to(document_controller::get_paginated_documents))
                            .route(web::post().to(document_controller::upload)),
                    )
                    .service(
                        web::resource("/{id}")
                            .route(web::get().to(document_controller::get_document))
                            .route(web::delete().to(document_controller::delete_document)),
                    ),
            ),
    );
}
