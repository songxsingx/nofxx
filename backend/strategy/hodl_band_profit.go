package strategy

import (
	"fmt"
	"log"
	"nextrade/backend/exchange"
	"sync"
	"time"
)

type HODLBandProfitStrategy struct {
	cfg         *HODLBandProfitConfig
	ex          exchange.Exchange
	lastBuy     time.Time
	baseQty     float64
	totalQty    float64
	baseCost    float64
	lastTrigger time.Time
	mu          sync.RWMutex // 并发安全保护
	retryCount  int          // 失败重试次数
	maxRetries  int          // 最大重试次数
}

func NewHODLBandProfitStrategy(cfg *HODLBandProfitConfig, ex exchange.Exchange) *HODLBandProfitStrategy {
	// 参数验证
	if cfg == nil {
		log.Printf("⚠️ HODL策略配置为空，使用默认配置")
		cfg = &HODLBandProfitConfig{
			Symbol:           "BTCUSDT",
			BaseAmountUSDT:   100,
			ProfitTriggerPct: 10,
			ReinvestRatio:    0.5,
			IntervalHours:    1,
			TakeProfitPct:    100,
			StopLossPct:      10,
		}
	}

	// 配置有效性检查
	if cfg.BaseAmountUSDT <= 0 {
		log.Printf("⚠️ 基础金额无效 %.2f，设置为默认值 100", cfg.BaseAmountUSDT)
		cfg.BaseAmountUSDT = 100
	}
	if cfg.ReinvestRatio < 0 || cfg.ReinvestRatio > 1 {
		log.Printf("⚠️ 再投资比例无效 %.2f，设置为默认值 0.5", cfg.ReinvestRatio)
		cfg.ReinvestRatio = 0.5
	}

	return &HODLBandProfitStrategy{
		cfg:        cfg,
		ex:         ex,
		maxRetries: 3, // 最多重试3次
	}
}

func (s *HODLBandProfitStrategy) Execute() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// 带重试的获取行情
	ticker, err := s.getTickerWithRetry()
	if err != nil {
		return fmt.Errorf("获取行情失败: %w", err)
	}

	// 验证价格有效性
	if ticker.Last <= 0 {
		return fmt.Errorf("无效的价格: %.2f", ticker.Last)
	}

	currentValue := s.totalQty * ticker.Last
	baseValue := s.baseQty * ticker.Last
	totalPnL := 0.0
	if s.baseCost > 0 {
		totalPnL = (currentValue - s.baseCost) / s.baseCost * 100
	}

	// 初始买入逻辑
	if s.baseQty == 0 {
		qty := s.cfg.BaseAmountUSDT / ticker.Last
		if qty <= 0 {
			return fmt.Errorf("计算数量失败: 基础金额 %.2f, 价格 %.2f", s.cfg.BaseAmountUSDT, ticker.Last)
		}

		// 验证最小交易数量
		minQty := 0.001 // 假设最小交易数量
		if qty < minQty {
			return fmt.Errorf("交易数量 %.6f 小于最小值 %.6f", qty, minQty)
		}

		_, err = s.placeOrderWithRetry(s.cfg.Symbol, "buy", "market", qty, 0)
		if err != nil {
			return fmt.Errorf("初始买入失败: %w", err)
		}
		s.baseQty = qty
		s.totalQty = qty
		s.baseCost = s.cfg.BaseAmountUSDT
		s.lastBuy = now
		s.lastTrigger = now
		fmt.Printf("[波段盈利定投] 初始买入 %.6f %s @ $%.2f\n", qty, s.cfg.Symbol, ticker.Last)
		return nil
	}

	// 盈利再投资逻辑
	profitPnL := (baseValue - s.baseCost) / s.baseCost * 100
	if profitPnL >= s.cfg.ProfitTriggerPct && now.Sub(s.lastTrigger).Hours() >= float64(s.cfg.IntervalHours) {
		profitUSDT := (baseValue - s.baseCost) * s.cfg.ReinvestRatio
		if profitUSDT > 0 {
			qty := profitUSDT / ticker.Last
			if qty >= 0.001 { // 最小交易数量检查
				_, err = s.placeOrderWithRetry(s.cfg.Symbol, "buy", "market", qty, 0)
				if err != nil {
					return fmt.Errorf("盈利再投资失败: %w", err)
				}
				s.totalQty += qty
				s.lastTrigger = now
				fmt.Printf("[波段盈利定投] 盈利 %.2f%% → 再投 $%.2f (%.6f %s)\n", profitPnL, profitUSDT, qty, s.cfg.Symbol)
			}
		}
	}

	// 止盈止损逻辑
	if s.totalQty > 0 {
		shouldSell := false
		reason := ""

		if s.cfg.TakeProfitPct > 0 && totalPnL >= s.cfg.TakeProfitPct {
			shouldSell = true
			reason = fmt.Sprintf("止盈 (%.2f%%)", s.cfg.TakeProfitPct)
		} else if s.cfg.StopLossPct > 0 && totalPnL <= -s.cfg.StopLossPct {
			shouldSell = true
			reason = fmt.Sprintf("止损 (%.2f%%)", s.cfg.StopLossPct)
		}

		if shouldSell {
			_, err = s.placeOrderWithRetry(s.cfg.Symbol, "sell", "market", s.totalQty, 0)
			if err != nil {
				return fmt.Errorf("%s卖出失败: %w", reason, err)
			}
			fmt.Printf("[波段盈利定投] %s卖出 %.6f %s @ $%.2f (PnL: %.2f%%)\n", reason, s.totalQty, s.cfg.Symbol, ticker.Last, totalPnL)
			s.totalQty = 0
			s.baseQty = 0
			s.baseCost = 0
		}
	}

	return nil
}

// getTickerWithRetry 带重试机制的获取行情
func (s *HODLBandProfitStrategy) getTickerWithRetry() (*exchange.Ticker, error) {
	var ticker *exchange.Ticker
	var err error

	for i := 0; i <= s.maxRetries; i++ {
		ticker, err = s.ex.GetTicker(s.cfg.Symbol)
		if err == nil {
			s.retryCount = 0 // 重置重试计数
			return ticker, nil
		}

		if i < s.maxRetries {
			waitTime := time.Duration(i+1) * time.Second
			log.Printf("⚠️ 获取行情失败，%d秒后重试 (%d/%d): %v", waitTime/time.Second, i+1, s.maxRetries, err)
			time.Sleep(waitTime)
		}
	}

	s.retryCount++
	return nil, fmt.Errorf("重试%d次后仍失败: %w", s.maxRetries, err)
}

// placeOrderWithRetry 带重试机制的下单
func (s *HODLBandProfitStrategy) placeOrderWithRetry(symbol, side, orderType string, qty, price float64) (*exchange.Order, error) {
	var order *exchange.Order
	var err error

	for i := 0; i <= s.maxRetries; i++ {
		order, err = s.ex.PlaceOrder(symbol, side, orderType, qty, price)
		if err == nil {
			return order, nil
		}

		if i < s.maxRetries {
			waitTime := time.Duration(i+1) * 2 * time.Second // 下单重试间隔更长
			log.Printf("⚠️ 下单失败，%d秒后重试 (%d/%d): %v", waitTime/time.Second, i+1, s.maxRetries, err)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("重试%d次后仍失败: %w", s.maxRetries, err)
}

// GetState 获取策略状态（用于监控）
func (s *HODLBandProfitStrategy) GetState() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"base_qty":     s.baseQty,
		"total_qty":    s.totalQty,
		"base_cost":    s.baseCost,
		"last_buy":     s.lastBuy,
		"last_trigger": s.lastTrigger,
		"retry_count":  s.retryCount,
	}
}
