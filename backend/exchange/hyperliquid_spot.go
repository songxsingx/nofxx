package exchange

type HyperliquidSpotClient struct{}

func NewHyperliquidSpotClient() Exchange { return &HyperliquidSpotClient{} }

func (c *HyperliquidSpotClient) Type() ExchangeType { return ExchangeHyperliquidSpot }
func (c *HyperliquidSpotClient) IsSpot() bool { return true }
func (c *HyperliquidSpotClient) GetKlines(symbol, interval string, limit int) ([]Kline, error) { return []Kline{{Time: 1730700000000, Open: 30000, Close: 30500, Volume: 100}}, nil }
func (c *HyperliquidSpotClient) GetTicker(symbol string) (*Ticker, error) { return &Ticker{Symbol: symbol, Last: 30500}, nil }
func (c *HyperliquidSpotClient) GetBalance() (map[string]float64, error) { return map[string]float64{"USDC": 1000}, nil }
func (c *HyperliquidSpotClient) PlaceOrder(symbol, side, orderType string, qty, price float64) (*Order, error) { return &Order{ID: "sim-123", Status: "FILLED"}, nil }
func (c *HyperliquidSpotClient) GetOpenPositions() ([]Position, error) { return []Position{}, nil }
func (c *HyperliquidSpotClient) CancelOrder(orderID string) error { return nil }
func (c *HyperliquidSpotClient) GetPrecision(symbol string) (float64, float64, error) { return 6, 6, nil }
