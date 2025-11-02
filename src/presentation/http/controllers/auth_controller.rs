use crate::presentation::http::dtos::LoginRequest;
use crate::{application::auth::AuthService, core::request_context::RequestContext};
use actix_web::{HttpResponse, Responder, post, web};
use tracing::instrument;

// #[post("/register")]
// async fn register(
//     auth_service: web::Data<AuthService>,
//     payload: web::Json<RegisterRequest>,
// ) -> impl Responder {
//     match auth_service
//         .register_user(&payload.email, &payload.full_name, &payload.password)
//         .await
//     {
//         Ok(user) => HttpResponse::Ok().json(user),
//         Err(e) => HttpResponse::BadRequest().body(e),
//     }
// }

#[post("/login")]
#[instrument(name = "login", skip(ctx, auth_service, payload))]
async fn login(
    ctx: RequestContext,
    auth_service: web::Data<AuthService>,
    payload: web::Json<LoginRequest>,
) -> impl Responder {
    tracing::info!(context = format!("{}", ctx));
    match auth_service.login_user(&payload.email, &payload.password).await {
        Ok(auth) => HttpResponse::Ok().json(serde_json::json!({"access_token": auth.access_token, "refresh_token": auth.refresh_token})),
        Err(e) => HttpResponse::Unauthorized().body(e),
    }
}
