package strategy

type StrategyType string

const (
    StrategyAI             StrategyType = "ai"
    StrategyHODLBandProfit StrategyType = "hodl_band_profit"
)

type HODLBandProfitConfig struct {
    Symbol             string  `json:"symbol"`
    BaseAmountUSDT     float64 `json:"base_amount_usdt"`
    ProfitTriggerPct   float64 `json:"profit_trigger_pct"`
    ReinvestRatio      float64 `json:"reinvest_ratio"`
    IntervalHours      int     `json:"interval_hours"`
    TakeProfitPct      float64 `json:"take_profit_pct"`
    StopLossPct        float64 `json:"stop_loss_pct"`
}
