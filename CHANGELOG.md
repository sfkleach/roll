# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced dice notation parsing with support for complex expressions
  - Single die notation: `d20` (implicit count of 1)
  - Multiple dice groups: `2d10 d6`, `1d20,7d4`, `3d6+2d4`
  - Mixed separators: space, comma, and plus signs
- Comprehensive test coverage for all dice notation formats
- Fyne-based GUI framework setup for cross-platform support
- Full test coverage for dice functionality
- Just command runner with development workflow
- GitHub Actions CI/CD workflows
- Cross-platform build support (Linux, macOS, Windows)
- golangci-lint configuration for code quality

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - 2025-08-24

### Added
- Project initialization
- Basic project structure and documentation
- MIT License
- Go module setup with Fyne dependency
