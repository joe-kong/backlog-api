import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';

// 認証情報の型定義
interface User {
  id: string;
  name: string;
  roleType: number;
  lang?: string;
  mailAddress?: string;
}

interface AuthToken {
  accessToken: string;
  tokenType: string;
  expiresAt: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: User | null;
  token: AuthToken | null;
  login: (code: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
  error: string | null;
}

// 初期値
const initialAuthContext: AuthContextType = {
  isAuthenticated: false,
  user: null,
  token: null,
  login: async () => {},
  logout: () => {},
  loading: false,
  error: null,
};

// コンテキストの作成
const AuthContext = createContext<AuthContextType>(initialAuthContext);

// プロパティの型定義
interface AuthProviderProps {
  children: ReactNode;
}

// コンテキストプロバイダー
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<AuthToken | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // ローカルストレージから認証状態を復元
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const storedToken = localStorage.getItem('accessToken');
    const storedTokenType = localStorage.getItem('tokenType');
    const storedExpiresAt = localStorage.getItem('expiresAt');

    if (storedUser && storedToken && storedTokenType && storedExpiresAt) {
      setUser(JSON.parse(storedUser));
      setToken({
        accessToken: storedToken,
        tokenType: storedTokenType,
        expiresAt: storedExpiresAt,
      });
      setIsAuthenticated(true);
    }
  }, []);

  // ログイン処理
  const login = async (code: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`http://localhost:8080/api/auth/callback?code=${code}`);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Authentication failed');
      }

      const { token, user } = data;

      // ローカルストレージに保存
      localStorage.setItem('user', JSON.stringify(user));
      localStorage.setItem('accessToken', token.accessToken);
      localStorage.setItem('tokenType', token.tokenType);
      localStorage.setItem('expiresAt', token.expiresAt);

      setUser(user);
      setToken(token);
      setIsAuthenticated(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  // ログアウト処理
  const logout = () => {
    // ローカルストレージをクリア
    localStorage.removeItem('user');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('tokenType');
    localStorage.removeItem('expiresAt');

    setUser(null);
    setToken(null);
    setIsAuthenticated(false);

    // サーバーサイドのログアウト処理（オプション）
    if (user) {
      fetch(`http://localhost:8080/api/auth/logout/${user.id}`, { method: 'GET' }).catch(console.error);
    }
  };

  const value = {
    isAuthenticated,
    user,
    token,
    login,
    logout,
    loading,
    error,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// カスタムフック
export const useAuth = () => useContext(AuthContext); 