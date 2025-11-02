use argon2::{
    Argon2, PasswordVerifier,
    password_hash::{self, PasswordHasher, SaltString, rand_core::OsRng},
};

#[derive(std::default::Default)]
pub struct PasswordService {
    argon2: Argon2<'static>,
}

impl PasswordService {
    pub fn new() -> Self {
        Self {
            argon2: Argon2::default(),
        }
    }

    pub fn hash_password(&self, password: &str) -> Result<String, password_hash::Error> {
        let salt = SaltString::generate(&mut OsRng);
        let password_hash = self.argon2.hash_password(password.as_bytes(), &salt)?;
        Ok(password_hash.serialize().to_string())
    }

    pub fn verify_password(&self, password: &str, hash: &str) -> bool {
        if let Ok(parsed) = password_hash::PasswordHash::new(hash) {
            self.argon2
                .verify_password(password.as_bytes(), &parsed)
                .is_ok()
        } else {
            false
        }
    }
}
