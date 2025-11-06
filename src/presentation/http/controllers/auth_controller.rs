use crate::application::auth::errors::AuthError;
use crate::application::auth::{create_expired_refresh_cookie, create_refresh_cookie};
use crate::presentation::http::dtos::{LoginRequest, RegisterRequest, TokenResponse, UserResponse};
use crate::{application::auth::AuthService, core::request_context::RequestContext};
use actix_web::{HttpRequest, HttpResponse, web};
use tracing::instrument;

#[instrument(name = "register", skip(auth_service))]
pub async fn register(
    ctx: RequestContext,
    auth_service: web::Data<AuthService>,
    payload: web::Json<RegisterRequest>,
) -> Result<HttpResponse, AuthError> {
    let user = auth_service
        .register_user(&ctx, &payload.email, &payload.full_name, &payload.password)
        .await?;

    Ok(HttpResponse::Ok().json(UserResponse::from(user)))
}

#[instrument(name = "login", skip(auth_service))]
pub async fn login(
    ctx: RequestContext,
    auth_service: web::Data<AuthService>,
    payload: web::Json<LoginRequest>,
) -> Result<HttpResponse, AuthError> {
    let auth = auth_service
        .login_user(&ctx, &payload.email, &payload.password)
        .await?;

    let refresh_cookie = create_refresh_cookie(&auth.refresh_token);

    Ok(HttpResponse::Ok()
        .cookie(refresh_cookie)
        .json(TokenResponse::from(auth)))
}

#[instrument(name = "logout")]
pub async fn logout(_ctx: RequestContext) -> HttpResponse {
    let expired_cookie = create_expired_refresh_cookie();
    HttpResponse::Ok().cookie(expired_cookie).finish()
}

#[instrument(name = "refresh", skip(req, auth_service))]
pub async fn refresh(
    ctx: RequestContext,
    req: HttpRequest,
    auth_service: web::Data<AuthService>,
) -> Result<HttpResponse, AuthError> {
    let Some(refresh_cookie) = req.cookie("refresh_token") else {
        return Ok(HttpResponse::Unauthorized().body("Missing refresh token cookie"));
    };

    let refresh_token = refresh_cookie.value();
    let auth = auth_service.refresh_flow(&ctx, refresh_token).await?;

    let refresh_cookie = create_refresh_cookie(&auth.refresh_token);

    Ok(HttpResponse::Ok()
        .cookie(refresh_cookie)
        .json(TokenResponse::from(auth)))
}
