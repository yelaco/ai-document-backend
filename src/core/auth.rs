use actix_web::HttpMessage;
use actix_web::dev::ServiceRequest;
use actix_web::error::ErrorUnauthorized;
use actix_web_httpauth::extractors::basic::BasicAuth;
use actix_web_httpauth::extractors::bearer::BearerAuth;
use uuid::Uuid;

use crate::application::auth::AuthService;

pub struct AuthContext {
    pub user_id: Uuid,
    pub role: Role,
}

pub enum Role {
    Admin,
    User,
}

pub async fn basic_auth_validator(
    req: ServiceRequest,
    credentials: BasicAuth,
) -> Result<actix_web::dev::ServiceRequest, (actix_web::Error, actix_web::dev::ServiceRequest)> {
    let email = credentials.user_id();
    let password = credentials.password().unwrap_or("");

    let auth_service = req
        .app_data::<AuthService>()
        .expect("AuthService must be set");

    let Ok(Some(user)) = auth_service.verify_basic_credentials(email, password).await else {
        return Err((ErrorUnauthorized("invalid credentials"), req));
    };

    req.extensions_mut().insert(AuthContext {
        user_id: user.id,
        role: Role::User,
    });

    Ok(req)
}

pub async fn jwt_auth_validator(
    req: ServiceRequest,
    credentials: BearerAuth,
) -> Result<actix_web::dev::ServiceRequest, (actix_web::Error, actix_web::dev::ServiceRequest)> {
    let token = credentials.token();

    tracing::info!("JWT Token received: {}", token);

    Ok(req)
}
