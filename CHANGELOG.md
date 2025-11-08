# Changelog

All notable changes to the NexTrade project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [3.0.0] - 2025-01-15

This is a **major breaking update** that completely transforms NexTrade from a static config-based system to a modern web-based trading platform.

### ✨ Major Features

- **Web-Based UI**: Complete React + TypeScript dashboard with real-time monitoring
- **Multi-Agent Competition**: Multiple AI agents compete in real-time with performance ranking
- **Unified Exchange API**: Single interface supporting Binance, Hyperliquid, and Aster DEX
- **Live/Paper Trading**: Real account trading with full risk controls
- **AI Strategy Evolution**: Agents learn and improve from past decisions
- **Professional Dashboard**: Binance-style interface with equity curves and decision logs

### 🔄 Breaking Changes

- **Config Format**: Migrated from YAML configs to SQLite database
- **Auth System**: New JWT + OTP authentication system
- **API Structure**: Completely redesigned RESTful API
- **Frontend**: Replaced terminal UI with web dashboard

### 🐛 Fixes

- Fixed race conditions in position management
- Improved error handling and logging
- Enhanced security with proper input validation

---

## [2.1.0] - 2024-12-01

### ✨ Features

- **Hyperliquid Support**: Added support for Hyperliquid decentralized exchange
- **Multi-Model AI**: Support for DeepSeek, Qwen, and custom OpenAI-compatible APIs
- **Advanced Risk Controls**: Position limits, drawdown limits, and emergency stop
- **Enhanced Logging**: Detailed decision logs with full Chain of Thought reasoning

### 🐛 Fixes

- Fixed issue with position sizing calculation
- Improved stability under high market volatility
- Fixed bug in stop-loss execution timing

---

## [2.0.0] - 2024-11-15

### ✨ Major Features

- **Multi-Exchange Support**: Unified API for multiple exchanges
- **Configurable AI Prompts**: Customizable trading strategies via prompts
- **Real-time Dashboard**: Terminal-based dashboard with live updates
- **Risk Management**: Advanced position sizing and stop-loss mechanisms

### 🔄 Breaking Changes

- **New Config Format**: Migrated from hardcoded values to YAML config file
- **Improved Architecture**: Modular design for easier maintenance and extension

---

## [1.0.0] - 2024-10-01

### ✨ Initial Release

- **Basic Trading**: Simple long/short trading bot
- **Binance Futures**: Support for Binance Futures exchange
- **Technical Analysis**: Basic TA indicators (MACD, RSI, BB)
- **Risk Controls**: Simple stop-loss and position sizing

---

## How to Use This Changelog

### For Users
- Check the [Unreleased] section for upcoming features
- Review version sections to understand what changed
- Follow migration guides for breaking changes

### For Contributors
When making changes, add them to the [Unreleased] section under appropriate categories:
- **Added** - New features
- **Changed** - Changes to existing functionality
- **Deprecated** - Features that will be removed
- **Removed** - Features that were removed
- **Fixed** - Bug fixes
- **Security** - Security fixes

When releasing a new version, move [Unreleased] items to a new version section with date.

---

## Links

- [Documentation](docs/README.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)
- [GitHub Repository](https://github.com/tinkle-community/nextrade)

---

**Last Updated:** 2025-11-01
