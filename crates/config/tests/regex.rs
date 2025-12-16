use mcim_config::regex::*;

#[test]
fn test_name_regex() {
    assert!(NAME_REGEX.is_match("a"));
    assert!(NAME_REGEX.is_match("a-b"));
    assert!(NAME_REGEX.is_match("a-b-c"));
    assert!(NAME_REGEX.is_match("a-9"));
    assert!(NAME_REGEX.is_match("a_b"));
    assert!(NAME_REGEX.is_match("valid-name"));
    assert!(NAME_REGEX.is_match("another-valid-name_with_numbers-123"));

    assert!(!NAME_REGEX.is_match("A"));
    assert!(!NAME_REGEX.is_match("a-"));
    assert!(!NAME_REGEX.is_match("a_"));
    assert!(!NAME_REGEX.is_match("-invalid-name"));
    assert!(!NAME_REGEX.is_match("_invalid-name"));
    assert!(!NAME_REGEX.is_match("invalid name"));
    assert!(!NAME_REGEX.is_match("a".repeat(65).as_str()));
}

#[test]
fn test_author_regex() {
    assert!(AUTHOR_REGEX.is_match("John Doe"));
    assert!(AUTHOR_REGEX.is_match("John Doe <john.doe@example.com>"));
    assert!(AUTHOR_REGEX.is_match("John Doe <john.doe@example.com> (http://johndoe.com)"));
    assert!(AUTHOR_REGEX.is_match("John Doe (http://johndoe.com)"));

    assert!(!AUTHOR_REGEX.is_match("<john.doe@example.com>"));
    assert!(!AUTHOR_REGEX.is_match("(http://johndoe.com)"));
    assert!(!AUTHOR_REGEX.is_match("John Doe <john.doe@example.com> ()"));
    assert!(!AUTHOR_REGEX.is_match("John Doe <> (http://johndoe.com)"));
}

#[test]
fn test_semver_regex() {
    assert!(SEMVER_REGEX.is_match("1.0.0"));
    assert!(SEMVER_REGEX.is_match("1.0.0-alpha"));
    assert!(SEMVER_REGEX.is_match("1.0.0-alpha.1"));
    assert!(SEMVER_REGEX.is_match("1.0.0-0.3.7"));
    assert!(SEMVER_REGEX.is_match("1.0.0-x.7.z.92"));
    assert!(SEMVER_REGEX.is_match("1.0.0-alpha+build.1"));
    assert!(SEMVER_REGEX.is_match("1.0.0+build.1"));
    assert!(SEMVER_REGEX.is_match("1.0.0-beta+exp.sha.5114f85"));

    assert!(!SEMVER_REGEX.is_match("1"));
    assert!(!SEMVER_REGEX.is_match("1.0"));
    assert!(!SEMVER_REGEX.is_match("1.0.0-"));
    assert!(!SEMVER_REGEX.is_match("1.0.0+"));
    assert!(!SEMVER_REGEX.is_match("1.0.0-alpha."));
    assert!(!SEMVER_REGEX.is_match("1.0.0-alpha..1"));
    assert!(!SEMVER_REGEX.is_match("1.0.0-alpha+build."));
    assert!(!SEMVER_REGEX.is_match("1.0.0-alpha+build..1"));
}
