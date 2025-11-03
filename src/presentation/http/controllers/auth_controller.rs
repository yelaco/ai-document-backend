use crate::application::auth::errors::AuthError;
use crate::presentation::http::dtos::{LoginRequest, LoginResponse, RegisterRequest, UserResponse};
use crate::{application::auth::AuthService, core::request_context::RequestContext};
use actix_web::{HttpResponse, web};
use tracing::instrument;

#[instrument(name = "register", skip(ctx, auth_service, payload))]
pub async fn register(
    ctx: RequestContext,
    auth_service: web::Data<AuthService>,
    payload: web::Json<RegisterRequest>,
) -> Result<HttpResponse, AuthError> {
    let user = auth_service
        .register_user(ctx, &payload.email, &payload.full_name, &payload.password)
        .await?;

    Ok(HttpResponse::Ok().json(UserResponse::from(user)))
}

#[instrument(name = "login", skip(ctx, auth_service, payload))]
pub async fn login(
    ctx: RequestContext,
    auth_service: web::Data<AuthService>,
    payload: web::Json<LoginRequest>,
) -> Result<HttpResponse, AuthError> {
    let auth = auth_service
        .login_user(ctx, &payload.email, &payload.password)
        .await?;

    Ok(HttpResponse::Ok().json(LoginResponse::from(auth)))
}
