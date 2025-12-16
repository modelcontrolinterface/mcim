use mcim_config::ConfigError;
use std::io;
use toml;

#[test]
fn test_config_error_display() {
    let validation_error = ConfigError::Validation("Invalid field".to_string());

    assert_eq!(
        format!("{}", validation_error),
        "validation error: Invalid field"
    );

    let io_error = ConfigError::Io(io::Error::new(io::ErrorKind::NotFound, "File not found"));

    assert_eq!(format!("{}", io_error), "IO error: File not found");

    let toml_error: toml::de::Error = toml::from_str::<toml::Value>("invalid toml").unwrap_err();
    let parse_error = ConfigError::Parse(toml_error);

    assert!(format!("{}", parse_error).starts_with("parse error:"));
}
