package trader

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
)

// SpotTrader 币安现货交易器
type SpotTrader struct {
	client *binance.Client

	// 余额缓存
	cachedBalance     map[string]interface{}
	balanceCacheTime  time.Time
	balanceCacheMutex sync.RWMutex

	// 缓存有效期（15秒）
	cacheDuration time.Duration
}

// NewSpotTrader 创建现货交易器
func NewSpotTrader(apiKey, secretKey string, testnet bool) *SpotTrader {
	client := binance.NewClient(apiKey, secretKey)

	// 如果启用测试网，切换到测试网地址
	if testnet {
		client.BaseURL = "https://testnet.binance.vision"
		log.Printf("🧪 [币安现货] 使用测试网模式: %s", client.BaseURL)
	}

	return &SpotTrader{
		client:        client,
		cacheDuration: 15 * time.Second, // 15秒缓存
	}
}

// GetBalance 获取账户余额（带缓存）
func (t *SpotTrader) GetBalance() (map[string]interface{}, error) {
	// 先检查缓存是否有效
	t.balanceCacheMutex.RLock()
	if t.cachedBalance != nil && time.Since(t.balanceCacheTime) < t.cacheDuration {
		cacheAge := time.Since(t.balanceCacheTime)
		t.balanceCacheMutex.RUnlock()
		log.Printf("✓ 使用缓存的账户余额（缓存时间: %.1f秒前）", cacheAge.Seconds())
		return t.cachedBalance, nil
	}
	t.balanceCacheMutex.RUnlock()

	// 缓存过期或不存在，调用API
	log.Printf("🔄 缓存过期，正在调用币安现货API获取账户余额...")
	account, err := t.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("❌ 币安现货API调用失败: %v", err)
		return nil, fmt.Errorf("获取账户信息失败: %w", err)
	}

	// 计算总余额和可用余额
	totalBalance := 0.0
	availableBalance := 0.0

	for _, balance := range account.Balances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)

		// 只关注USDT余额（现货通常用USDT）
		if balance.Asset == "USDT" {
			availableBalance = free
			totalBalance = free + locked
			log.Printf("✓ 币安现货API返回: USDT余额=%s (可用=%s, 冻结=%s)",
				strconv.FormatFloat(totalBalance, 'f', 2, 64),
				balance.Free,
				balance.Locked)
		}
	}

	result := make(map[string]interface{})
	result["totalWalletBalance"] = totalBalance
	result["availableBalance"] = availableBalance
	result["totalUnrealizedProfit"] = 0.0 // 现货没有未实现盈亏

	// 更新缓存
	t.balanceCacheMutex.Lock()
	t.cachedBalance = result
	t.balanceCacheTime = time.Now()
	t.balanceCacheMutex.Unlock()

	return result, nil
}

// GetPositions 获取所有持仓（现货没有持仓概念，返回资产列表）
func (t *SpotTrader) GetPositions() ([]map[string]interface{}, error) {
	log.Printf("🔄 正在调用币安现货API获取资产信息...")
	account, err := t.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取资产失败: %w", err)
	}

	var result []map[string]interface{}
	for _, balance := range account.Balances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)
		total := free + locked

		// 只显示有余额的资产
		if total > 0 {
			assetMap := make(map[string]interface{})
			assetMap["symbol"] = balance.Asset
			assetMap["free"] = free
			assetMap["locked"] = locked
			assetMap["total"] = total
			result = append(result, assetMap)
		}
	}

	return result, nil
}

// SetMarginMode 现货不支持保证金模式
func (t *SpotTrader) SetMarginMode(symbol string, isCrossMargin bool) error {
	log.Printf("⚠️ 现货交易不支持保证金模式设置，跳过")
	return nil
}

// SetLeverage 现货不支持杠杆
func (t *SpotTrader) SetLeverage(symbol string, leverage int) error {
	log.Printf("⚠️ 现货交易不支持杠杆设置，跳过")
	return nil
}

// PlaceOrder 下单（现货）
func (t *SpotTrader) PlaceOrder(symbol, side, orderType string, qty, price float64) (string, error) {
	log.Printf("📝 准备下现货订单: %s %s %.8f @ %.2f", side, symbol, qty, price)

	var order *binance.CreateOrderResponse
	var err error

	switch orderType {
	case "market":
		// 市价单
		if side == "buy" {
			order, err = t.client.NewCreateOrderService().
				Symbol(symbol).
				Side(binance.SideTypeBuy).
				Type(binance.OrderTypeMarket).
				Quantity(fmt.Sprintf("%.8f", qty)).
				Do(context.Background())
		} else {
			order, err = t.client.NewCreateOrderService().
				Symbol(symbol).
				Side(binance.SideTypeSell).
				Type(binance.OrderTypeMarket).
				Quantity(fmt.Sprintf("%.8f", qty)).
				Do(context.Background())
		}
	case "limit":
		// 限价单
		if side == "buy" {
			order, err = t.client.NewCreateOrderService().
				Symbol(symbol).
				Side(binance.SideTypeBuy).
				Type(binance.OrderTypeLimit).
				TimeInForce(binance.TimeInForceTypeGTC).
				Quantity(fmt.Sprintf("%.8f", qty)).
				Price(fmt.Sprintf("%.2f", price)).
				Do(context.Background())
		} else {
			order, err = t.client.NewCreateOrderService().
				Symbol(symbol).
				Side(binance.SideTypeSell).
				Type(binance.OrderTypeLimit).
				TimeInForce(binance.TimeInForceTypeGTC).
				Quantity(fmt.Sprintf("%.8f", qty)).
				Price(fmt.Sprintf("%.2f", price)).
				Do(context.Background())
		}
	default:
		return "", fmt.Errorf("不支持的订单类型: %s", orderType)
	}

	if err != nil {
		log.Printf("❌ 下单失败: %v", err)
		return "", err
	}

	log.Printf("✅ 订单成功: OrderID=%d, Status=%s", order.OrderID, order.Status)
	return fmt.Sprintf("%d", order.OrderID), nil
}

// CancelOrder 取消订单
func (t *SpotTrader) CancelOrder(symbol, orderID string) error {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return fmt.Errorf("无效的订单ID: %s", orderID)
	}

	_, err = t.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderIDInt).
		Do(context.Background())

	if err != nil {
		log.Printf("❌ 取消订单失败: %v", err)
		return err
	}

	log.Printf("✅ 订单已取消: OrderID=%s", orderID)
	return nil
}

// ClosePosition 现货不支持平仓（需要卖出资产）
func (t *SpotTrader) ClosePosition(symbol, side string) error {
	log.Printf("⚠️ 现货交易不支持平仓操作，请使用卖单")
	return fmt.Errorf("现货交易不支持平仓操作")
}

// GetKlines 获取K线数据
func (t *SpotTrader) GetKlines(symbol, interval string, limit int) ([]map[string]interface{}, error) {
	klines, err := t.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("获取K线失败: %w", err)
	}

	var result []map[string]interface{}
	for _, k := range klines {
		kline := make(map[string]interface{})
		kline["time"] = k.OpenTime
		kline["open"], _ = strconv.ParseFloat(k.Open, 64)
		kline["high"], _ = strconv.ParseFloat(k.High, 64)
		kline["low"], _ = strconv.ParseFloat(k.Low, 64)
		kline["close"], _ = strconv.ParseFloat(k.Close, 64)
		kline["volume"], _ = strconv.ParseFloat(k.Volume, 64)
		result = append(result, kline)
	}

	return result, nil
}

// GetPrice 获取最新价格
func (t *SpotTrader) GetPrice(symbol string) (float64, error) {
	prices, err := t.client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("获取价格失败: %w", err)
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("未找到价格信息")
	}

	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("价格解析失败: %w", err)
	}

	return price, nil
}

// OpenLong 开多仓（现货买入）
func (t *SpotTrader) OpenLong(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	log.Printf("📝 现货买入: %s %.8f", symbol, quantity)

	orderID, err := t.PlaceOrder(symbol, "buy", "market", quantity, 0)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["orderId"] = orderID
	result["symbol"] = symbol
	result["side"] = "buy"
	result["quantity"] = quantity
	return result, nil
}

// OpenShort 开空仓（现货不支持做空）
func (t *SpotTrader) OpenShort(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	return nil, fmt.Errorf("现货交易不支持做空")
}

// CloseLong 平多仓（现货卖出）
func (t *SpotTrader) CloseLong(symbol string, quantity float64) (map[string]interface{}, error) {
	log.Printf("📝 现货卖出: %s %.8f", symbol, quantity)

	// 如果quantity为0，获取当前持有量全部卖出
	if quantity == 0 {
		// 获取当前持仓
		positions, err := t.GetPositions()
		if err != nil {
			return nil, err
		}

		for _, pos := range positions {
			if asset, ok := pos["symbol"].(string); ok && asset+"USDT" == symbol {
				if total, ok := pos["total"].(float64); ok {
					quantity = total
					break
				}
			}
		}

		if quantity == 0 {
			return nil, fmt.Errorf("没有持仓可以平仓")
		}
	}

	orderID, err := t.PlaceOrder(symbol, "sell", "market", quantity, 0)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["orderId"] = orderID
	result["symbol"] = symbol
	result["side"] = "sell"
	result["quantity"] = quantity
	return result, nil
}

// CloseShort 平空仓（现货不支持）
func (t *SpotTrader) CloseShort(symbol string, quantity float64) (map[string]interface{}, error) {
	return nil, fmt.Errorf("现货交易不支持平空仓")
}

// GetMarketPrice 获取市场价格（别名）
func (t *SpotTrader) GetMarketPrice(symbol string) (float64, error) {
	return t.GetPrice(symbol)
}

// SetStopLoss 设置止损单（现货使用限价单模拟）
func (t *SpotTrader) SetStopLoss(symbol string, positionSide string, quantity, stopPrice float64) error {
	log.Printf("⚠️ 现货交易暂不支持自动止损，建议手动设置止损限价单")
	return nil
}

// SetTakeProfit 设置止盈单（现货使用限价单模拟）
func (t *SpotTrader) SetTakeProfit(symbol string, positionSide string, quantity, takeProfitPrice float64) error {
	log.Printf("⚠️ 现货交易暂不支持自动止盈，建议手动设置止盈限价单")
	return nil
}

// CancelAllOrders 取消该币种的所有挂单
func (t *SpotTrader) CancelAllOrders(symbol string) error {
	// 获取所有挂单
	openOrders, err := t.client.NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return fmt.Errorf("获取挂单失败: %w", err)
	}

	if len(openOrders) == 0 {
		log.Printf("✓ %s 没有挂单需要取消", symbol)
		return nil
	}

	// 批量取消
	for _, order := range openOrders {
		err := t.CancelOrder(symbol, fmt.Sprintf("%d", order.OrderID))
		if err != nil {
			log.Printf("⚠️ 取消订单 %d 失败: %v", order.OrderID, err)
		}
	}

	log.Printf("✅ 已取消 %s 的 %d 个挂单", symbol, len(openOrders))
	return nil
}

// FormatQuantity 格式化数量到正确的精度
func (t *SpotTrader) FormatQuantity(symbol string, quantity float64) (string, error) {
	// 获取交易对信息
	exchangeInfo, err := t.client.NewExchangeInfoService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("获取交易对信息失败: %w", err)
	}

	if len(exchangeInfo.Symbols) == 0 {
		return fmt.Sprintf("%.8f", quantity), nil
	}

	// 使用交易对的数量精度
	precision := exchangeInfo.Symbols[0].BaseAssetPrecision
	formatStr := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(formatStr, quantity), nil
}
