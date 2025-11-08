import { useState, useEffect } from 'react';
import type { AIModel, Exchange, CreateTraderRequest, HODLBandProfitConfig } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { t } from '../i18n/translations';

// 提取下划线后面的名称部分
function getShortName(fullName: string): string {
  const parts = fullName.split('_');
  return parts.length > 1 ? parts[parts.length - 1] : fullName;
}

interface TraderConfigData {
  trader_id?: string;
  trader_name: string;
  ai_model: string;
  exchange_id: string;
  btc_eth_leverage: number;
  altcoin_leverage: number;
  trading_symbols: string;
  custom_prompt: string;
  override_base_prompt: boolean;
  system_prompt_template: string;
  is_cross_margin: boolean;
  use_coin_pool: boolean;
  use_oi_top: boolean;
  initial_balance: number;
  scan_interval_minutes: number;
  // 现货特有配置
  spot_order_type?: 'market' | 'limit';  // 订单类型：市价单/限价单
  spot_position_size_pct?: number;       // 每次交易使用余额百分比（0-100）
  spot_take_profit_pct?: number;         // 止盈百分比
  spot_stop_loss_pct?: number;           // 止损百分比
  // HODL波段盈利定投策略配置
  strategy?: string;                     // "ai" or "hodl_band_profit"
  strategy_config?: HODLBandProfitConfig;
}

interface TraderConfigModalProps {
  isOpen: boolean;
  onClose: () => void;
  traderData?: TraderConfigData | null;
  isEditMode?: boolean;
  availableModels?: AIModel[];
  availableExchanges?: Exchange[];
  onSave?: (data: CreateTraderRequest) => Promise<void>;
  tradingMode?: 'spot' | 'futures' | '';
}

export function TraderConfigModal({ 
  isOpen, 
  onClose, 
  traderData, 
  isEditMode = false,
  availableModels = [],
  availableExchanges = [],
  onSave,
  tradingMode = ''
}: TraderConfigModalProps) {
  const { language } = useLanguage();
  
  // 根据交易模式过滤交易所，如果未设置模式则显示全部
  const filteredExchanges = !tradingMode ? availableExchanges : availableExchanges.filter(exchange => {
    const isSpot = exchange.id?.includes('_spot');
    return tradingMode === 'spot' ? isSpot : !isSpot;
  });
  const [formData, setFormData] = useState<TraderConfigData>({
    trader_name: '',
    ai_model: '',
    exchange_id: '',
    btc_eth_leverage: 5,
    altcoin_leverage: 3,
    trading_symbols: '',
    custom_prompt: '',
    override_base_prompt: false,
    system_prompt_template: 'default',
    is_cross_margin: true,
    use_coin_pool: false,
    use_oi_top: false,
    initial_balance: 1000,
    scan_interval_minutes: 3,
    // 现货默认配置
    spot_order_type: 'market',
    spot_position_size_pct: 100,
    spot_take_profit_pct: 20,
    spot_stop_loss_pct: 10,
    // HODL波段盈利定投策略默认配置
    strategy: 'ai',
    strategy_config: {
      symbol: '',
      base_amount_usdt: 100,
      profit_trigger_pct: 10,
      reinvest_ratio: 0.5,
      interval_hours: 1,
      take_profit_pct: 100,
      stop_loss_pct: 10,
    },
  });
  const [isSaving, setIsSaving] = useState(false);
  const [availableCoins, setAvailableCoins] = useState<string[]>([]);
  const [selectedCoins, setSelectedCoins] = useState<string[]>([]);
  const [showCoinSelector, setShowCoinSelector] = useState(false);
  const [promptTemplates, setPromptTemplates] = useState<{name: string}[]>([]);

  // 判断当前选择的交易所是否为现货
  const selectedExchange = availableExchanges.find(e => e.id === formData.exchange_id);
  const isSpotExchange = selectedExchange?.id?.includes('_spot') || false;

  // 首次打开模态框时初始化默认值
  useEffect(() => {
    if (isOpen && !isEditMode && !traderData) {
      const defaultExchangeId = filteredExchanges.length > 0 ? filteredExchanges[0].id : '';
      const defaultModelId = availableModels.length > 0 ? availableModels[0].id : '';
      
      setFormData({
        trader_name: '',
        ai_model: defaultModelId,
        exchange_id: defaultExchangeId,
        btc_eth_leverage: 5,
        altcoin_leverage: 3,
        trading_symbols: '',
        custom_prompt: '',
        override_base_prompt: false,
        system_prompt_template: 'default',
        is_cross_margin: true,
        use_coin_pool: false,
        use_oi_top: false,
        initial_balance: 1000,
        scan_interval_minutes: 3,
        // 现货默认配置
        spot_order_type: 'market',
        spot_position_size_pct: 100,
        spot_take_profit_pct: 20,
        spot_stop_loss_pct: 10,
        // HODL波段盈利定投策略默认配置
        strategy: 'ai',
        strategy_config: {
          symbol: '',
          base_amount_usdt: 100,
          profit_trigger_pct: 10,
          reinvest_ratio: 0.5,
          interval_hours: 1,
          take_profit_pct: 100,
          stop_loss_pct: 10,
        },
      });
    }
  }, [isOpen, isEditMode, traderData]);

  useEffect(() => {
    if (traderData) {
      setFormData(traderData);
      // 设置已选择的币种
      if (traderData.trading_symbols) {
        const coins = traderData.trading_symbols.split(',').map(s => s.trim()).filter(s => s);
        setSelectedCoins(coins);
      }
    } else if (!isEditMode && isOpen) {
      // 只在首次打开模态框时初始化，不要在filteredExchanges变化时重置
      const defaultExchangeId = filteredExchanges.length > 0 ? filteredExchanges[0].id : '';
      const defaultModelId = availableModels.length > 0 ? availableModels[0].id : '';
      
      setFormData(prev => ({
        ...prev,
        ai_model: prev.ai_model || defaultModelId,
        exchange_id: prev.exchange_id || defaultExchangeId,
      }));
    }
    // 确保旧数据也有默认的 system_prompt_template
    if (traderData && !traderData.system_prompt_template) {
      setFormData(prev => ({
        ...prev,
        system_prompt_template: 'default'
      }));
    } else if (!traderData && !isEditMode && isOpen) {
      // 初始化时设置默认的 system_prompt_template
      setFormData(prev => ({
        ...prev,
        system_prompt_template: 'default'
      }));
    }
  }, [traderData, isEditMode, isOpen]);

  // 获取系统配置中的币种列表
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch('/api/config');
        const config = await response.json();
        if (config.default_coins) {
          setAvailableCoins(config.default_coins);
        }
      } catch (error) {
        console.error('Failed to fetch config:', error);
        // 使用默认币种列表
        setAvailableCoins(['BTCUSDT', 'ETHUSDT', 'SOLUSDT', 'BNBUSDT', 'XRPUSDT', 'DOGEUSDT', 'ADAUSDT']);
      }
    };
    fetchConfig();
  }, []);

  // 获取系统提示词模板列表
  useEffect(() => {
    const fetchPromptTemplates = async () => {
      try {
        const response = await fetch('/api/prompt-templates');
        const data = await response.json();
        if (data.templates) {
          setPromptTemplates(data.templates);
        }
      } catch (error) {
        console.error('Failed to fetch prompt templates:', error);
        // 使用默认模板列表
        setPromptTemplates([{name: 'default'}, {name: 'aggressive'}]);
      }
    };
    fetchPromptTemplates();
  }, []);

  // 当选择的币种改变时，更新输入框
  useEffect(() => {
    const symbolsString = selectedCoins.join(',');
    setFormData(prev => ({ ...prev, trading_symbols: symbolsString }));
  }, [selectedCoins]);

  if (!isOpen) return null;

  console.log('🔍 TraderConfigModal 渲染:', {
    tradingMode,
    availableModels: availableModels.length,
    availableExchanges: availableExchanges.length,
    filteredExchanges: filteredExchanges.length,
    formData_ai_model: formData.ai_model,
    formData_exchange_id: formData.exchange_id,
  });

  // 检查是否有可用的模型和交易所
  if (availableModels.length === 0 || filteredExchanges.length === 0) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm">
        <div className="bg-[#1E2329] border border-[#2B3139] rounded-xl shadow-2xl max-w-md w-full mx-4 p-6">
          <div className="text-center">
            <div className="w-16 h-16 rounded-full bg-red-500/10 flex items-center justify-center mx-auto mb-4">
              <span className="text-4xl">⚠️</span>
            </div>
            <h3 className="text-xl font-bold text-[#EAECEF] mb-2">无法创建交易员</h3>
            <div className="text-[#848E9C] mb-6 space-y-2">
              {availableModels.length === 0 && (
                <p>• 请先配置 AI 模型</p>
              )}
              {filteredExchanges.length === 0 && tradingMode === 'spot' && (
                <p>• 请先配置现货交易所（带 _spot 后缀）</p>
              )}
              {filteredExchanges.length === 0 && tradingMode === 'futures' && (
                <p>• 请先配置合约交易所</p>
              )}
              {filteredExchanges.length === 0 && !tradingMode && (
                <p>• 请先配置交易所</p>
              )}
            </div>
            <button
              onClick={onClose}
              className="w-full px-4 py-2 bg-[#F0B90B] text-black rounded-lg font-semibold hover:bg-[#E1A706] transition-colors"
            >
              确定
            </button>
          </div>
        </div>
      </div>
    );
  }

  const handleInputChange = (field: keyof TraderConfigData, value: any) => {
    console.log('📝 输入变更:', field, '=', value);
    setFormData(prev => ({ ...prev, [field]: value }));
    
    // 如果是直接编辑trading_symbols，同步更新selectedCoins
    if (field === 'trading_symbols') {
      const coins = value.split(',').map((s: string) => s.trim()).filter((s: string) => s);
      setSelectedCoins(coins);
    }
  };

  // 处理HODL策略配置变更
  const handleStrategyConfigChange = (field: keyof HODLBandProfitConfig, value: any) => {
    setFormData(prev => ({
      ...prev,
      strategy_config: {
        ...prev.strategy_config!,
        [field]: value
      }
    }));
  };

  const handleCoinToggle = (coin: string) => {
    setSelectedCoins(prev => {
      if (prev.includes(coin)) {
        return prev.filter(c => c !== coin);
      } else {
        return [...prev, coin];
      }
    });
  };

  const handleSave = async () => {
    if (!onSave) return;

    setIsSaving(true);
    try {
      const saveData: CreateTraderRequest = {
        name: formData.trader_name,
        ai_model_id: formData.ai_model,
        exchange_id: formData.exchange_id,
        btc_eth_leverage: formData.btc_eth_leverage || 5,
        altcoin_leverage: formData.altcoin_leverage || 3,
        trading_symbols: formData.trading_symbols || '',
        custom_prompt: formData.custom_prompt || '',
        override_base_prompt: formData.override_base_prompt || false,
        system_prompt_template: formData.system_prompt_template || 'default',
        is_cross_margin: formData.is_cross_margin !== undefined ? formData.is_cross_margin : true,
        use_coin_pool: formData.use_coin_pool || false,
        use_oi_top: formData.use_oi_top || false,
        initial_balance: formData.initial_balance || 1000,
        scan_interval_minutes: formData.scan_interval_minutes || 3,
        // HODL波段盈利定投策略配置
        strategy: formData.strategy || 'ai',
        strategy_config: formData.strategy === 'hodl_band_profit' ? formData.strategy_config : undefined,
        // 现货交易配置
        spot_order_type: formData.spot_order_type || 'market',
        spot_position_size_pct: formData.spot_position_size_pct || 100,
        spot_take_profit_pct: formData.spot_take_profit_pct || 20,
        spot_stop_loss_pct: formData.spot_stop_loss_pct || 10,
      };
      
      console.log('📤 发送保存请求:', saveData);
      
      await onSave(saveData);
      onClose();
    } catch (error) {
      console.error('保存失败:', error);
      alert(`保存失败: ${error instanceof Error ? error.message : '未知错误'}`);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm">
      <div 
        className="bg-[#1E2329] border border-[#2B3139] rounded-xl shadow-2xl max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-[#2B3139] bg-gradient-to-r from-[#1E2329] to-[#252B35]">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[#F0B90B] to-[#E1A706] flex items-center justify-center">
              <span className="text-lg">{isEditMode ? '✏️' : '➕'}</span>
            </div>
            <div>
              <h2 className="text-xl font-bold text-[#EAECEF]">
                {isEditMode ? '修改交易员' : '创建交易员'}
              </h2>
              <p className="text-sm text-[#848E9C] mt-1">
                {isEditMode ? '修改交易员配置参数' : '配置新的AI交易员'}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="w-8 h-8 rounded-lg text-[#848E9C] hover:text-[#EAECEF] hover:bg-[#2B3139] transition-colors flex items-center justify-center"
          >
            ✕
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-8">
          {/* Basic Info */}
          <div className="bg-[#0B0E11] border border-[#2B3139] rounded-lg p-5">
            <h3 className="text-lg font-semibold text-[#EAECEF] mb-5 flex items-center gap-2">
              🤖 基础配置
            </h3>
            <div className="space-y-4">
              <div>
                <label className="text-sm text-[#EAECEF] block mb-2">交易员名称</label>
                <input
                  type="text"
                  value={formData.trader_name}
                  onChange={(e) => handleInputChange('trader_name', e.target.value)}
                  className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                  placeholder="请输入交易员名称"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm text-[#EAECEF] block mb-2">AI模型</label>
                  <select
                    value={formData.ai_model}
                    onChange={(e) => handleInputChange('ai_model', e.target.value)}
                    className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                  >
                    {availableModels.map(model => (
                      <option key={model.id} value={model.id}>
                        {getShortName(model.name || model.id).toUpperCase()}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-sm text-[#EAECEF] block mb-2">交易所</label>
                  <select
                    value={formData.exchange_id}
                    onChange={(e) => handleInputChange('exchange_id', e.target.value)}
                    className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                  >
                    {filteredExchanges.map(exchange => (
                      <option key={exchange.id} value={exchange.id}>
                        {getShortName(exchange.name || exchange.id).toUpperCase()}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            </div>
          </div>

          {/* Trading Configuration */}
          <div className="bg-[#0B0E11] border border-[#2B3139] rounded-lg p-5">
            <h3 className="text-lg font-semibold text-[#EAECEF] mb-5 flex items-center gap-2">
              ⚖️ 交易配置
            </h3>
            <div className="space-y-4">
              {/* 第一行：保证金模式和初始余额 */}
              <div className="grid grid-cols-2 gap-4">
                {/* 仅合约交易所显示保证金模式 */}
                {!isSpotExchange && (
                  <div>
                    <label className="text-sm text-[#EAECEF] block mb-2">保证金模式</label>
                    <div className="flex gap-2">
                      <button
                        type="button"
                        onClick={() => handleInputChange('is_cross_margin', true)}
                        className={`flex-1 px-3 py-2 rounded text-sm ${
                          formData.is_cross_margin 
                            ? 'bg-[#F0B90B] text-black' 
                            : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                        }`}
                      >
                        全仓
                      </button>
                      <button
                        type="button"
                        onClick={() => handleInputChange('is_cross_margin', false)}
                        className={`flex-1 px-3 py-2 rounded text-sm ${
                          !formData.is_cross_margin 
                            ? 'bg-[#F0B90B] text-black' 
                            : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                        }`}
                      >
                        逐仓
                      </button>
                    </div>
                  </div>
                )}
                <div>
                  <label className="text-sm text-[#EAECEF] block mb-2">初始余额 ($)</label>
                  <input
                    type="number"
                    value={formData.initial_balance}
                    onChange={(e) => handleInputChange('initial_balance', Number(e.target.value))}
                    className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                    min="100"
                    step="100"
                  />
                </div>
              </div>

              {/* 第二行：AI 扫描决策间隔 */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm text-[#EAECEF] block mb-2">{t('aiScanInterval', language)}</label>
                  <input
                    type="number"
                    value={formData.scan_interval_minutes}
                    onChange={(e) => handleInputChange('scan_interval_minutes', Number(e.target.value))}
                    className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                    min="1"
                    max="60"
                    step="1"
                  />
                  <p className="text-xs text-gray-500 mt-1">{t('scanIntervalRecommend', language)}</p>
                </div>
                <div></div>
              </div>

              {/* 第三行：杠杆设置（仅合约交易所显示） */}
              {!isSpotExchange && (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm text-[#EAECEF] block mb-2">BTC/ETH 杠杆</label>
                    <input
                      type="number"
                      value={formData.btc_eth_leverage}
                      onChange={(e) => handleInputChange('btc_eth_leverage', Number(e.target.value))}
                      className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                      min="1"
                      max="125"
                    />
                  </div>
                  <div>
                    <label className="text-sm text-[#EAECEF] block mb-2">山寨币杠杆</label>
                    <input
                      type="number"
                      value={formData.altcoin_leverage}
                      onChange={(e) => handleInputChange('altcoin_leverage', Number(e.target.value))}
                      className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                      min="1"
                      max="75"
                    />
                  </div>
                </div>
              )}

              {/* 第三行：交易币种 */}
              <div>
                <div className="flex items-center justify-between mb-2">
                  <label className="text-sm text-[#EAECEF]">交易币种 (用逗号分隔，留空使用默认)</label>
                  <button
                    type="button"
                    onClick={() => setShowCoinSelector(!showCoinSelector)}
                    className="px-3 py-1 text-xs bg-[#F0B90B] text-black rounded hover:bg-[#E1A706] transition-colors"
                  >
                    {showCoinSelector ? '收起选择' : '快速选择'}
                  </button>
                </div>
                <input
                  type="text"
                  value={formData.trading_symbols}
                  onChange={(e) => handleInputChange('trading_symbols', e.target.value)}
                  className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                  placeholder="例如: BTCUSDT,ETHUSDT,ADAUSDT"
                />
                
                {/* 币种选择器 */}
                {showCoinSelector && (
                  <div className="mt-3 p-3 bg-[#0B0E11] border border-[#2B3139] rounded">
                    <div className="text-xs text-[#848E9C] mb-2">点击选择币种：</div>
                    <div className="flex flex-wrap gap-2">
                      {availableCoins.map(coin => (
                        <button
                          key={coin}
                          type="button"
                          onClick={() => handleCoinToggle(coin)}
                          className={`px-2 py-1 text-xs rounded transition-colors ${
                            selectedCoins.includes(coin)
                              ? 'bg-[#F0B90B] text-black'
                              : 'bg-[#1E2329] text-[#848E9C] border border-[#2B3139] hover:border-[#F0B90B]'
                          }`}
                        >
                          {coin.replace('USDT', '')}
                        </button>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* 现货特有配置（仅现货交易所显示） */}
              {isSpotExchange && (
                <>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="text-sm text-[#EAECEF] block mb-2">订单类型</label>
                      <div className="flex gap-2">
                        <button
                          type="button"
                          onClick={() => handleInputChange('spot_order_type', 'market')}
                          className={`flex-1 px-3 py-2 rounded text-sm ${
                            formData.spot_order_type === 'market'
                              ? 'bg-[#F0B90B] text-black' 
                              : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                          }`}
                        >
                          市价单
                        </button>
                        <button
                          type="button"
                          onClick={() => handleInputChange('spot_order_type', 'limit')}
                          className={`flex-1 px-3 py-2 rounded text-sm ${
                            formData.spot_order_type === 'limit'
                              ? 'bg-[#F0B90B] text-black' 
                              : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                          }`}
                        >
                          限价单
                        </button>
                      </div>
                      <p className="text-xs text-[#848E9C] mt-1">市价单快速成交，限价单可设定价格</p>
                    </div>
                    <div>
                      <label className="text-sm text-[#EAECEF] block mb-2">仓位比例 (%)</label>
                      <input
                        type="number"
                        value={formData.spot_position_size_pct}
                        onChange={(e) => handleInputChange('spot_position_size_pct', Number(e.target.value))}
                        className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                        min="1"
                        max="100"
                        step="1"
                      />
                      <p className="text-xs text-[#848E9C] mt-1">每次交易使用余额的百分比</p>
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="text-sm text-[#EAECEF] block mb-2">止盈比例 (%)</label>
                      <input
                        type="number"
                        value={formData.spot_take_profit_pct}
                        onChange={(e) => handleInputChange('spot_take_profit_pct', Number(e.target.value))}
                        className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                        min="1"
                        max="1000"
                        step="1"
                      />
                      <p className="text-xs text-[#848E9C] mt-1">达到该涨幅后自动卖出</p>
                    </div>
                    <div>
                      <label className="text-sm text-[#EAECEF] block mb-2">止损比例 (%)</label>
                      <input
                        type="number"
                        value={formData.spot_stop_loss_pct}
                        onChange={(e) => handleInputChange('spot_stop_loss_pct', Number(e.target.value))}
                        className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                        min="1"
                        max="100"
                        step="1"
                      />
                      <p className="text-xs text-[#848E9C] mt-1">跌破该比例后自动卖出</p>
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* HODL 波段盈利定投策略配置（仅现货交易所显示） */}
          {isSpotExchange && (
            <div className="bg-[#0B0E11] border border-[#2B3139] rounded-lg p-5">
              <h3 className="text-lg font-semibold text-[#EAECEF] mb-5 flex items-center gap-2">
                💰 HODL 单币波段盈利定投
              </h3>
              <div className="space-y-4">
                {/* 策略开关 */}
                <div className="flex items-center justify-between p-4 bg-[#1E2329] rounded-lg border border-[#2B3139]">
                  <div className="flex-1">
                    <div className="text-sm font-semibold text-[#EAECEF] mb-1">启用波段盈利定投策略</div>
                    <div className="text-xs text-[#848E9C]">
                      只囤一个币（如BTC），盈利{formData.strategy_config?.profit_trigger_pct || 10}% → {(formData.strategy_config?.reinvest_ratio || 0.5) * 100}% 再投，本金不动，盈利滚雪球
                    </div>
                    <div className="text-xs text-green-500 mt-1">
                      ✅ 双轨模式：HODL后台监控（囤1小时检查） + AI主动交易并行运行
                    </div>
                    <div className="text-xs text-yellow-500 mt-1">
                      ⚠️ 合约交易员不支持HODL策略（有爆仓风险）
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      onClick={() => handleInputChange('strategy', 'hodl_band_profit')}
                      className={`px-4 py-2 rounded text-sm ${
                        formData.strategy === 'hodl_band_profit'
                          ? 'bg-[#F0B90B] text-black' 
                          : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                      }`}
                    >
                      启用
                    </button>
                    <button
                      type="button"
                      onClick={() => handleInputChange('strategy', 'ai')}
                      className={`px-4 py-2 rounded text-sm ${
                        formData.strategy === 'ai'
                          ? 'bg-[#F0B90B] text-black' 
                          : 'bg-[#0B0E11] text-[#848E9C] border border-[#2B3139]'
                      }`}
                    >
                      关闭
                    </button>
                  </div>
                </div>

                {/* 策略详细配置（仅启用时显示） */}
                {formData.strategy === 'hodl_band_profit' && (
                  <>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">
                          目标币种 <span className="text-[#848E9C]">(可选)</span>
                        </label>
                        <input
                          type="text"
                          value={formData.strategy_config?.symbol || ''}
                          onChange={(e) => handleStrategyConfigChange('symbol', e.target.value.toUpperCase())}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          placeholder="如：BTCUSDT（留空则AI自行决策）"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">留空则使用AI波段盈利策略</p>
                      </div>
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">初始投入金额 ($)</label>
                        <input
                          type="number"
                          value={formData.strategy_config?.base_amount_usdt || 100}
                          onChange={(e) => handleStrategyConfigChange('base_amount_usdt', Number(e.target.value))}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="10"
                          step="10"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">首次买入使用的USDT金额</p>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">盈利触发比例 (%)</label>
                        <input
                          type="number"
                          value={formData.strategy_config?.profit_trigger_pct || 10}
                          onChange={(e) => handleStrategyConfigChange('profit_trigger_pct', Number(e.target.value))}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="1"
                          max="1000"
                          step="1"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">盈利达到此比例时触发再投资</p>
                      </div>
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">再投资比例 (%)</label>
                        <input
                          type="number"
                          value={(formData.strategy_config?.reinvest_ratio || 0.5) * 100}
                          onChange={(e) => handleStrategyConfigChange('reinvest_ratio', Number(e.target.value) / 100)}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="1"
                          max="100"
                          step="1"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">盈利中用于再投资的比例</p>
                      </div>
                    </div>

                    <div className="grid grid-cols-3 gap-4">
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">检查间隔 (小时)</label>
                        <input
                          type="number"
                          value={formData.strategy_config?.interval_hours || 1}
                          onChange={(e) => handleStrategyConfigChange('interval_hours', Number(e.target.value))}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="1"
                          max="24"
                          step="1"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">执行策略的时间间隔</p>
                      </div>
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">止盈比例 (%)</label>
                        <input
                          type="number"
                          value={formData.strategy_config?.take_profit_pct || 100}
                          onChange={(e) => handleStrategyConfigChange('take_profit_pct', Number(e.target.value))}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="1"
                          max="1000"
                          step="1"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">全部持仓止盈比例</p>
                      </div>
                      <div>
                        <label className="text-sm text-[#EAECEF] block mb-2">止损比例 (%)</label>
                        <input
                          type="number"
                          value={formData.strategy_config?.stop_loss_pct || 10}
                          onChange={(e) => handleStrategyConfigChange('stop_loss_pct', Number(e.target.value))}
                          className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                          min="1"
                          max="100"
                          step="1"
                        />
                        <p className="text-xs text-[#848E9C] mt-1">全部持仓止损比例</p>
                      </div>
                    </div>

                    {/* 策略说明 */}
                    <div className="p-4 bg-[#1E2329] rounded-lg border border-[#2B3139]">
                      <div className="text-xs text-[#848E9C] space-y-1">
                        <div className="flex items-start gap-2">
                          <span className="text-[#F0B90B]">💡</span>
                          <span>策略模式：{formData.strategy_config?.symbol ? `单币囤币（${formData.strategy_config.symbol}）` : 'AI智能波段（多币种自动选择）'}</span>
                        </div>
                        <div className="flex items-start gap-2">
                          <span className="text-[#F0B90B]">📈</span>
                          <span>盈利{formData.strategy_config?.profit_trigger_pct || 10}%时，将盈利的{(formData.strategy_config?.reinvest_ratio || 0.5) * 100}%再投入，实现复利增长</span>
                        </div>
                        <div className="flex items-start gap-2">
                          <span className="text-[#F0B90B]">🔒</span>
                          <span>本金不动，仅用盈利滚雪球，降低风险</span>
                        </div>
                      </div>
                    </div>
                  </>
                )}
              </div>
            </div>
          )}

          {/* Signal Sources */}
          <div className="bg-[#0B0E11] border border-[#2B3139] rounded-lg p-5">
            <h3 className="text-lg font-semibold text-[#EAECEF] mb-5 flex items-center gap-2">
              📡 信号源配置
            </h3>
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  checked={formData.use_coin_pool}
                  onChange={(e) => handleInputChange('use_coin_pool', e.target.checked)}
                  className="w-4 h-4"
                />
                <label className="text-sm text-[#EAECEF]">使用 Coin Pool 信号</label>
              </div>
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  checked={formData.use_oi_top}
                  onChange={(e) => handleInputChange('use_oi_top', e.target.checked)}
                  className="w-4 h-4"
                />
                <label className="text-sm text-[#EAECEF]">使用 OI Top 信号</label>
              </div>
            </div>
          </div>

          {/* Trading Prompt */}
          <div className="bg-[#0B0E11] border border-[#2B3139] rounded-lg p-5">
            <h3 className="text-lg font-semibold text-[#EAECEF] mb-5 flex items-center gap-2">
              💬 交易策略提示词
            </h3>
            <div className="space-y-4">
              {/* 系统提示词模板选择 */}
              <div>
                <label className="text-sm text-[#EAECEF] block mb-2">系统提示词模板</label>
                <select
                  value={formData.system_prompt_template}
                  onChange={(e) => handleInputChange('system_prompt_template', e.target.value)}
                  className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none"
                >
                  {promptTemplates.map(template => (
                    <option key={template.name} value={template.name}>
                      {template.name === 'default' ? 'Default (默认稳健)' :
                       template.name === 'aggressive' ? 'Aggressive (激进)' :
                       template.name === 'spot' ? 'Spot (现货专用 - 无杠杆长期持有)' :
                       template.name === 'adaptive' ? 'Adaptive (自适应)' :
                       template.name === 'nextrade' ? 'NexTrade (极简主义)' :
                       template.name === 'taro_long_prompts' ? 'Taro (高级策略)' :
                       template.name.charAt(0).toUpperCase() + template.name.slice(1)}
                    </option>
                  ))}
                </select>
                <p className="text-xs text-[#848E9C] mt-1">
                  选择预设的交易策略模板（包含交易哲学、风控原则等）
                </p>
              </div>

              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  checked={formData.override_base_prompt}
                  onChange={(e) => handleInputChange('override_base_prompt', e.target.checked)}
                  className="w-4 h-4"
                />
                <label className="text-sm text-[#EAECEF]">覆盖默认提示词</label>
                <span className="text-xs text-[#F0B90B] inline-flex items-center gap-1"><svg xmlns="http://www.w3.org/2000/svg" className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0Z"/><line x1="12" x2="12" y1="9" y2="13"/><line x1="12" x2="12.01" y1="17" y2="17"/></svg> 启用后将完全替换默认策略</span>
              </div>
              <div>
                <label className="text-sm text-[#EAECEF] block mb-2">
                  {formData.override_base_prompt ? '自定义提示词' : '附加提示词'}
                </label>
                <textarea
                  value={formData.custom_prompt}
                  onChange={(e) => handleInputChange('custom_prompt', e.target.value)}
                  className="w-full px-3 py-2 bg-[#0B0E11] border border-[#2B3139] rounded text-[#EAECEF] focus:border-[#F0B90B] focus:outline-none h-24 resize-none"
                  placeholder={formData.override_base_prompt ? "输入完整的交易策略提示词..." : "输入额外的交易策略提示..."}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex justify-end gap-3 p-6 border-t border-[#2B3139] bg-gradient-to-r from-[#1E2329] to-[#252B35]">
          <button
            onClick={onClose}
            className="px-6 py-3 bg-[#2B3139] text-[#EAECEF] rounded-lg hover:bg-[#404750] transition-all duration-200 border border-[#404750]"
          >
            取消
          </button>
          {onSave && (
            <button
              onClick={handleSave}
              disabled={isSaving || !formData.trader_name || !formData.ai_model || !formData.exchange_id}
              className="px-8 py-3 bg-gradient-to-r from-[#F0B90B] to-[#E1A706] text-black rounded-lg hover:from-[#E1A706] hover:to-[#D4951E] transition-all duration-200 disabled:bg-[#848E9C] disabled:cursor-not-allowed font-medium shadow-lg"
            >
              {isSaving ? '保存中...' : (isEditMode ? '保存修改' : '创建交易员')}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
