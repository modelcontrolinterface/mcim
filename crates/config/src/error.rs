use thiserror::Error;

#[derive(Debug, Error)]
pub enum ConfigError {
    #[error("validation error: {0}")]
    Validation(String),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("parse error: {0}")]
    Parse(#[from] toml::de::Error),

    #[error("write error: {0}")]
    Write(String),
}
