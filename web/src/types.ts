export interface SystemStatus {
  trader_id: string;
  trader_name: string;
  ai_model: string;
  is_running: boolean;
  start_time: string;
  runtime_minutes: number;
  call_count: number;
  initial_balance: number;
  scan_interval: string;
  stop_until: string;
  last_reset_time: string;
  ai_provider: string;
}

export interface AccountInfo {
  total_equity: number;
  wallet_balance: number;
  unrealized_profit: number;
  available_balance: number;
  total_pnl: number;
  total_pnl_pct: number;
  total_unrealized_pnl: number;
  initial_balance: number;
  daily_pnl: number;
  position_count: number;
  margin_used: number;
  margin_used_pct: number;
}

export interface Position {
  symbol: string;
  side: string;
  entry_price: number;
  mark_price: number;
  quantity: number;
  leverage: number;
  unrealized_pnl: number;
  unrealized_pnl_pct: number;
  liquidation_price: number;
  margin_used: number;
}

export interface DecisionAction {
  action: string;
  symbol: string;
  quantity: number;
  leverage: number;
  price: number;
  order_id: number;
  timestamp: string;
  success: boolean;
  error?: string;
}

export interface AccountSnapshot {
  total_balance: number;
  available_balance: number;
  total_unrealized_profit: number;
  position_count: number;
  margin_used_pct: number;
}

export interface DecisionRecord {
  timestamp: string;
  cycle_number: number;
  input_prompt: string;
  cot_trace: string;
  decision_json: string;
  account_state: AccountSnapshot;
  positions: any[];
  candidate_coins: string[];
  decisions: DecisionAction[];
  execution_log: string[];
  success: boolean;
  error_message?: string;
}

export interface Statistics {
  total_cycles: number;
  successful_cycles: number;
  failed_cycles: number;
  total_open_positions: number;
  total_close_positions: number;
}

// AI Trading相关类型
export interface TraderInfo {
  trader_id: string;
  trader_name: string;
  ai_model: string;
  exchange_id?: string;
  is_running?: boolean;
  custom_prompt?: string;
}

export interface AIModel {
  id: string;
  name: string;
  provider: string;
  enabled: boolean;
  apiKey?: string;
  customApiUrl?: string;
  customModelName?: string;
}

export interface Exchange {
  id: string;
  name: string;
  type: 'cex' | 'dex';
  enabled: boolean;
  apiKey?: string;
  secretKey?: string;
  testnet?: boolean;
  // Hyperliquid 特定字段
  hyperliquidWalletAddr?: string;
  // Aster 特定字段
  asterUser?: string;
  asterSigner?: string;
  asterPrivateKey?: string;
  // Gate.io 特定字段
  gateioPassphrase?: string;
  // OKX 特定字段
  okxPassphrase?: string;
}

export interface CreateTraderRequest {
  name: string;
  ai_model_id: string;
  exchange_id: string;
  initial_balance: number;
  scan_interval_minutes?: number;
  btc_eth_leverage?: number;
  altcoin_leverage?: number;
  trading_symbols?: string;
  custom_prompt?: string;
  override_base_prompt?: boolean;
  system_prompt_template?: string;
  is_cross_margin?: boolean;
  use_coin_pool?: boolean;
  use_oi_top?: boolean;
  // HODL波段盈利定投策略配置
  strategy?: string; // "ai" or "hodl_band_profit"
  strategy_config?: HODLBandProfitConfig;
  // 现货交易配置
  spot_order_type?: 'market' | 'limit';
  spot_position_size_pct?: number;
  spot_take_profit_pct?: number;
  spot_stop_loss_pct?: number;
}

// HODL波段盈利定投策略配置
export interface HODLBandProfitConfig {
  symbol: string;              // 目标币种，如"BTCUSDT"
  base_amount_usdt: number;    // 初始投入金额（USDT）
  profit_trigger_pct: number;  // 盈利触发百分比（如10表示10%）
  reinvest_ratio: number;      // 再投资比例（如0.5表示盈利的50%）
  interval_hours: number;      // 检查间隔（小时）
  take_profit_pct: number;     // 止盈百分比（如100表示100%）
  stop_loss_pct: number;       // 止损百分比（如10表示-10%）
}

export interface UpdateModelConfigRequest {
  models: {
    [key: string]: {
      enabled: boolean;
      api_key: string;
      custom_api_url?: string;
      custom_model_name?: string;
    };
  };
}

export interface UpdateExchangeConfigRequest {
  exchanges: {
    [key: string]: {
      enabled: boolean;
      api_key: string;
      secret_key: string;
      testnet?: boolean;
      // Hyperliquid 特定字段
      hyperliquid_wallet_addr?: string;
      // Aster 特定字段
      aster_user?: string;
      aster_signer?: string;
      aster_private_key?: string;
      // Gate.io 特定字段
      gateio_passphrase?: string;
      // OKX 特定字段
      okx_passphrase?: string;
    };
  };
}

// Competition related types
export interface CompetitionTraderData {
  trader_id: string;
  trader_name: string;
  ai_model: string;
  exchange: string;
  total_equity: number;
  total_pnl: number;
  total_pnl_pct: number;
  position_count: number;
  margin_used_pct: number;
  is_running: boolean;
}

export interface CompetitionData {
  traders: CompetitionTraderData[];
  count: number;
}

// Trader Configuration Data for View Modal
export interface TraderConfigData {
  trader_id?: string;
  trader_name: string;
  ai_model: string;
  exchange_id: string;
  btc_eth_leverage: number;
  altcoin_leverage: number;
  trading_symbols: string;
  custom_prompt: string;
  override_base_prompt: boolean;
  is_cross_margin: boolean;
  use_coin_pool: boolean;
  use_oi_top: boolean;
  initial_balance: number;
  scan_interval_minutes: number;
  is_running: boolean;
  // HODL波段盈利定投策略配置
  strategy?: string;
  strategy_config?: HODLBandProfitConfig;
  // 现货交易配置
  spot_order_type?: 'market' | 'limit';
  spot_position_size_pct?: number;
  spot_take_profit_pct?: number;
  spot_stop_loss_pct?: number;
}
