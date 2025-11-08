# 🤝 Contributing to NexTrade

Thank you for your interest in contributing to NexTrade! This document provides guidelines and workflows for contributing to the project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

## 🤝 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- npm 9+
- TA-Lib (for technical analysis indicators)

### Installation

```bash
# Clone the repository
git clone https://github.com/tinkle-community/nextrade.git
cd nextrade

# Backend setup
go mod tidy

# Frontend setup
cd web
npm install
```

### Running Development Servers

```bash
# Start backend (from project root)
go run main.go

# Start frontend (from web directory)
cd web
npm run dev
```

## 🔧 Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📥 Pull Request Process

### Before Submitting

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you

### PR Template

```
## Description

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context.

Fixes # (issue)

## Type of Change

Please delete options that are not relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

## How Has This Been Tested?

Please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce.

## Checklist:

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules
```

## 💻 Coding Standards

### Go Backend

NexTrade/
├── cmd/               # Main applications
├── internal/          # Private code
│   ├── exchange/      # Exchange adapters
│   ├── trader/        # Trading logic
│   ├── ai/           # AI integrations
│   └── api/          # API handlers
├── pkg/              # Public libraries
├── web/              # Frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   └── utils/
│   └── public/
└── docs/             # Documentation

---

## 📋 Commit Message Guidelines

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Examples

```
feat(exchange): add OKX futures API integration

- Implement order placement and cancellation
- Add balance and position retrieval
- Support leverage configuration

Closes #123
```

```
fix(trader): prevent duplicate position opening

The trader was opening multiple positions in the same direction
for the same symbol. Added check to prevent this behavior.

Fixes #456
```

```
docs: update Docker deployment guide

- Add troubleshooting section
- Update environment variables
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

## 🧪 Testing

### Go Backend Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Frontend Testing

```bash
# Run frontend tests
cd web
npm test

# Run tests with coverage
npm test -- --coverage
```

**Best Practices:**
- Use TypeScript strict mode
- Define interfaces for all data structures
- Avoid `any` type
- Use functional components with hooks
- Follow React best practices
- Run `npm run lint` before committing

## 📚 Documentation

- Keep documentation up to date with code changes
- Use clear, concise language
- Include examples where appropriate
- Follow the existing documentation style

## 👥 Community

### Communication

- **GitHub Issues**: For bug reports and feature requests
- **Pull Requests**: For code contributions
- **Telegram**: [NexTrade Developer Community](https://t.me/nofx_dev_community)

### Getting Help

If you need help, please:
1. Check existing issues and documentation
2. Ask in the [Telegram community](https://t.me/nofx_dev_community)
3. Create a new issue if needed

Your contributions make NexTrade better for everyone. We appreciate your time and effort!
