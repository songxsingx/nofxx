import { useState } from 'react';

interface TraderFormProps {
  onCreate: (formData: any) => void;
}

export default function TraderForm({ onCreate }: TraderFormProps) {
  const [form, setForm] = useState({
    exchange: 'binance_spot',
    strategy: 'ai',
    symbol: 'BTCUSDT',
    base_amount_usdt: 100,
    profit_trigger_pct: 10,
    reinvest_ratio: 0.5,
    interval_hours: 24,
    take_profit_pct: 200,
    stop_loss_pct: 30
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  return (
    <form onSubmit={(e) => { e.preventDefault(); onCreate(form); }}>
      <select name="exchange" value={form.exchange} onChange={handleChange}>
        <option value="binance_spot">Binance 现货</option>
        <option value="gateio_spot">Gate.io 现货</option>
        <option value="gateio_futures">Gate.io 合约</option>
        <option value="hyperliquid_spot">Hyperliquid 现货</option>
        <option value="aster_spot">Aster 现货</option>
      </select>

      <select name="strategy" value={form.strategy} onChange={handleChange}>
        <option value="ai">AI 决策</option>
        <option value="hodl_band_profit">波段盈利定投</option>
      </select>

      {form.strategy === 'hodl_band_profit' && (
        <div style={{ border: '1px solid #ccc', padding: '15px', margin: '10px 0' }}>
          <h4>波段盈利定投配置</h4>
          <input name="symbol" placeholder="BTCUSDT" value={form.symbol} onChange={handleChange} />
          <input type="number" name="base_amount_usdt" placeholder="初始本金 USDT" value={form.base_amount_usdt} onChange={handleChange} />
          <input type="number" name="profit_trigger_pct" placeholder="触发盈利 %" value={form.profit_trigger_pct} onChange={handleChange} />
          <input type="number" name="reinvest_ratio" placeholder="再投比例 (0.5=50%)" value={form.reinvest_ratio} onChange={handleChange} />
          <input type="number" name="interval_hours" placeholder="间隔小时" value={form.interval_hours} onChange={handleChange} />
          <input type="number" name="take_profit_pct" placeholder="总止盈 %" value={form.take_profit_pct} onChange={handleChange} />
          <input type="number" name="stop_loss_pct" placeholder="总止损 %" value={form.stop_loss_pct} onChange={handleChange} />
        </div>
      )}

      <button type="submit">创建交易员</button>
    </form>
  );
}
