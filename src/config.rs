use config::{ConfigError, Environment};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub address: String,
    pub log_level: String,

    pub key_path: Option<String>,
    pub cert_path: Option<String>,
}

impl Config {
    pub fn from_env() -> Result<Self, ConfigError> {
        let s = config::Config::builder()
            .set_default("address", "0.0.0.0:7687")?
            .set_default("log_level", "info")?
            .add_source(Environment::with_prefix("MCI").separator("__"))
            .build()?;

        s.try_deserialize()
    }
}
