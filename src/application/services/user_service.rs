use std::sync::Arc;

use crate::{domain::User, interfaces::UserRepository};

pub struct UserService {
    repository: Arc<dyn UserRepository>,
}

impl UserService {
    pub fn new(repository: Arc<dyn UserRepository>) -> Self {
        Self { repository }
    }

    pub async fn get_user_by_id(&self, user_id: &str) -> Result<Option<User>, String> {
        self.repository.get_user_by_id(user_id)
    }
}
