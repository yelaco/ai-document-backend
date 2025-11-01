use crate::api::http::dtos::UserResponse;
use crate::application::services::UserService;
use actix_web::{HttpResponse, Responder, get, web};

#[get("/{id}")]
#[tracing::instrument(name = "Get user by id", skip(service))]
pub async fn get_user(service: web::Data<UserService>, id: web::Path<String>) -> impl Responder {
    match service.get_user_by_id(id.as_str()).await {
        Ok(Some(user)) => HttpResponse::Ok().json(UserResponse::from(user)),
        Ok(None) => HttpResponse::NotFound().body("User not found"),
        Err(err) => {
            tracing::error!("Failed to get user: {}", err);
            HttpResponse::InternalServerError().finish()
        }
    }
}
