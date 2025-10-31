# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-10-30

### Fixed
- Fixed README.md Known Issues section to correctly reflect that Device Groups now support Create, Update, and Delete operations
- Corrected grammar and spelling errors in documentation ("doesnt" → "doesn't", "no be" → "cannot be")
- Clarified device name update workaround requirements
- Improved clarity about Custom Profile and Profile update limitations for Assignment Groups and Devices

### Changed
- Updated documentation formatting and consistency across all markdown files
- Applied go fmt to ensure consistent code formatting

## [Unreleased]

### Deprecated

- **Assignment Group Resource**: Added deprecation warnings for `group_type` and `install_type` fields
  - `group_type`: This field is deprecated by the SimpleMDM API and may be ignored for accounts using the New Groups Experience
  - `install_type`: The SimpleMDM API recommends setting install_type per-app using the Assign App endpoint instead of at the group level
  - Both fields remain supported for backward compatibility but their behavior may vary by account type
  - See documentation for migration guidance and alternative approaches

### Added
### Changed
### Removed
### Fixed
### Security
