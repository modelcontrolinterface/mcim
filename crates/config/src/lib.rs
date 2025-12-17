pub mod error;
pub mod loader;
pub mod regex;
pub mod v1;
pub mod writer;

pub use error::ConfigError;
pub use loader::load_config;
pub use v1::{ConfigV1, ModuleCategoryEnum, ModuleConfig, PackageConfig, SchemaTypeEnum};
pub use writer::write_config;

pub type Result<T> = std::result::Result<T, ConfigError>;

pub trait ConfigVersion {
    fn version(&self) -> u8;
    fn validate(&self) -> Result<()>;
}
