import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import './HomePage.css';

// 型定義
interface BacklogItem {
  id: string;
  projectId: string;
  projectName: string;
  type: string;
  contentSummary: string;
  createdUser: {
    id: string;
    name: string;
  };
  created: string;
  isFavorite: boolean;
}

const HomePage: React.FC = () => {
  const { user, token } = useAuth();
  const [items, setItems] = useState<BacklogItem[]>([]);
  const [favorites, setFavorites] = useState<BacklogItem[]>([]);
  const [keyword, setKeyword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'all' | 'favorites'>('all');

  // 初回ロード時に全データを取得
  useEffect(() => {
    if (user) {
      fetchItems();
      fetchFavorites();
    }
  }, [user]);

  // 検索実行
  const handleSearch = () => {
    fetchItems();
  };

  // 検索フォームのEnterキー対応
  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  // 更新情報データ取得
  const fetchItems = async () => {
    if (!user) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/items?userId=${user.id}&keyword=${keyword}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to fetch items');
      }
      
      setItems(data.items || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch items');
    } finally {
      setLoading(false);
    }
  };

  // お気に入りデータ取得
  const fetchFavorites = async () => {
    if (!user) return;
    
    try {
      const response = await fetch(`/api/favorites/${user.id}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to fetch favorites');
      }
      
      setFavorites(data.items || []);
    } catch (err) {
      console.error('Failed to fetch favorites:', err);
    }
  };

  // お気に入り追加/削除
  const toggleFavorite = async (itemId: string, isFavorite: boolean) => {
    if (!user) return;
    
    try {
      const method = isFavorite ? 'DELETE' : 'POST';
      const response = await fetch(`/api/favorites/${user.id}/${itemId}`, { method });
      
      if (!response.ok) {
        throw new Error('Failed to update favorite');
      }
      
      // 現在の項目の状態を更新
      const updatedItems = items.map(item => 
        item.id === itemId ? { ...item, isFavorite: !isFavorite } : item
      );
      setItems(updatedItems);
      
      // お気に入りリストも更新
      fetchFavorites();
    } catch (err) {
      console.error('Failed to toggle favorite:', err);
    }
  };

  return (
    <div className="home-page">
      <div className="container">
        <h1>Backlog 更新情報検索</h1>
        
        <div className="search-bar">
          <input
            type="text"
            placeholder="キーワードを入力"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onKeyPress={handleKeyPress}
          />
          <button onClick={handleSearch}>検索</button>
        </div>
        
        <div className="tabs">
          <button
            className={`tab ${activeTab === 'all' ? 'active' : ''}`}
            onClick={() => setActiveTab('all')}
          >
            全ての更新情報
          </button>
          <button
            className={`tab ${activeTab === 'favorites' ? 'active' : ''}`}
            onClick={() => setActiveTab('favorites')}
          >
            お気に入り
          </button>
        </div>
        
        {loading && <div className="loading">Loading...</div>}
        
        {error && <div className="error">{error}</div>}
        
        {!loading && !error && (
          <div className="results">
            <table className="items-table">
              <thead>
                <tr>
                  <th className="favorite-col">お気に入り</th>
                  <th>ID</th>
                  <th>プロジェクト</th>
                  <th>種別</th>
                  <th>内容</th>
                  <th>作成者</th>
                  <th>作成日</th>
                </tr>
              </thead>
              <tbody>
                {(activeTab === 'all' ? items : favorites).map((item) => (
                  <tr key={item.id}>
                    <td className="favorite-col">
                      <button
                        className={`favorite-btn ${item.isFavorite ? 'active' : ''}`}
                        onClick={() => toggleFavorite(item.id, item.isFavorite)}
                      >
                        ★
                      </button>
                    </td>
                    <td>{item.id}</td>
                    <td>{item.projectName}</td>
                    <td>{item.type}</td>
                    <td>{item.contentSummary}</td>
                    <td>{item.createdUser.name}</td>
                    <td>{new Date(item.created).toLocaleString()}</td>
                  </tr>
                ))}
                {(activeTab === 'all' ? items : favorites).length === 0 && (
                  <tr>
                    <td colSpan={7} className="no-data">表示する項目がありません</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default HomePage; 