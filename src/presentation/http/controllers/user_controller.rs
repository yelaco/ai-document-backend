use crate::presentation::http::dtos::UserResponse;
use crate::{application::user::UserService, core::request_context::RequestContext};
use actix_web::{HttpResponse, Responder, web};
use uuid::Uuid;

#[tracing::instrument(name = "Get user by id", skip(service))]
pub async fn get_user(
    ctx: RequestContext,
    service: web::Data<UserService>,
    id: web::Path<Uuid>,
) -> impl Responder {
    match service.get_user_by_id(&ctx, id.into_inner()).await {
        Ok(Some(user)) => HttpResponse::Ok().json(UserResponse::from(user)),
        Ok(None) => HttpResponse::NotFound().body("User not found"),
        Err(err) => {
            tracing::error!("Failed to get user: {}", err);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[tracing::instrument(name = "Get current user", skip(service))]
pub async fn get_me(ctx: RequestContext, service: web::Data<UserService>) -> impl Responder {
    let Some(user_id) = ctx.user_id else {
        return HttpResponse::Unauthorized().body("Unauthorized");
    };

    match service.get_user_by_id(&ctx, user_id).await {
        Ok(Some(user)) => HttpResponse::Ok().json(UserResponse::from(user)),
        Ok(None) => HttpResponse::NotFound().body("User not found"),
        Err(err) => {
            tracing::error!("Failed to get user: {}", err);
            HttpResponse::InternalServerError().finish()
        }
    }
}
