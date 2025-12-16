use std::io::Write;
use tempfile::{Builder, NamedTempFile};

#[allow(dead_code)]
pub fn create_temp_config_file(content: &str) -> NamedTempFile {
    let file = Builder::new()
        .prefix("config")
        .suffix(".toml")
        .tempfile()
        .unwrap();
    let mut tmpfile = file.reopen().unwrap();

    tmpfile.write_all(content.as_bytes()).unwrap();

    file
}
