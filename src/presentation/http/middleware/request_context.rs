use actix_web::dev::ServiceResponse;
use actix_web::error::ErrorUnauthorized;
use actix_web::{Error, HttpMessage};
use actix_web::{body::MessageBody, dev::ServiceRequest, middleware::Next};
use uuid::Uuid;

use crate::core::auth::AuthContext;

pub struct RequestContext {
    pub request_id: Uuid,
    pub user_id: Option<Uuid>,
}

pub async fn attach_request_context(
    req: ServiceRequest,
    next: Next<impl MessageBody>,
) -> Result<ServiceResponse<impl MessageBody>, Error> {
    let request_id = uuid::Uuid::new_v4();
    let user_id = req
        .extensions()
        .get::<AuthContext>()
        .map(|auth_context| auth_context.user_id);

    let context = RequestContext {
        request_id,
        user_id,
    };
    req.extensions_mut().insert(context);
    next.call(req).await
}
