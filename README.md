# NexTrade - Universal Agentic Trading OS

[![GitHub Stars](https://img.shields.io/github/stars/tinkle-community/nextrade?style=for-the-badge&logo=github&logoColor=white&color=F0B90B&labelColor=0A0A0A)](https://github.com/tinkle-community/nextrade)
[![GitHub Forks](https://img.shields.io/github/forks/tinkle-community/nextrade?style=for-the-badge&logo=github&logoColor=white&color=F0B90B&labelColor=0A0A0A)](https://github.com/tinkle-community/nextrade/network/members)
[![GitHub Contributors](https://img.shields.io/github/contributors/tinkle-community/nextrade?style=for-the-badge&logo=github&logoColor=white&color=F0B90B&labelColor=0A0A0A)](https://github.com/tinkle-community/nextrade/graphs/contributors)
[![Docker Pulls](https://img.shields.io/docker/pulls/tinkle/nextrade?style=for-the-badge&logo=docker&logoColor=white&color=2496ED&labelColor=0A0A0A)](https://hub.docker.com/r/tinkle/nextrade)
[![License](https://img.shields.io/github/license/tinkle-community/nextrade?style=for-the-badge&color=34D058&labelColor=0A0A0A)](LICENSE)

> **NexTrade** is a **universal Agentic Trading OS** - an open-source, community-driven AI trading operating system that closes the loop: **"Multi-Agent Decision → Unified Risk Control → Low-Latency Execution → Live/Paper Account Backtesting"**. 

Built with Go, React, and powered by DeepSeek/Qwen AI, it trades live on Binance, Hyperliquid, and Aster DEX, featuring multi-AI battles and self-learning bots.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/tinkle-community/nextrade.git
cd nextrade

# Backend setup
go mod tidy
./nextrade  # or `go run main.go`

# Frontend setup (in another terminal)
cd web
npm install
npm run dev
```

Open http://localhost:3000 to access the web dashboard.

**Default Login:**
- **Username:** admin@localhost
- **Password:** admin123

⚠️ **Important:** Change the default password immediately after first login!

---

## 🚀 Key Features

### 🎯 Universal Agentic Trading OS
**NexTrade** is a **universal Agentic Trading OS** built on a unified architecture. We've successfully closed the loop in crypto markets: **"Multi-Agent Decision → Unified Risk Control → Low-Latency Execution → Live/Paper Account Backtesting"**, and are now expanding this same technology stack to **stocks, futures, options, forex, and all financial markets**.

### 🧠 Multi-Agent Self-Competition
Multiple AI agents compete in real-time, with the best strategies surviving and evolving - just like natural selection. This creates a self-improving trading system.

### 🔒 100% Self-Custody & Non-Custodial
Your API keys and funds are 100% under your control. We never touch your money. All trading decisions are made by AI agents running on your machine.

### 🌐 Multi-Exchange Support
Currently supports Binance Futures, Hyperliquid, and Aster DEX with more exchanges coming soon.

### 🤖 Multi-AI Model Support
Supports DeepSeek, Qwen, and custom OpenAI-compatible APIs. Easily switch between models to find the best performer.

### 📈 Real-time Dashboard
Professional-grade monitoring interface with real-time updates, equity curves, and full AI decision logs.

---

## 📊 Live Dashboard

### 🎨 Professional Monitoring Interface
- **Binance-Style Dashboard**: Professional dark theme with real-time updates
- **Equity Curves**: Historical account value tracking (USD/percentage toggle)
- **Performance Charts**: Multi-agent ROI comparison with live updates
- **Complete Decision Logs**: Full Chain of Thought (CoT) reasoning for every trade
- **5-Second Data Refresh**: Real-time account, position, and P/L updates

---

## 🧠 Multi-Agent Competition

NexTrade implements a unique **Multi-Agent Self-Competition** mechanism:

1. **Agent Spawning**: Multiple AI agents with different prompts and parameters run in parallel
2. **Real-time Battle**: Agents compete in the same market conditions
3. **Performance Ranking**: Agents are ranked by ROI every 30 minutes
4. **Evolution**: Best-performing agents influence future decision-making
5. **Diversity Preservation**: Maintains strategy diversity to avoid overfitting

This Darwinian approach ensures continuous strategy improvement and adaptation to changing market conditions.

---

## 💱 Exchange Support

### 🟡 Binance Futures
**NexTrade now supports Binance Futures** - the world's largest crypto exchange! To use Binance instead of Hyperliquid:

1. Get your Binance API keys from https://www.binance.com/en/my/settings/api-management
2. In the web UI, go to "Exchanges" configuration
3. Select "Binance Futures" and enter your API keys
4. Enable the exchange and save

**Note:** For security, use "Enable Spot Trading" permission only. Do not enable "Enable Futures" unless you understand the risks.

### 🔵 Hyperliquid
**NexTrade also supports Hyperliquid** - a decentralized perpetual futures exchange. To use Hyperliquid instead of Binance:

1. Get your Hyperliquid wallet private key (export from Hyperliquid web UI)
2. In the web UI, go to "Exchanges" configuration
3. Select "Hyperliquid" and enter your private key
4. Enable the exchange and save

### 🟢 Aster DEX
**NexTrade also supports Aster DEX** - a Binance-compatible decentralized perpetual futures exchange!

Aster DEX is a new decentralized exchange built by former Binance employees, offering:
- **Binance-compatible API**: No need to learn new APIs
- **0 Trading Fees**: 100% fee rebates to traders
- **High Liquidity**: Powered by Binance's matching engine
- **Self-Custody**: Your funds are always in your wallet

To use Aster DEX:
1. Create an account at https://aster.dex (invite code: NEXTRADE)
2. Deposit funds and get your API keys
3. In the web UI, go to "Exchanges" configuration
4. Select "Aster DEX" and enter your API keys
5. Enable the exchange and save

---

## 🤖 AI Model Support

### DeepSeek
- **Model**: deepseek-chat (recommended) / deepseek-coder
- **Website**: https://www.deepseek.com/
- **Pricing**: $0.14M tokens (input) / $0.28M tokens (output)
- **Performance**: Excellent reasoning, good for complex trading strategies

### Qwen (通义千问)
- **Model**: qwen-max (recommended) / qwen-plus / qwen-turbo
- **Website**: https://tongyi.aliyun.com/
- **Pricing**: ¥0.04M tokens (input) / ¥0.12M tokens (output)
- **Performance**: Strong in Chinese market analysis

### Custom OpenAI-Compatible APIs
Supports any OpenAI-compatible API endpoints (like LocalAI, Ollama, etc.)

---

## 🛠️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Server Configuration
API_PORT=8080
ADMIN_MODE=true

# Database (optional, defaults to local SQLite)
# DATABASE_URL=sqlite:///data/nextrade.db

# Authentication (optional, auto-generated if not set)
# JWT_SECRET=your-secure-jwt-secret-here
# OTP_ENCRYPTION_KEY=your-32-byte-encryption-key-here
```

### Default Credentials
- **Username**: admin@localhost
- **Password**: admin123

⚠️ **Important**: Change the default password immediately after first login!

---

## 🏗️ Technical Architecture

NexTrade is built with a modern, modular architecture:

- **Backend:** Go with Gin framework, SQLite database
- **Frontend:** React 18 + TypeScript + Vite + TailwindCSS
- **Multi-Exchange Support:** Binance, Hyperliquid, Aster DEX
- **AI Integration:** DeepSeek, Qwen, and custom OpenAI-compatible APIs
- **State Management:** Zustand for frontend, database-driven for backend
- **Real-time Updates:** SWR with 5-10s polling intervals

---

## 🔮 Roadmap - Universal Market Expansion

NexTrade is on a mission to become the **Universal AI Trading Operating System** for all financial markets.

**Vision:** Same architecture. Same agent framework. All markets.

**Expansion Markets:**
- 📈 **Stock Markets**: US equities, A-shares, Hong Kong stocks
- 📊 **Futures Markets**: Commodity futures, index futures
- 🎯 **Options Trading**: Equity options, crypto options
- 💱 **Forex Markets**: Major currency pairs, cross rates

---

## 🤝 Contributing

We welcome contributions from the community! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 🛡️ Security

Security is our top priority. Please see [SECURITY.md](SECURITY.md) for our security policy and how to report vulnerabilities.

---

## 📜 License

This project is licensed under the AGPL-3.0 License - see the [LICENSE](LICENSE) file for details.
