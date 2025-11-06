use std::sync::Arc;

use uuid::Uuid;

use crate::{core::request_context::RequestContext, domain::User, interfaces::UserRepository};

pub struct UserService {
    repository: Arc<dyn UserRepository>,
}

impl UserService {
    pub fn new(repository: Arc<dyn UserRepository>) -> Self {
        Self { repository }
    }

    pub async fn get_user_by_id(
        &self,
        _ctx: &RequestContext,
        user_id: Uuid,
    ) -> Result<Option<User>, String> {
        self.repository.get_user_by_id(user_id).await
    }
}
