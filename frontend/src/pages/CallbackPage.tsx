import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './CallbackPage.css';

const CallbackPage: React.FC = () => {
  const { login, isAuthenticated, loading: authLoading, error: authError } = useAuth();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // すでに認証済みの場合はホームページにリダイレクト
    if (isAuthenticated && !authLoading) {
      navigate('/');
      return;
    }

    const processOAuthCallback = async () => {
      try {
        const params = new URLSearchParams(location.search);
        const code = params.get('code');
        const errorParam = params.get('error');

        if (errorParam) {
          throw new Error(`Authorization error: ${errorParam}`);
        }

        if (!code) {
          throw new Error('Authorization code is missing');
        }

        // 認可コードを使ってログイン処理を実行
        await login(code);
        navigate('/');
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Authentication failed');
      } finally {
        setLoading(false);
      }
    };

    processOAuthCallback();
  }, [location, login, navigate, isAuthenticated, authLoading]);

  return (
    <div className="callback-page">
      <div className="callback-container">
        <h1>認証処理中...</h1>
        
        {(loading || authLoading) && (
          <div className="loading-area">
            <div className="loading-spinner"></div>
            <p>Backlogアカウントでの認証処理を完了しています</p>
          </div>
        )}
        
        {(error || authError) && (
          <div className="error-area">
            <h3>認証エラー</h3>
            <p className="error-message">{error || authError}</p>
            <button onClick={() => navigate('/login')} className="back-button">
              ログインページに戻る
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default CallbackPage; 