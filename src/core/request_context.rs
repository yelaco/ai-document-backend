use std::fmt::Display;

use actix_web::HttpMessage;
use actix_web::{FromRequest, HttpRequest};
use futures::future::{Ready, ready};
use uuid::Uuid;

use crate::core::auth::AuthContext;

#[derive(Clone, Debug)]
pub struct RequestContext {
    pub request_id: Uuid,
    pub user_id: Option<Uuid>,
}

impl Display for RequestContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "RequestContext {{ request_id: {}, user_id: {:?} }}",
            self.request_id, self.user_id
        )
    }
}

impl FromRequest for RequestContext {
    type Error = actix_web::Error;
    type Future = Ready<Result<Self, Self::Error>>;

    fn from_request(req: &HttpRequest, _payload: &mut actix_web::dev::Payload) -> Self::Future {
        if let Some(ctx) = req.extensions().get::<AuthContext>() {
            return ready(Ok(RequestContext {
                // TODO: get request id from middleware
                request_id: uuid::Uuid::new_v4(),
                user_id: Some(ctx.user_id),
            }));
        }

        // fallback value if middleware didn't set it
        ready(Ok(RequestContext {
            request_id: uuid::Uuid::new_v4(),
            user_id: None,
        }))
    }
}
