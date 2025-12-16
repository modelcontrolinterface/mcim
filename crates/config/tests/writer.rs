use mcim_config::{ConfigV1, ModuleConfig, PackageConfig, SchemaTypeEnum, writer::write_config, ConfigError};
use std::path::Path;
use tempfile::tempdir;

fn basic_config(schema_file: &Path) -> ConfigV1 {
    ConfigV1 {
        manifest_version: 1,
        package: PackageConfig {
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
        },
        server: vec![ModuleConfig {
            name: "my-server".to_string(),
            schema: schema_file.to_path_buf(),
            schema_type: SchemaTypeEnum::Jsonschema,
        }],
        sandbox: vec![],
        interceptor: vec![],
    }
}

#[test]
fn test_write_config_success() {
    let temp_dir = tempdir().unwrap();
    let schema_path = temp_dir.path().join("schema.json");

    std::fs::write(&schema_path, "{}").unwrap();

    let config = basic_config(&schema_path);
    let output_path = temp_dir.path().join("config.toml");
    let result = write_config(&config, &output_path);

    assert!(result.is_ok());

    let content = std::fs::read_to_string(&output_path).unwrap();
    let loaded_config: ConfigV1 = toml::from_str(&content).unwrap();

    assert_eq!(config.package.name, loaded_config.package.name);
    assert_eq!(config.server[0].name, loaded_config.server[0].name);
}

#[test]
fn test_write_config_permission_denied() {
    let temp_dir = tempdir().unwrap();
    let schema_path = temp_dir.path().join("schema.json");

    std::fs::write(&schema_path, "{}").unwrap();

    let config = basic_config(&schema_path);
    let read_only_dir = tempdir().unwrap();
    let mut perms = std::fs::metadata(read_only_dir.path())
        .unwrap()
        .permissions();

    perms.set_readonly(true);
    std::fs::set_permissions(read_only_dir.path(), perms).unwrap();

    let output_path = read_only_dir.path().join("config.toml");
    let result = write_config(&config, &output_path);

    assert!(result.is_err());
    let err = result.unwrap_err();
    assert!(matches!(err, ConfigError::Write(_)));
    assert!(err.to_string().contains("Failed to write configuration to file"));
}
