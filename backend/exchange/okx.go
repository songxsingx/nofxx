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
	return &Order{ID: fmt.Sprintf("okx-%s", symbol), Status: "FILLED"}, nil
}
func (o *OKXClient) CancelOrder(orderID string) error { return nil }
func (o *OKXClient) GetPrecision(symbol string) (float64, float64, error) {
	return 8, 8, nil
}
