# Roll - Virtual Dice Rolling Application

A cross-platform virtual dice rolling application built with Go and Fyne, allowing you to roll dice, save dice sets, and load them for later use.

## Features

- **Dice Rolling**: Use compact notation (e.g., "3d6") to roll multiple dice
- **Visual Results**: See individual die results and total sum
- **Save/Load**: Store dice configurations for quick access
- **Cross-Platform**: Runs on Linux, Mac, Windows, iOS, and Android

## Installation

### From Source

```bash
git clone https://github.com/sfkleach/roll.git
cd roll
go build
```

### Prerequisites

- Go 1.21 or later
- Fyne dependencies for your platform

## Usage

1. Enter dice notation in the input field (e.g., "3d6" for three six-sided dice)
2. Click the "Roll" button to simulate the dice roll
3. View individual die results and the total sum
4. Save frequently used dice sets for quick access

### Dice Notation

The application supports flexible dice notation:

**Basic notation:**
- `3d6` - Roll three six-sided dice
- `2d10` - Roll two ten-sided dice
- `1d20` - Roll one twenty-sided die
- `d20` - Roll one twenty-sided die (count defaults to 1)

**Complex expressions:**
- `2d10 d6` - Roll two ten-sided dice and one six-sided die (space-separated)
- `1d20,7d4` - Roll one twenty-sided die and seven four-sided dice (comma-separated)
- `3d6+2d4` - Roll three six-sided dice and two four-sided dice (plus-separated)
- `d20 2d6 d4` - Mixed notation with implicit counts

## Development

This project uses [Just](https://github.com/casey/just) as a command runner for development tasks.

### Prerequisites

- Go 1.22 or later
- [Just](https://github.com/casey/just) command runner
- Fyne dependencies for your platform:
  - Linux: `sudo apt-get install libgl1-mesa-dev xorg-dev`
  - macOS: Xcode command line tools
  - Windows: TDM-GCC or similar

### Available Commands

```bash
# Show all available commands
just

# Build the application
just build

# Run tests
just test

# Run tests with race detection
just test-race

# Run tests with coverage report
just test-coverage

# Run the application
just run

# Format code
just fmt

# Run linter (requires golangci-lint)
just lint

# Tidy dependencies
just tidy

# Clean build artifacts
just clean

# Install development dependencies
just install-deps

# Cross-compile for all platforms
just build-all

# Run all CI checks
just ci
```

### Building

```bash
# Build with automatic version detection
just build

# Build with specific version
just build "1.4.0"

# Build for current platform using build script
./scripts/build.sh

# Build with custom version using build script
./scripts/build.sh "1.4.0-beta"
```

### Version Management

The application version is automatically managed through git tags:

- **Development builds**: Show version as "dev"
- **Tagged releases**: Version is extracted from git tags (e.g., `v1.4.0` becomes `1.4.0`)
- **GitHub releases**: Version is automatically injected during the release workflow

To create a new release:
1. Create and push a git tag: `git tag v1.4.0 && git push origin v1.4.0`
2. The GitHub Actions workflow will automatically build and release binaries for all platforms

The version is injected at build time using Go's `-ldflags`, so no manual version updates are required in the source code.

### Testing

```bash
just test
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and changes.