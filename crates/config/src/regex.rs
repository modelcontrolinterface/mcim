use lazy_static::lazy_static;
use regex::Regex;

lazy_static! {
    pub static ref NAME_REGEX: Regex = Regex::new(r"^[a-z](?:[a-z0-9_-]{0,62}[a-z0-9])?$").unwrap();
    pub static ref AUTHOR_REGEX: Regex = Regex::new(r"^([^<(]+?)(?:\s*<([^>]+)>)?(?:\s*\(([^)]+)\))?$").unwrap();
    pub static ref SEMVER_REGEX: Regex = Regex::new(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$").unwrap();
}
