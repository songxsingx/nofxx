package exchange

import (
	"fmt"
)

func NewExchange(typ ExchangeType, config map[string]interface{}) (Exchange, error) {
	switch typ {
	case ExchangeBinanceSpot:
		key := config["binance_api_key"].(string)
		secret := config["binance_secret_key"].(string)
		return NewBinanceSpotClient(key, secret), nil
	case ExchangeGateioSpot, ExchangeGateioFutures:
		// 安全的类型断言，带默认值
		apiKey, _ := config["gateio_api_key"].(string)
		secretKey, _ := config["gateio_secret_key"].(string)
		passphrase, _ := config["gateio_passphrase"].(string)
		testnet, _ := config["gateio_testnet"].(bool)
		settle, ok := config["gateio_settle"].(string)
		if !ok || settle == "" {
			settle = "usdt" // 默认USDT结算
		}

		cfg := &GateioConfig{
			ApiKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
			Testnet:    testnet,
			Settle:     settle,
		}
		return NewGateioClient(typ, cfg)
	case ExchangeOKXSpot, ExchangeOKXFutures:
		cfg := &OKXConfig{
			ApiKey:     config["okx_api_key"].(string),
			SecretKey:  config["okx_secret_key"].(string),
			Passphrase: config["okx_passphrase"].(string),
			Testnet:    config["okx_testnet"].(bool),
		}
		return NewOKXClient(typ, cfg)
	case ExchangeHyperliquidSpot:
		// Hyperliquid使用无参构造函数（模拟实现）
		return NewHyperliquidSpotClient(), nil
	case ExchangeAsterSpot:
		// Aster使用apiKey和secretKey参数（模拟实现）
		apiKey := config["aster_api_key"].(string)
		secretKey := config["aster_secret_key"].(string)
		return NewAsterSpotClient(apiKey, secretKey), nil
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", typ)
	}
}
