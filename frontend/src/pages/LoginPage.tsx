import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './LoginPage.css';

// API URLの環境変数（デフォルトはproxyの設定を使用）
const API_URL = process.env.REACT_APP_API_URL || '';

const LoginPage: React.FC = () => {
  const { isAuthenticated, error: authError } = useAuth();
  const [authUrl, setAuthUrl] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  // 認証済みの場合はホームページにリダイレクト
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/');
    }
  }, [isAuthenticated, navigate]);

  // 認可URLを取得
  useEffect(() => {
    const fetchAuthUrl = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${API_URL}/api/auth/url`);
        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.error || 'Failed to get authorization URL');
        }

        setAuthUrl(data.url);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to get authorization URL');
      } finally {
        setLoading(false);
      }
    };

    fetchAuthUrl();
  }, []);

  return (
    <div className="login-page">
      <div className="login-container">
        <h1>Backlog にログイン</h1>
        <p className="description">
          以下のボタンからBacklogアカウントでログインしてください。
        </p>

        {(error || authError) && (
          <div className="error-message">
            {error || authError}
          </div>
        )}

        {loading ? (
          <div className="loading-spinner">Loading...</div>
        ) : (
          <a 
            href={authUrl} 
            className="login-button"
            rel="noopener noreferrer"
          >
            Backlogアカウントでログイン
          </a>
        )}

        {/* <div className="info-box">
          <h3>OAuth 2.0認証について</h3>
          <p>
            ログイン後、Backlogから明示的に許可を求められます。
            アプリはあなたのBacklogデータへの読み取り権限のみを要求します。
          </p>
        </div> */}
      </div>
    </div>
  );
};

export default LoginPage; 