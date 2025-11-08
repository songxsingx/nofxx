package exchange

import (
	"fmt"
)

type BinanceSpotClient struct{}

func NewBinanceSpotClient(apiKey, secretKey string) Exchange {
	return &BinanceSpotClient{}
}

func (c *BinanceSpotClient) Type() ExchangeType { return ExchangeBinanceSpot }
func (c *BinanceSpotClient) IsSpot() bool       { return true }
func (c *BinanceSpotClient) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	return []Kline{{Time: 1730700000000, Open: 30000, High: 31000, Low: 29500, Close: 30500, Volume: 100}}, nil
}
func (c *BinanceSpotClient) GetTicker(symbol string) (*Ticker, error) {
	return &Ticker{Symbol: symbol, Last: 30500}, nil
}
func (c *BinanceSpotClient) GetBalance() (map[string]float64, error) {
	return map[string]float64{"USDT": 1000, "BTC": 0.1}, nil
}
func (c *BinanceSpotClient) GetOpenPositions() ([]Position, error) { return []Position{}, nil }
func (c *BinanceSpotClient) PlaceOrder(symbol, side, orderType string, qty, price float64) (*Order, error) {
	// 参数验证
	if symbol == "" {
		return nil, fmt.Errorf("symbol不能为空")
	}
	if qty <= 0 {
		return nil, fmt.Errorf("数量必须大于0, 当前: %.8f", qty)
	}
	if side != "BUY" && side != "SELL" && side != "buy" && side != "sell" {
		return nil, fmt.Errorf("无效的方向: %s (有效值: BUY/SELL)", side)
	}
	if orderType != "market" && orderType != "limit" {
		return nil, fmt.Errorf("无效的订单类型: %s (有效值: market/limit)", orderType)
	}

	// ⚠️ 关键:Binance现货市价单不应设置价格和TimeInForce
	if orderType == "market" {
		if price > 0 {
			return nil, fmt.Errorf("Binance现货市价单不应指定价格")
		}
	} else if orderType == "limit" {
		if price <= 0 {
			return nil, fmt.Errorf("Binance现货限价单必须指定价格")
		}
	}

	// TODO: 实现真实的Binance Spot API调用
	return &Order{ID: fmt.Sprintf("binance-%s", symbol), Status: "FILLED"}, nil
}
func (c *BinanceSpotClient) CancelOrder(orderID string) error                     { return nil }
func (c *BinanceSpotClient) GetPrecision(symbol string) (float64, float64, error) { return 8, 8, nil }
