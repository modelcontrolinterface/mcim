use crate::{ConfigError, Result, v1::ConfigV1};
use std::path::Path;

pub fn write_config(config: &ConfigV1, path: &Path) -> Result<()> {
    let content = toml::to_string_pretty(config).map_err(|e| {
        ConfigError::Write(format!("Failed to serialize configuration to TOML: {}", e))
    })?;
    std::fs::write(path, content).map_err(|e| {
        ConfigError::Write(format!(
            "Failed to write configuration to file at {:?}: {}",
            path, e
        ))
    })?;
    Ok(())
}
