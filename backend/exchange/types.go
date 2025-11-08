package exchange

type ExchangeType string

const (
    ExchangeBinance          ExchangeType = "binance"
    ExchangeHyperliquid      ExchangeType = "hyperliquid"
    ExchangeAster            ExchangeType = "aster"
    ExchangeBinanceSpot      ExchangeType = "binance_spot"
    ExchangeHyperliquidSpot  ExchangeType = "hyperliquid_spot"
    ExchangeAsterSpot        ExchangeType = "aster_spot"
    ExchangeGateioSpot       ExchangeType = "gateio_spot"
    ExchangeGateioFutures    ExchangeType = "gateio_futures"
    ExchangeOKXSpot          ExchangeType = "okx_spot"
    ExchangeOKXFutures       ExchangeType = "okx_futures"
)

type GateioConfig struct {
    ApiKey      string `json:"gateio_api_key"`
    SecretKey   string `json:"gateio_secret_key"`
    Passphrase  string `json:"gateio_passphrase"`
    Testnet     bool   `json:"gateio_testnet"`
    Settle      string `json:"gateio_settle"`
}

type Exchange interface {
    Type() ExchangeType
    IsSpot() bool
    GetKlines(symbol, interval string, limit int) ([]Kline, error)
    GetTicker(symbol string) (*Ticker, error)
    GetBalance() (map[string]float64, error)
    GetOpenPositions() ([]Position, error)
    PlaceOrder(symbol, side, orderType string, qty, price float64) (*Order, error)
    CancelOrder(orderID string) error
    GetPrecision(symbol string) (float64, float64, error)
}

type Kline struct { Time int64; Open, High, Low, Close, Volume float64 }
type Ticker struct { Symbol string; Last float64 }
type Position struct { Symbol, Side string; Size, EntryPrice float64 }
type Order struct { ID, Status string; Filled float64 }

func (t ExchangeType) IsSpot() bool {
    return t == ExchangeBinanceSpot || t == ExchangeHyperliquidSpot || t == ExchangeAsterSpot || t == ExchangeGateioSpot || t == ExchangeOKXSpot
}
