import React, { createContext, useContext, useState, useEffect } from 'react';

interface User {
  id: string;
  email: string;
}

type TradingMode = 'spot' | 'futures' | '';

interface AuthContextType {
  user: User | null;
  token: string | null;
  tradingMode: TradingMode;
  setTradingMode: (mode: TradingMode) => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  // 从 localStorage 恢复交易模式状态
  const [tradingMode, setTradingModeState] = useState<TradingMode>(() => {
    const saved = localStorage.getItem('tradingMode');
    return (saved as TradingMode) || '';
  });
  const [isLoading, setIsLoading] = useState(true);

  // 包装 setTradingMode 以持久化到 localStorage
  const setTradingMode = (mode: TradingMode) => {
    setTradingModeState(mode);
    if (mode) {
      localStorage.setItem('tradingMode', mode);
    } else {
      localStorage.removeItem('tradingMode');
    }
  };

  useEffect(() => {
    // 始终使用管理员模式,无需登录
    setUser({ id: 'admin', email: 'admin@localhost' });
    setToken('admin-mode');
    setIsLoading(false);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        tradingMode,
        setTradingMode,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}