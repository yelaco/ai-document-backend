use actix_web::Error;
use actix_web::dev::ServiceRequest;
use actix_web_httpauth::extractors::bearer::BearerAuth;
use uuid::Uuid;

pub struct AuthContext {
    pub user_id: Uuid,
    pub role: Role,
}

pub enum Role {
    Admin,
    User,
}

pub async fn jwt_auth_validator(
    req: ServiceRequest,
    credentials: BearerAuth,
) -> Result<ServiceRequest, (Error, ServiceRequest)> {
    let token = credentials.token();

    tracing::info!("JWT Token received: {}", token);

    Ok(req)
}
