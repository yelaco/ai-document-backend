mod settings;

use crate::config::settings::Settings;
use config::{Config, ConfigError, Environment, File};

pub fn load() -> Result<Settings, ConfigError> {
    let app_env = std::env::var("APP_ENV").unwrap_or("development".into());

    let builder = Config::builder()
        .add_source(File::with_name("config/settings/default"))
        .add_source(File::with_name(&format!("config/settings/{app_env}")).required(false))
        .add_source(Environment::default().separator("__"))
        .set_override("app_env", app_env)
        .unwrap();

    builder.build()?.try_deserialize()
}
