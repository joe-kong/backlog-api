import React from 'react';
import { Link } from 'react-router-dom';
import './NotFoundPage.css';

const NotFoundPage: React.FC = () => {
  return (
    <div className="not-found-page">
      <div className="not-found-container">
        <h1>404</h1>
        <h2>ページが見つかりません</h2>
        <p>
          お探しのページは存在しないか、移動した可能性があります。
        </p>
        <Link to="/" className="home-button">
          ホームページに戻る
        </Link>
      </div>
    </div>
  );
};

export default NotFoundPage; 