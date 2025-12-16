use crate::{ConfigError, ConfigVersion, Result, regex::*};
use serde::{Deserialize, Serialize};
use std::path::{Path, PathBuf};
use validator::{Validate, ValidationError};

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum ModuleTypeEnum {
    Server,
    Sandbox,
    Interceptor,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum SchemaTypeEnum {
    Openapi,
    Asyncapi,
    Jsonschema,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct ConfigV1 {
    #[serde(default = "default_manifest_version")]
    pub manifest_version: u8,

    #[validate(nested)]
    pub package: PackageConfig,

    #[serde(default)]
    #[validate(nested)]
    pub server: Vec<ModuleConfig>,

    #[serde(default)]
    #[validate(nested)]
    pub sandbox: Vec<ModuleConfig>,

    #[serde(default)]
    #[validate(nested)]
    pub interceptor: Vec<ModuleConfig>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct PackageConfig {
    #[validate(length(min = 3, max = 64))]
    #[validate(regex(path = *NAME_REGEX))]
    pub name: String,

    #[serde(default = "default_package_version")]
    #[validate(regex(path = *SEMVER_REGEX))]
    pub version: String,

    #[serde(default)]
    #[validate(length(max = 500))]
    pub description: String,

    #[serde(default)]
    #[validate(length(max = 5))]
    pub keywords: Vec<String>,

    #[serde(default)]
    #[validate(length(max = 3))]
    pub types: Vec<ModuleTypeEnum>,

    #[serde(default)]
    #[validate(custom(function = "validate_authors"))]
    pub authors: Vec<String>,

    #[serde(default)]
    #[validate(url)]
    pub repository: Option<String>,

    #[serde(default)]
    #[validate(url)]
    pub homepage: Option<String>,

    #[serde(default)]
    pub license: Option<String>,

    #[serde(default)]
    #[validate(custom(function = "validate_path_exists"))]
    pub license_file: Option<PathBuf>,

    #[serde(default)]
    #[validate(custom(function = "validate_path_exists"))]
    pub readme: Option<PathBuf>,

    #[serde(default)]
    #[validate(custom(function = "validate_path_exists"))]
    pub changelog: Option<PathBuf>,

    #[serde(default)]
    pub publish: Vec<String>,

    #[serde(default)]
    pub build: Option<String>,

    #[serde(default = "default_package_abi_version")]
    #[validate(regex(path = *SEMVER_REGEX))]
    pub abi_version: String,

    #[validate(custom(function = "validate_path_exists"))]
    pub wasm: Option<PathBuf>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct ModuleConfig {
    #[validate(length(min = 3, max = 64))]
    #[validate(regex(path = *NAME_REGEX))]
    pub name: String,

    #[validate(custom(function = "validate_path_exists"))]
    pub schema: PathBuf,

    #[serde(default = "default_module_schema_type")]
    pub schema_type: SchemaTypeEnum,
}

fn default_manifest_version() -> u8 {
    1
}

fn default_package_version() -> String {
    "0.0.0".to_string()
}

fn default_package_abi_version() -> String {
    "0.2.0".to_string()
}

fn default_module_schema_type() -> SchemaTypeEnum {
    SchemaTypeEnum::Jsonschema
}

fn validate_authors(authors: &[String]) -> std::result::Result<(), ValidationError> {
    let mut errors = Vec::new();

    for (idx, author) in authors.iter().enumerate() {
        let author = author.trim();
        if !AUTHOR_REGEX.is_match(author) {
            errors.push(format!("Invalid author at index {}: '{}'", idx, author));
        }
    }

    if errors.is_empty() {
        return Ok(());
    }

    let mut err = ValidationError::new("invalid_author");
    err.message = Some(errors.join("; ").into());
    Err(err)
}

fn validate_path_exists(path: &Path) -> std::result::Result<(), ValidationError> {
    if !path.exists() {
        let mut err = ValidationError::new("file_not_found");
        err.message = Some(format!("File not found: {}", path.display()).into());
        return Err(err);
    }

    if !path.is_file() {
        let mut err = ValidationError::new("not_a_file");
        err.message = Some(format!("Path is not a file: {}", path.display()).into());
        return Err(err);
    }

    Ok(())
}

impl ConfigVersion for ConfigV1 {
    fn version(&self) -> u8 {
        self.manifest_version
    }

    fn validate(&self) -> Result<()> {
        Validate::validate(self).map_err(|e| ConfigError::Validation(e.to_string()))?;

        self.validate_manifest_version()?;
        self.validate_has_modules()?;
        self.validate_unique_module_names()?;
        self.validate_license_specification()?;

        Ok(())
    }
}

impl ConfigV1 {
    pub fn resolve_paths(&mut self, config_dir: &std::path::Path) {
        let resolve = |path: &mut PathBuf| {
            if path.is_relative() {
                *path = config_dir.join(&path);
            }
        };

        if let Some(path) = &mut self.package.license_file {
            resolve(path);
        }
        if let Some(path) = &mut self.package.readme {
            resolve(path);
        }
        if let Some(path) = &mut self.package.changelog {
            resolve(path);
        }
        if let Some(path) = &mut self.package.wasm {
            resolve(path);
        }

        for module in self.server.iter_mut() {
            resolve(&mut module.schema);
        }
        for module in self.sandbox.iter_mut() {
            resolve(&mut module.schema);
        }
        for module in self.interceptor.iter_mut() {
            resolve(&mut module.schema);
        }
    }

    fn validate_manifest_version(&self) -> std::result::Result<(), ConfigError> {
        if self.manifest_version != 1 {
            return Err(ConfigError::Validation(
                "Unsupported manifest_version, expected 1".into(),
            ));
        }

        Ok(())
    }

    fn validate_has_modules(&self) -> std::result::Result<(), ConfigError> {
        if self.server.is_empty() && self.sandbox.is_empty() && self.interceptor.is_empty() {
            return Err(ConfigError::Validation(
                "At least one module must be defined".into(),
            ));
        }

        Ok(())
    }

    fn validate_unique_module_names(&self) -> std::result::Result<(), ConfigError> {
        use std::collections::HashSet;

        let mut names = HashSet::new();

        for module in self
            .server
            .iter()
            .chain(&self.sandbox)
            .chain(&self.interceptor)
        {
            if !names.insert(&module.name) {
                return Err(ConfigError::Validation(format!(
                    "Duplicate module name found: '{}'",
                    module.name
                )));
            }
        }

        Ok(())
    }

    fn validate_license_specification(&self) -> std::result::Result<(), ConfigError> {
        if self.package.license.is_none() && self.package.license_file.is_none() {
            return Err(ConfigError::Validation(
                "Either 'license' or 'license_file' must be specified".into(),
            ));
        }

        Ok(())
    }
}
