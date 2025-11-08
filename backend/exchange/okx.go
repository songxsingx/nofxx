package exchange

import (
	"fmt"
)

type OKXConfig struct {
	ApiKey     string `json:"okx_api_key"`
	SecretKey  string `json:"okx_secret_key"`
	Passphrase string `json:"okx_passphrase"`
	Testnet    bool   `json:"okx_testnet"`
}

type OKXClient struct {
	typ ExchangeType
	cfg *OKXConfig
}

func NewOKXClient(typ ExchangeType, cfg *OKXConfig) (Exchange, error) {
	return &OKXClient{typ: typ, cfg: cfg}, nil
}

func (o *OKXClient) Type() ExchangeType { return o.typ }
func (o *OKXClient) IsSpot() bool       { return o.typ.IsSpot() }
func (o *OKXClient) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	return []Kline{{Time: 1730700000000, Open: 30000, Close: 30500, Volume: 100}}, nil
}
func (o *OKXClient) GetTicker(symbol string) (*Ticker, error) {
	return &Ticker{Symbol: symbol, Last: 30500}, nil
}
func (o *OKXClient) GetBalance() (map[string]float64, error) {
	return map[string]float64{"USDT": 1000}, nil
}
func (o *OKXClient) GetOpenPositions() ([]Position, error) {
	if o.IsSpot() {
		return []Position{}, nil
	}
	return []Position{}, nil
}
func (o *OKXClient) PlaceOrder(symbol, side, orderType string, qty, price float64) (*Order, error) {
	// 参数验证
	if symbol == "" {
		return nil, fmt.Errorf("symbol不能为空")
	}
	if qty <= 0 {
		return nil, fmt.Errorf("数量必须大于0, 当前: %.8f", qty)
	}
	if side != "BUY" && side != "SELL" && side != "LONG" && side != "SHORT" {
		return nil, fmt.Errorf("无效的方向: %s (有效值: BUY/SELL/LONG/SHORT)", side)
	}
	if orderType != "market" && orderType != "limit" {
		return nil, fmt.Errorf("无效的订单类型: %s (有效值: market/limit)", orderType)
	}

	// OKX现货和合约的区分处理
	if o.IsSpot() {
		// 现货下单逻辑
		// ⚠️ 关键：现货市价单不应设置 TimeInForce
		if orderType == "market" {
			// 市价单：不设置价格和TimeInForce
			if price > 0 {
				return nil, fmt.Errorf("OKX现货市价单不应指定价格")
			}
		} else if orderType == "limit" {
			// 限价单：必须设置价格
			if price <= 0 {
				return nil, fmt.Errorf("OKX现货限价单必须指定价格")
			}
		}
	} else {
		// 合约下单逻辑
		// TODO: 实现真实的OKX合约下单逻辑
	}

	// TODO: 实现真实的OKX API调用
	return &Order{ID: fmt.Sprintf("okx-%s", symbol), Status: "FILLED"}, nil
}
func (o *OKXClient) CancelOrder(orderID string) error { return nil }
func (o *OKXClient) GetPrecision(symbol string) (float64, float64, error) {
	return 8, 8, nil
}
