package exchange

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	gateapi "github.com/gateio/gateapi-go/v6"
)

type GateioClient struct {
	typ        ExchangeType
	cfg        *GateioConfig
	spotClient *gateapi.APIClient
	futClient  *gateapi.APIClient
	ctx        context.Context
}

func NewGateioClient(typ ExchangeType, cfg *GateioConfig) (Exchange, error) {
	if cfg.ApiKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("Gate.io API密钥不能为空")
	}

	// 合约交易必须设置Settle
	if typ == ExchangeGateioFutures {
		if cfg.Settle == "" {
			cfg.Settle = "usdt" // 默认使用USDT结算
			log.Printf("⚠️ Settle未设置，默认使用: usdt")
		}
		log.Printf("💰 Gate.io合约Settle: %s", cfg.Settle)
		cfg.Settle = strings.ToLower(cfg.Settle)
	}

	// 创建配置
	config := gateapi.NewConfiguration()

	// 设置测试网或正式网
	if cfg.Testnet {
		// Gate.io测试网
		if typ == ExchangeGateioSpot {
			// 现货测试网暂不可用，使用正式网
			config.BasePath = "https://api.gateio.ws/api/v4"
			log.Printf("⚠️ Gate.io现货测试网暂不可用，使用正式网环境")
		} else {
			// 合约测试网
			config.BasePath = "https://fx-api-testnet.gateio.ws/api/v4"
			log.Printf("🧪 Gate.io合约使用测试网环境: %s", config.BasePath)
		}
	} else {
		config.BasePath = "https://api.gateio.ws/api/v4" // 正式网
		log.Printf("🔧 Gate.io使用正式网环境")
	}

	// 创建客户端
	client := gateapi.NewAPIClient(config)

	// 创建认证上下文（先清洗密钥，移除空格与BOM）
	cleanKey := sanitizeKey(cfg.ApiKey)
	cleanSecret := sanitizeKey(cfg.SecretKey)
	ctx := context.WithValue(context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    cleanKey,
			Secret: cleanSecret,
		})

	// 输出API密钥信息用于调试（隐藏大部分内容）
	apiKeyPreview := ""
	if len(cleanKey) > 8 {
		apiKeyPreview = cleanKey[:4] + "***" + cleanKey[len(cleanKey)-4:]
	} else {
		apiKeyPreview = "***"
	}
	log.Printf("🔑 Gate.io API Key: %s (长度: %d)", apiKeyPreview, len(cleanKey))

	log.Printf("✅ Gate.io客户端初始化成功 (类型: %s, 测试网: %v)", typ, cfg.Testnet)

	return &GateioClient{
		typ:        typ,
		cfg:        cfg,
		spotClient: client,
		futClient:  client,
		ctx:        ctx,
	}, nil
}

func (g *GateioClient) Type() ExchangeType { return g.typ }
func (g *GateioClient) IsSpot() bool       { return g.typ.IsSpot() }

// sanitizeKey 清洗API密钥，移除首尾空格与BOM
func sanitizeKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\uFEFF")
	// 移除所有空白字符（空格、制表、换行、回车）
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

// GetKlines 获取K线数据
func (g *GateioClient) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	symbol = g.normalizeSymbol(symbol)

	if g.IsSpot() {
		// 现货K线
		intervalStr := g.convertInterval(interval)
		from := time.Now().Add(-time.Duration(limit) * g.intervalDuration(interval)).Unix()
		to := time.Now().Unix()

		opts := &gateapi.ListCandlesticksOpts{
			Interval: optional.NewString(intervalStr),
			From:     optional.NewInt64(from),
			To:       optional.NewInt64(to),
			Limit:    optional.NewInt32(int32(limit)),
		}

		candlesticks, _, err := g.spotClient.SpotApi.ListCandlesticks(g.ctx, symbol, opts)
		if err != nil {
			return nil, fmt.Errorf("获取Gate.io现货K线失败: %w", err)
		}

		klines := make([]Kline, 0, len(candlesticks))
		for _, candle := range candlesticks {
			open, _ := strconv.ParseFloat(candle[5], 64)   // open
			high, _ := strconv.ParseFloat(candle[3], 64)   // high
			low, _ := strconv.ParseFloat(candle[4], 64)    // low
			close, _ := strconv.ParseFloat(candle[2], 64)  // close
			volume, _ := strconv.ParseFloat(candle[6], 64) // volume
			timestamp, _ := strconv.ParseInt(candle[0], 10, 64)

			klines = append(klines, Kline{
				Time:   timestamp * 1000, // 转换为毫秒
				Open:   open,
				High:   high,
				Low:    low,
				Close:  close,
				Volume: volume,
			})
		}
		return klines, nil
	} else {
		// 合约K线
		contract := g.normalizeContract(symbol)
		intervalStr := g.convertInterval(interval)
		from := time.Now().Add(-time.Duration(limit) * g.intervalDuration(interval)).Unix()
		to := time.Now().Unix()

		opts := &gateapi.ListFuturesCandlesticksOpts{
			Interval: optional.NewString(intervalStr),
			From:     optional.NewInt64(from),
			To:       optional.NewInt64(to),
			Limit:    optional.NewInt32(int32(limit)),
		}

		candlesticks, _, err := g.futClient.FuturesApi.ListFuturesCandlesticks(g.ctx, g.cfg.Settle, contract, opts)
		if err != nil {
			return nil, fmt.Errorf("获取Gate.io合约K线失败: %w", err)
		}

		klines := make([]Kline, 0, len(candlesticks))
		for _, candle := range candlesticks {
			open, _ := strconv.ParseFloat(candle.O, 64)
			high, _ := strconv.ParseFloat(candle.H, 64)
			low, _ := strconv.ParseFloat(candle.L, 64)
			close, _ := strconv.ParseFloat(candle.C, 64)
			volume := float64(candle.V)

			klines = append(klines, Kline{
				Time:   int64(candle.T) * 1000, // 转换为毫秒
				Open:   open,
				High:   high,
				Low:    low,
				Close:  close,
				Volume: volume,
			})
		}
		return klines, nil
	}
}

// GetTicker 获取行情数据
func (g *GateioClient) GetTicker(symbol string) (*Ticker, error) {
	symbol = g.normalizeSymbol(symbol)

	if g.IsSpot() {
		// 现货行情
		opts := &gateapi.ListTickersOpts{
			CurrencyPair: optional.NewString(symbol),
		}
		tickers, _, err := g.spotClient.SpotApi.ListTickers(g.ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("获取Gate.io现货行情失败: %w", err)
		}
		if len(tickers) == 0 {
			return nil, fmt.Errorf("未找到币种 %s 的行情", symbol)
		}

		last, _ := strconv.ParseFloat(tickers[0].Last, 64)
		return &Ticker{
			Symbol: symbol,
			Last:   last,
		}, nil
	} else {
		// 合约行情
		contract := g.normalizeContract(symbol)
		opts := &gateapi.ListFuturesTickersOpts{
			Contract: optional.NewString(contract),
		}
		tickers, _, err := g.futClient.FuturesApi.ListFuturesTickers(g.ctx, g.cfg.Settle, opts)
		if err != nil {
			return nil, fmt.Errorf("获取Gate.io合约行情失败: %w", err)
		}
		if len(tickers) == 0 {
			return nil, fmt.Errorf("未找到合约 %s 的行情", contract)
		}

		last, _ := strconv.ParseFloat(tickers[0].Last, 64)
		return &Ticker{
			Symbol: symbol,
			Last:   last,
		}, nil
	}
}

// GetBalance 获取账户余额
func (g *GateioClient) GetBalance() (map[string]float64, error) {
	balances := make(map[string]float64)

	if g.IsSpot() {
		// 现货账户余额
		log.Printf("🔍 正在获取Gate.io现货余额...")
		accounts, resp, err := g.spotClient.SpotApi.ListSpotAccounts(g.ctx, nil)
		if err != nil {
			log.Printf("❌ Gate.io现货余额获取失败: %v", err)
			if resp != nil {
				log.Printf("🔴 HTTP响应码: %d", resp.StatusCode)
			}
			if strings.Contains(err.Error(), "code=-2014") {
				log.Printf("🔎 提示: 请确认使用Gate.io正式网现货密钥，且密钥无空白/不可见字符")
			}
			return nil, fmt.Errorf("获取Gate.io现货余额失败: %w", err)
		}

		for _, account := range accounts {
			available, _ := strconv.ParseFloat(account.Available, 64)
			if available > 0 {
				balances[account.Currency] = available
			}
		}

		log.Printf("📊 Gate.io现货余额: %+v", balances)
	} else {
		// 合约账户余额（返回更完整的信息）
		log.Printf("🔍 正在获取Gate.io合约余额 (settle=%s)...", g.cfg.Settle)
		accounts, resp, err := g.futClient.FuturesApi.ListFuturesAccounts(g.ctx, g.cfg.Settle)
		if err != nil {
			log.Printf("❌ Gate.io合约余额获取失败: %v", err)
			if resp != nil {
				log.Printf("🔴 HTTP响应码: %d", resp.StatusCode)
			}
			if strings.Contains(err.Error(), "code=-2014") {
				log.Printf("🔎 提示: 请确认使用Gate.io合约测试网密钥(来自 fx-testnet.gateio.ws)，且密钥无空白/不可见字符；同时确认Settle=%s", g.cfg.Settle)
			}
			return nil, fmt.Errorf("获取Gate.io合约余额失败: %w", err)
		}

		// 解析合约账户字段
		total, _ := strconv.ParseFloat(accounts.Total, 64)
		available, _ := strconv.ParseFloat(accounts.Available, 64)
		unrealizedPnl, _ := strconv.ParseFloat(accounts.UnrealisedPnl, 64)

		// 返回多个字段（使用特殊键名避免冲突）
		balances[g.cfg.Settle] = available                       // 可用余额
		balances[g.cfg.Settle+"_total"] = total                  // 总余额
		balances[g.cfg.Settle+"_unrealized_pnl"] = unrealizedPnl // 未实现盈亏

		log.Printf("📊 Gate.io合约余额 (%s): Total=%.2f, Available=%.2f, UnrealizedPnL=%.2f",
			g.cfg.Settle, total, available, unrealizedPnl)
	}

	return balances, nil
}

// GetOpenPositions 获取持仓
func (g *GateioClient) GetOpenPositions() ([]Position, error) {
	if g.IsSpot() {
		// 现货没有持仓概念
		return []Position{}, nil
	}

	// 合约持仓
	positions, _, err := g.futClient.FuturesApi.ListPositions(g.ctx, g.cfg.Settle, nil)
	if err != nil {
		return nil, fmt.Errorf("获取Gate.io合约持仓失败: %w", err)
	}

	result := make([]Position, 0)
	for _, pos := range positions {
		size := float64(pos.Size)
		if size == 0 {
			continue // 跳过无持仓
		}

		entryPrice, _ := strconv.ParseFloat(pos.EntryPrice, 64)
		side := "LONG"
		if size < 0 {
			side = "SHORT"
			size = -size
		}

		result = append(result, Position{
			Symbol:     g.denormalizeSymbol(pos.Contract),
			Side:       side,
			Size:       size,
			EntryPrice: entryPrice,
		})
	}

	return result, nil
}

// PlaceOrder 下单
func (g *GateioClient) PlaceOrder(symbol, side, orderType string, qty, price float64) (*Order, error) {
	symbol = g.normalizeSymbol(symbol)

	if g.IsSpot() {
		// 现货下单
		order := gateapi.Order{
			CurrencyPair: symbol,
			Side:         strings.ToLower(side), // buy/sell
			Amount:       fmt.Sprintf("%.8f", qty),
			Type:         orderType, // limit/market
		}

		if orderType == "limit" && price > 0 {
			order.Price = fmt.Sprintf("%.8f", price)
		}

		result, _, err := g.spotClient.SpotApi.CreateOrder(g.ctx, order, nil)
		if err != nil {
			return nil, fmt.Errorf("Gate.io现货下单失败: %w", err)
		}

		return &Order{
			ID:     result.Id,
			Status: result.Status,
		}, nil
	} else {
		// 合约下单
		contract := g.normalizeContract(symbol)
		size := int64(qty)
		if side == "SELL" || side == "SHORT" {
			size = -size
		}

		order := gateapi.FuturesOrder{
			Contract: contract,
			Size:     size,
			Tif:      "ioc", // Immediate or Cancel
		}

		if price > 0 {
			order.Price = fmt.Sprintf("%.8f", price)
		} else {
			order.Price = "0" // 市价单
		}

		result, _, err := g.futClient.FuturesApi.CreateFuturesOrder(g.ctx, g.cfg.Settle, order, nil)
		if err != nil {
			return nil, fmt.Errorf("Gate.io合约下单失败: %w", err)
		}

		return &Order{
			ID:     fmt.Sprintf("%d", result.Id),
			Status: result.Status,
		}, nil
	}
}

// CancelOrder 撤单
func (g *GateioClient) CancelOrder(orderID string) error {
	if g.IsSpot() {
		_, _, err := g.spotClient.SpotApi.CancelOrder(g.ctx, orderID, "", nil)
		if err != nil {
			return fmt.Errorf("Gate.io现货撤单失败: %w", err)
		}
	} else {
		id, _ := strconv.ParseInt(orderID, 10, 64)
		_, _, err := g.futClient.FuturesApi.CancelFuturesOrder(g.ctx, g.cfg.Settle, orderID, nil)
		if err != nil {
			return fmt.Errorf("Gate.io合约撤单失败 (ID: %d): %w", id, err)
		}
	}
	return nil
}

// GetPrecision 获取精度
func (g *GateioClient) GetPrecision(symbol string) (float64, float64, error) {
	symbol = g.normalizeSymbol(symbol)

	if g.IsSpot() {
		// 现货精度
		pairs, _, err := g.spotClient.SpotApi.ListCurrencyPairs(g.ctx)
		if err != nil {
			return 8, 8, fmt.Errorf("获取Gate.io现货精度失败: %w", err)
		}

		for _, pair := range pairs {
			if pair.Id == symbol {
				return float64(pair.AmountPrecision), float64(pair.Precision), nil
			}
		}
		return 8, 8, nil // 默认精度
	} else {
		// 合约精度
		contract := g.normalizeContract(symbol)
		contracts, _, err := g.futClient.FuturesApi.ListFuturesContracts(g.ctx, g.cfg.Settle, nil)
		if err != nil {
			return 8, 8, fmt.Errorf("获取Gate.io合约精度失败: %w", err)
		}

		for _, c := range contracts {
			if c.Name == contract {
				// 合约的数量精度通常是整数，价格精度从价格字段推断
				return 0, 8, nil // 合约数量为整数，价格8位小数
			}
		}
		return 0, 8, nil // 默认精度
	}
}

// 辅助方法：标准化币种符号
func (g *GateioClient) normalizeSymbol(symbol string) string {
	// BTCUSDT -> BTC_USDT (Gate.io现货格式)
	symbol = strings.ToUpper(symbol)
	if !strings.Contains(symbol, "_") {
		if strings.HasSuffix(symbol, "USDT") {
			base := strings.TrimSuffix(symbol, "USDT")
			return base + "_USDT"
		}
	}
	return symbol
}

// 辅助方法：标准化合约符号
func (g *GateioClient) normalizeContract(symbol string) string {
	// BTCUSDT -> BTC_USDT (Gate.io合约格式)
	symbol = strings.ToUpper(symbol)
	if !strings.Contains(symbol, "_") {
		if strings.HasSuffix(symbol, "USDT") {
			base := strings.TrimSuffix(symbol, "USDT")
			return base + "_USDT"
		}
	}
	return symbol
}

// 辅助方法：反标准化符号
func (g *GateioClient) denormalizeSymbol(symbol string) string {
	// BTC_USDT -> BTCUSDT
	return strings.ReplaceAll(symbol, "_", "")
}

// 辅助方法：转换时间间隔
func (g *GateioClient) convertInterval(interval string) string {
	// 1m, 5m, 15m, 1h, 4h, 1d
	switch interval {
	case "1m":
		return "1m"
	case "5m":
		return "5m"
	case "15m":
		return "15m"
	case "1h":
		return "1h"
	case "4h":
		return "4h"
	case "1d":
		return "1d"
	default:
		return "1h"
	}
}

// 辅助方法：获取时间间隔的时长
func (g *GateioClient) intervalDuration(interval string) time.Duration {
	switch interval {
	case "1m":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}
