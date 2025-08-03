```markdown
# Contributing to MyVue-Bayes

Thank you for your interest in contributing to MyVue-Bayes! This document provides guidelines for contributing to the project.

## Code of Conduct

Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in Issues
2. Create a new issue with a clear title and description
3. Include code examples and error messages
4. Specify your Go version and operating system

### Suggesting Features

1. Check if the feature has already been suggested
2. Create a new issue with the "enhancement" label
3. Clearly describe the feature and its use case
4. Provide examples of how it would work

### Pull Requests

1. Fork the repository
2. Create a new branch for your feature
3. Write tests for your changes
4. Ensure all tests pass
5. Update documentation as needed
6. Submit a pull request with a clear description

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Add comments for exported functions
- Keep functions focused and small

### Testing

- Write unit tests for new functionality
- Ensure coverage remains above 80%
- Test edge cases and error conditions

## Development Setup

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/myvue-bayes.git
cd myvue-bayes

# Create a branch
git checkout -b feature/your-feature-name

# Make changes and test
go test ./...

# Commit with a descriptive message
git commit -m "Add: description of your changes"

# Push to your fork
git push origin feature/your-feature-name
```

## Questions?

Feel free to open an issue for any questions about contributing.
```