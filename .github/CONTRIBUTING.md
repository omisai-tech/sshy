# Contributing to sshy

Thank you for your interest in contributing to sshy! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Conventions](#commit-conventions)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Feature Requests](#feature-requests)
- [Reporting Bugs](#reporting-bugs)
- [Forking](#forking)

## Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors. By participating, you agree to:

- Be respectful and inclusive
- Focus on constructive feedback
- Accept responsibility for mistakes
- Show empathy towards other contributors
- Help create a positive community

## How to Contribute

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a feature branch** from `master`
4. **Make your changes** following the guidelines below
5. **Write tests** for new functionality
6. **Ensure all tests pass**
7. **Commit your changes** using conventional commits
8. **Push to your fork**
9. **Create a Pull Request**

### Types of Contributions

- **Bug fixes**: Fix issues in existing functionality
- **Features**: Add new functionality that benefits the community
- **Documentation**: Improve documentation, README, examples
- **Tests**: Add or improve test coverage
- **Code quality**: Refactoring, performance improvements

## Development Setup

### Prerequisites

- Go 1.19 or later
- Git
- Make (optional, for convenience scripts)

### Setup Steps

1. **Clone the repository:**
   ```bash
   git clone https://github.com/omisai-tech/sshy.git
   cd sshy
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the project:**
   ```bash
   go build -o sshy
   ```

4. **Run tests:**
   ```bash
   go test ./...
   ```

### Development Workflow

- Use `go build` to compile
- Use `go test ./...` to run all tests
- Use `go fmt` to format code
- Use `go vet` to check for common errors

## Coding Standards

### Go Standards

This project follows standard Go conventions and best practices:

- **Code formatting**: Use `go fmt` and `gofmt`
- **Imports**: Group standard library, third-party, and local imports
- **Naming**: Use Go naming conventions (PascalCase for exported, camelCase for unexported)
- **Error handling**: Handle errors appropriately, don't ignore them
- **Documentation**: Document exported functions, types, and methods with comments
- **Concurrency**: Be careful with goroutines and channels

### Project-Specific Standards

- **Configuration separation**: Keep shared config (`servers.yaml`) separate from local config (`local.yaml`)
- **CLI consistency**: Maintain consistent command structure and help text
- **Error messages**: Provide clear, actionable error messages
- **Backward compatibility**: Don't break existing functionality without good reason

### Code Quality Tools

Before submitting a PR, ensure:

- Code is formatted with `go fmt`
- No linting errors from `go vet`
- All tests pass
- Code is reviewed for security issues

## Commit Conventions

This project uses [Conventional Commits](https://conventionalcommits.org/) to ensure clear and consistent commit messages.

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (formatting, etc.)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Examples

```
feat: add support for SSH key passphrases
fix: handle empty server list in fuzzy finder
docs: update installation instructions
refactor: simplify config loading logic
test: add unit tests for server validation
```

### Scope (Optional)

Scopes help categorize commits:

- `cmd`: CLI commands
- `config`: Configuration handling
- `models`: Data models
- `git`: Git operations (if applicable)

Example: `feat(cmd): add completion command`

## Testing

### Test Coverage

- **Minimum coverage**: Aim for 80%+ test coverage
- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test command interactions
- **Edge cases**: Test error conditions and boundary values

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Guidelines

- **Test file naming**: `*_test.go`
- **Test function naming**: `TestXxx` for unit tests, `TestXxxIntegration` for integration tests
- **Table-driven tests**: Use table-driven tests for multiple test cases
- **Mocking**: Use interfaces and dependency injection for testable code
- **Parallel tests**: Use `t.Parallel()` when tests can run concurrently

## Pull Request Process

1. **Ensure your PR is ready:**
   - All tests pass
   - Code is formatted and linted
   - Commit messages follow conventional commits
   - Documentation is updated if needed

2. **Create a descriptive PR:**
   - Clear title following conventional commit format
   - Detailed description of changes
   - Link to related issues
   - Screenshots for UI changes (if applicable)

3. **PR Review Process:**
   - Maintainers will review your PR
   - Address review comments
   - Once approved, a maintainer will merge

4. **PR Guidelines:**
   - Keep PRs focused on a single feature or fix
   - Update documentation if behavior changes
   - Add tests for new functionality
   - Ensure backward compatibility

## Feature Requests

### When to Submit a Feature Request

- The feature would benefit other users
- It's aligned with the project's goals
- It's not already implemented or planned

### How to Submit

1. Check existing issues for similar requests
2. Create a new issue with the "enhancement" label
3. Provide detailed description including:
   - Use case and benefits
   - Proposed implementation (if you have ideas)
   - Potential impact on existing functionality

### Feature Request Guidelines

- **Community benefit**: Features should help multiple users
- **Scope**: Keep features focused and avoid feature creep
- **Compatibility**: Consider backward compatibility
- **Complexity**: Balance benefit against implementation complexity

## Reporting Bugs

### Bug Report Template

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to reproduce**: Step-by-step instructions
3. **Expected behavior**: What should happen
4. **Actual behavior**: What actually happens
5. **Environment**: OS, Go version, sshy version
6. **Additional context**: Screenshots, logs, configuration

### Bug Priority

- **Critical**: Crashes, data loss, security issues
- **High**: Major functionality broken
- **Medium**: Feature partially broken
- **Low**: Minor issues, edge cases

## Forking

If you want to create a significantly different version of sshy:

1. **Create a fork** on GitHub
2. **Clearly document differences** in your fork's README
3. **Consider contributing upstream** if changes could benefit others
4. **Follow your own contribution guidelines** for your fork

### When to Fork vs Contribute

**Contribute upstream if:**
- Your changes fix bugs or add features that benefit the community
- Your changes align with the project's vision
- You're willing to maintain compatibility

**Fork if:**
- You want to create a substantially different tool
- Your changes conflict with the project's direction
- You need features that wouldn't be accepted upstream

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Conventional Commits](https://conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)

Thank you for contributing to sshy! 🎉</content>
<parameter name="filePath">/var/www/omisai/sshy/CONTRIBUTING.md