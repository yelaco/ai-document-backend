use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum AppEnv {
    Development,
    Staging,
    Production,
    Test,
}

impl AppEnv {
    pub fn is_prod(&self) -> bool {
        matches!(self, Self::Production)
    }
}

#[derive(Debug, Deserialize, Clone)]
pub struct Settings {
    pub app_name: String,
    pub app_env: AppEnv,
    pub port: u16,
    pub database: DatabaseSettings,
    pub vector_database: VectorDatabaseSettings,
    pub redis: RedisSettings,
    pub ai: AiSettings,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DatabaseSettings {
    pub url: String,
    pub max_connections: u32,
}

#[derive(Debug, Deserialize, Clone)]
pub struct VectorDatabaseSettings {
    pub host: String,
    pub port: u16,
}

#[derive(Debug, Deserialize, Clone)]
pub struct RedisSettings {
    pub host: String,
    pub port: u16,
    pub password: Option<String>,
}

#[derive(Debug, Deserialize, Clone)]
pub struct AiSettings {
    service_type: Option<String>,
    gemini_api_key: Option<String>,
}
