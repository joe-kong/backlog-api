import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './CallbackPage.css';

const CallbackPage: React.FC = () => {
  const { isAuthenticated, loading: authLoading, error: authError } = useAuth();
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

    const processCallback = async () => {
      try {
        const params = new URLSearchParams(location.search);
        const tokenBase64 = params.get('token');
        const userBase64 = params.get('user');
        const errorParam = params.get('error');

        if (errorParam) {
          throw new Error(`Authorization error: ${errorParam}`);
        }

        if (!tokenBase64 || !userBase64) {
          throw new Error('Token or user data is missing');
        }

        // Base64デコード
        const tokenJSON = atob(tokenBase64.replace(/-/g, '+').replace(/_/g, '/'));
        const userJSON = atob(userBase64.replace(/-/g, '+').replace(/_/g, '/'));
        
        // JSONをオブジェクトに変換
        const token = JSON.parse(tokenJSON);
        const user = JSON.parse(userJSON);
        
        // ローカルストレージに保存
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('accessToken', token.accessToken);
        localStorage.setItem('tokenType', token.tokenType);
        localStorage.setItem('expiresAt', token.expiresAt);
        
        // 短い遅延を追加して、状態が更新される時間を確保
        setTimeout(() => {
          // ホームページにリダイレクト
          navigate('/');
        }, 500);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Authentication failed');
      } finally {
        setLoading(false);
      }
    };

    processCallback();
  }, [location, navigate, isAuthenticated, authLoading]);

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