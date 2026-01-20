# Contributing to RSS Aggregator

Thank you for your interest in contributing to RSS Aggregator! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.21.3 or later
- PostgreSQL database
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- goose (`go install github.com/presslabs/goose/cmd/goose@latest`)

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/reinhardbuyabo/RSS-Aggregator.git
   cd RSS-Aggregator
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/reinhardbuyabo/RSS-Aggregator.git
   ```

4. Create a feature branch:
   ```bash
   git checkout -b my-feature-branch
   ```

5. Set up the database and run the application (see README.md)

## Development Workflow

### 1. Pick an Issue

- Look through [open issues](https://github.com/reinhardbuyabo/RSS-Aggregator/issues) for tasks
- Comment on the issue to claim it and avoid duplicate work
- For major changes, open a new issue first to discuss

### 2. Make Changes

- Follow Go coding conventions and style guides
- Write clear, commented code
- Keep changes focused and minimal

### 3. Testing

- Run existing tests to ensure nothing is broken:
  ```bash
  go test ./...
  ```
- Add tests for new functionality
- Verify code compiles:
  ```bash
  go build ./...
  ```

### 4. SQL Changes

If you modify SQL files in `sql/queries/` or `sql/schema/`:

1. Regenerate Go code:
   ```bash
   sqlc generate
   ```

2. Review generated code for correctness

### 5. Committing

- Use clear, descriptive commit messages
- Reference issue numbers where applicable
- Keep commits atomic (one logical change per commit)

### 6. Submitting Changes

1. Push your branch to your fork:
   ```bash
   git push origin my-feature-branch
   ```

2. Create a Pull Request against the `main` branch
3. Fill out the PR template with all requested information
4. Address any feedback from reviewers

## Code Standards

### Go Style

- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Keep functions small and focused
- Add Go doc comments for exported types and functions

### Error Handling

- Return errors with context
- Handle errors at the appropriate level
- Use the existing error patterns in the codebase

### Database

- Use sqlc for all database operations
- Write clear, efficient SQL queries
- Include proper indexes for performance
- Handle edge cases like duplicates

### Security

- Never commit secrets, API keys, or credentials
- Use environment variables for sensitive data
- Validate all user inputs
- Follow security best practices for web APIs

## Code Review Process

All contributions require review before merging. During review:

- Address all comments thoughtfully
- Ask for clarification if needed
- Keep discussions focused and constructive
- Be patient with feedback

## Community

- Be respectful and inclusive
- Help others when possible
- Follow the project's code of conduct
