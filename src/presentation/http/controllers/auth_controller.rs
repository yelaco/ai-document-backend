use crate::application::auth::services::AuthService;
use crate::presentation::http::dtos::LoginRequest;
use actix_web::{HttpResponse, Responder, post, web};

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
async fn login(auth: web::Data<AuthService>, payload: web::Json<LoginRequest>) -> impl Responder {
    match auth.login_user(&payload.email, &payload.password).await {
        Ok(auth) => HttpResponse::Ok().json(serde_json::json!({"access_token": auth.access_token, "refresh_token": auth.refresh_token})),
        Err(e) => HttpResponse::Unauthorized().body(e),
    }
}
