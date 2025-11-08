# Security Policy / 安全政策

**Languages:** [English](#english) | [中文](#中文)

---

# English

## 🛡️ Security Overview

NexTrade is an AI-powered trading system that handles real funds and API credentials. We take security seriously and appreciate the security community's efforts to responsibly disclose vulnerabilities.

**Critical Areas:**
- 🔑 API key storage and handling
- 💰 Trading execution and fund management
- 🔐 Authentication and authorization
- 🗄️ Database security (SQLite)
- 🌐 Web interface and API endpoints

---

## 📋 Supported Versions

We provide security updates for the following versions:

| Version | Supported          | Notes                |
| ------- | ------------------ | -------------------- |
| 3.x     | ✅ Fully supported | Current stable release |
| 2.x     | ⚠️ Limited support | Security fixes only |
| < 2.0   | ❌ Not supported   | Please upgrade       |

**Recommendation:** Always use the latest stable release (v3.x) for best security.

---

## 🔒 Reporting a Vulnerability

### ⚠️ Please DO NOT Publicly Disclose

If you discover a security vulnerability in NexTrade, please **DO NOT**:
- ❌ Open a public GitHub Issue
- ❌ Discuss it on social media (Twitter, Reddit, etc.)
- ❌ Share it in Telegram/Discord groups
- ❌ Post it on security forums before we've had time to fix it

Public disclosure before a fix is available puts all users at risk.

### ✅ Responsible Disclosure Process

**Step 1: Report Privately**

Send an email to security@nofx.ai with:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Your contact information (optional but preferred)

**Step 2: Confirmation**

We will acknowledge receipt of your report within 48 hours.

**Step 3: Analysis & Fix**

Our security team will:
- Analyze the vulnerability
- Develop a fix
- Test the fix
- Prepare a security advisory

**Step 4: Coordinated Release**

We will:
- Notify you when the fix is ready
- Release the fix and advisory simultaneously
- Credit you in the advisory (if you wish)

**Timeline:** We aim to resolve critical vulnerabilities within 30 days.

---

## 🔐 Deployment Security Best Practices

To keep your NexTrade deployment secure:

### 🔑 API Key Security
- Use exchange API keys with minimal permissions (trading only, no withdrawal)
- Rotate API keys regularly
- Never commit API keys to version control
- Use environment variables or secure vaults for key storage

### 🌐 Network Security
- Run NexTrade behind a firewall
- Use HTTPS for all web interfaces
- Restrict access to trusted IP addresses only
- Keep ports closed except for necessary ones (8080 for API, 3000 for web)

### 🗄️ Data Security
- Regularly backup the SQLite database
- Store backups in encrypted storage
- Use strong passwords for admin access
- Enable two-factor authentication

### ⬆️ Update Management
- Regularly update to the latest stable version
- Monitor security advisories
- Test updates in a staging environment first

### 📊 Monitoring
- Monitor trading activity for anomalies
- Enable detailed logging
- Set up alerts for unusual activity
- Regularly review logs

---

## 📞 Contact

**For security issues ONLY:**
- 📧 **Email:** security@nextrade.ai
- 🐦 **Twitter DM:** [@Web3Tinkle](https://x.com/Web3Tinkle)

**For general questions:**
- See [CONTRIBUTING.md](CONTRIBUTING.md)
- Join [Telegram Community](https://t.me/nofx_dev_community)

---

**Thank you for helping keep NexTrade secure!** 🔒

---

# 中文

## 🛡️ 安全概述

NexTrade 是一个处理真实资金和 API 凭证的 AI 交易系统。我们非常重视安全，并感谢安全社区负责任地披露漏洞的努力。

**关键领域：**
- 🔑 API 密钥存储和处理
- 💰 交易执行和资金管理
- 🔐 身份验证和授权
- 🗄️ 数据库安全（SQLite）
- 🌐 Web 界面和 API 端点

---

## 📋 支持的版本

我们为以下版本提供安全更新：

| 版本 | 支持情况 | 说明 |
| ------- | ------------------ | -------------------- |
| 3.x     | ✅ 完全支持 | 当前稳定版本 |
| 2.x     | ⚠️ 有限支持 | 仅安全修复 |
| < 2.0   | ❌ 不支持 | 请升级 |

**建议：** 始终使用最新的稳定版本 (v3.x) 以获得最佳安全性。

---

## 🔒 漏洞报告

### ⚠️ 请勿公开披露

如果您在 NexTrade 中发现安全漏洞，请**不要**：
- ❌ 在 GitHub 上公开问题
- ❌ 在社交媒体（Twitter、Reddit 等）上讨论
- ❌ 在 Telegram/Discord 群组中分享
- ❌ 在我们有时间修复之前在安全论坛上发布

在修复可用之前公开披露会使所有用户面临风险。

### ✅ 负责任的披露流程

**第 1 步：私下报告**

发送邮件至 security@nofx.ai，包含：
- 漏洞描述
- 重现步骤
- 潜在影响
- 您的联系信息（可选但推荐）

**第 2 步：确认**

我们将在 48 小时内确认收到您的报告。

**第 3 步：分析和修复**

我们的安全团队将：
- 分析漏洞
- 开发修复方案
- 测试修复
- 准备安全公告

**第 4 步：协调发布**

我们将：
- 在修复准备就绪时通知您
- 同时发布修复和公告
- 在公告中致谢您（如果您愿意）

**时间线：** 我们旨在 30 天内解决关键漏洞。

---

## 🔐 部署安全最佳实践

保护您的 NexTrade 部署安全：

### 🔑 API 密钥安全
- 使用权限最小的交易所 API 密钥（仅交易，无提现）
- 定期轮换 API 密钥
- 切勿将 API 密钥提交到版本控制
- 使用环境变量或安全保险库存储密钥

### 🌐 网络安全
- 在防火墙后运行 NexTrade
- 为所有 Web 界面使用 HTTPS
- 仅限受信任的 IP 地址访问
- 除必要端口外保持端口关闭（API 为 8080，Web 为 3000）

### 🗄️ 数据安全
- 定期备份 SQLite 数据库
- 将备份存储在加密存储中
- 为管理员访问使用强密码
- 启用双因素认证

### ⬆️ 更新管理
- 定期更新到最新稳定版本
- 监控安全公告
- 首先在暂存环境中测试更新

### 📊 监控
- 监控交易活动中的异常
- 启用详细日志记录
- 为异常活动设置警报
- 定期审查日志

---

## 📞 联系方式

**仅限安全问题：**
- 📧 **邮箱：** security@nextrade.ai
- 🐦 **Twitter 私信：** [@Web3Tinkle](https://x.com/Web3Tinkle)

**一般问题：**
- 参见 [CONTRIBUTING.md](CONTRIBUTING.md)
- 加入 [Telegram 社区](https://t.me/nofx_dev_community)

---

**感谢您帮助保持 NexTrade 的安全！** 🔒
