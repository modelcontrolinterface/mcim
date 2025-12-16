mod common;
use crate::common::create_temp_config_file;
use mcim_config::{
    ConfigV1, ConfigVersion, ModuleConfig, ModuleTypeEnum, PackageConfig, SchemaTypeEnum,
};
use std::path::PathBuf;

fn basic_package_config() -> PackageConfig {
    PackageConfig {
        name: "test-package".to_string(),
        version: "0.1.0".to_string(),
        description: "A test package".to_string(),
        keywords: vec![],
        types: vec![],
        authors: vec![],
        repository: None,
        homepage: None,
        license: Some("MIT".to_string()),
        license_file: None,
        readme: None,
        changelog: None,
        publish: vec![],
        build: None,
        abi_version: "0.1.0".to_string(),
        wasm: None,
    }
}

fn basic_module_config(name: &str, schema_file: &PathBuf) -> ModuleConfig {
    ModuleConfig {
        name: name.to_string(),
        schema: schema_file.clone(),
        schema_type: SchemaTypeEnum::Jsonschema,
    }
}

fn basic_config(schema_file: &PathBuf) -> ConfigV1 {
    ConfigV1 {
        manifest_version: 1,
        package: basic_package_config(),
        server: vec![basic_module_config("my-server", schema_file)],
        sandbox: vec![],
        interceptor: vec![],
    }
}

#[test]
fn test_validate_manifest_version() {
    let schema_file = create_temp_config_file("{}");
    let mut config = basic_config(&schema_file.path().to_path_buf());

    config.manifest_version = 2;

    assert!(config.validate().is_err());

    config.manifest_version = 1;

    assert!(config.validate().is_ok());
}

#[test]
fn test_validate_has_modules() {
    let schema_file = create_temp_config_file("{}");
    let mut config = basic_config(&schema_file.path().to_path_buf());

    config.server = vec![];

    assert!(config.validate().is_err());
}

#[test]
fn test_validate_unique_module_names() {
    let schema_file = create_temp_config_file("{}");
    let schema_path = schema_file.path().to_path_buf();
    let mut config = basic_config(&schema_path);

    config.sandbox = vec![basic_module_config("my-server", &schema_path)];

    assert!(config.validate().is_err());
}

#[test]
fn test_validate_license_specification() {
    let schema_file = create_temp_config_file("{}");
    let mut config = basic_config(&schema_file.path().to_path_buf());

    config.package.license = None;

    assert!(config.validate().is_err());
}

macro_rules! test_validation {
    ($field:ident, $value:expr, $should_pass:expr) => {
        let schema_file = create_temp_config_file("{}");
        let mut config = basic_config(&schema_file.path().to_path_buf());

        config.package.$field = $value;

        if $should_pass {
            assert!(config.validate().is_ok());
        } else {
            assert!(config.validate().is_err());
        }
    };
}

#[test]
fn test_package_config_validation() {
    test_validation!(name, "a".to_string(), false);
    test_validation!(name, "a".repeat(65), false);
    test_validation!(name, "Invalid-Name".to_string(), false);
    test_validation!(version, "1.0".to_string(), false);
    test_validation!(description, "a".repeat(501), false);
    test_validation!(keywords, vec!["a".to_string(); 6], false);
    test_validation!(types, vec![ModuleTypeEnum::Server; 4], false);
    test_validation!(authors, vec!["<invalid>".to_string()], false);
    test_validation!(repository, Some("invalid-url".to_string()), false);
    test_validation!(homepage, Some("invalid-url".to_string()), false);
}

#[test]
fn test_path_exists_validation() {
    let schema_file = create_temp_config_file("{}");
    let mut config = basic_config(&schema_file.path().to_path_buf());

    config.package.license_file = Some(PathBuf::from("non-existent-file"));

    assert!(config.validate().is_err());

    let dir = tempfile::tempdir().unwrap();

    config.package.license_file = Some(dir.path().to_path_buf());

    assert!(config.validate().is_err());
}

#[test]
fn test_defaults() {
    let config: ConfigV1 = toml::from_str(
        r#"
[package]
name = "test-package"
license = "MIT"
[[server]]
name = "my-server"
schema = "schema.json"
"#,
    )
    .unwrap();
    assert_eq!(config.manifest_version, 1);
    assert_eq!(config.package.version, "0.0.0");
    assert_eq!(config.package.abi_version, "0.2.0");
    assert!(matches!(
        config.server[0].schema_type,
        SchemaTypeEnum::Jsonschema
    ));
}
