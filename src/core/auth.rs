use actix_web::Error;
use actix_web::HttpMessage;
use actix_web::dev::ServiceRequest;
use actix_web::error::ErrorUnauthorized;
use actix_web_httpauth::extractors::bearer::BearerAuth;
use derive_more::derive::Display;
use uuid::Uuid;

pub struct AuthContext {
    pub user_id: Uuid,
    pub email: String,
    pub role: Role,
}

#[derive(Display)]
pub enum Role {
    #[display("admin")]
    Admin,

    #[display("user")]
    User,
}

pub async fn jwt_auth_validator(
    req: ServiceRequest,
    credentials: BearerAuth,
) -> Result<ServiceRequest, (Error, ServiceRequest)> {
    let token = credentials.token();

    let auth_service = req
        .app_data::<actix_web::web::Data<crate::application::auth::AuthService>>()
        .expect("AuthService not found");

    let Ok(claims) = auth_service.validate_access_token(token) else {
        tracing::error!("Invalid JWT Token: {}", token);
        return Err((actix_web::error::ErrorUnauthorized("Invalid token"), req));
    };

    // Extract auth info from claims
    let Ok(user_id) = Uuid::try_parse(&claims.sub) else {
        tracing::error!("Invalid user ID in JWT Token: {}", claims.sub);
        return Err((ErrorUnauthorized("Invalid token"), req));
    };

    let email = claims.email.clone();

    // TODO: better parsing
    let role = match claims.role.as_str() {
        "admin" => Role::Admin,
        _ => Role::User,
    };

    req.extensions_mut().insert(AuthContext {
        user_id,
        email,
        role,
    });

    Ok(req)
}
