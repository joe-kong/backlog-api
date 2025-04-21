import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Header.css';

const Header: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();

  return (
    <header className="header">
      <div className="header-content">
        <div className="logo">
          <Link to="/">Backlog 更新情報検索アプリ</Link>
        </div>
        <nav className="navigation">
          {isAuthenticated ? (
            <div className="user-info">
              <span className="user-name">{user?.name || 'ユーザー'}</span>
              <button className="logout-button" onClick={logout}>
                ログアウト
              </button>
            </div>
          ) : (
            <Link to="/login" className="login-button">
              ログイン
            </Link>
          )}
        </nav>
      </div>
    </header>
  );
};

export default Header; 