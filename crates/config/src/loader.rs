use crate::{ConfigVersion, v1::ConfigV1};
use anyhow::{Context, Result};
use std::path::Path;

pub fn load_config(path: &Path) -> Result<ConfigV1> {
    let content = std::fs::read_to_string(path)
        .with_context(|| format!("Failed to read configuration file at {:?}", path))?;
    let mut config: ConfigV1 =
        toml::from_str(&content).with_context(|| "Failed to parse TOML configuration")?;

    if let Some(config_dir) = path.parent() {
        config.resolve_paths(config_dir);
    }

    config
        .validate()
        .map_err(|e| anyhow::anyhow!("Validation error: {}", e))?;

    Ok(config)
}
