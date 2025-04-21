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
  isUpdating?: boolean;
}

// AI分析結果の型定義
interface AIAnalysis {
  summary: string;
  keyPoints: string[];
  nextActions: string[];
}

const HomePage: React.FC = () => {
  const { user, token } = useAuth();
  const [items, setItems] = useState<BacklogItem[]>([]);
  const [favorites, setFavorites] = useState<BacklogItem[]>([]);
  const [keyword, setKeyword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  
  // AI分析関連の状態
  const [selectedItem, setSelectedItem] = useState<BacklogItem | null>(null);
  const [aiAnalysis, setAiAnalysis] = useState<AIAnalysis | null>({
    summary: "",
    keyPoints: [],
    nextActions: []
  });
  const [aiLoading, setAiLoading] = useState<boolean>(false);
  const [showAiModal, setShowAiModal] = useState<boolean>(false);

  // 開発用デバッグ情報表示
  const [showDebug, setShowDebug] = useState<boolean>(false);

  // 初回ロード時にお気に入りデータを先に取得し、その後に全データを取得
  useEffect(() => {
    if (user) {
      const initializeData = async () => {
        setLoading(true);
        try {
          await fetchFavorites();
          await fetchItems();
        } catch (err) {
          console.error('初期化中にエラーが発生しました:', err);
        } finally {
          setLoading(false);
        }
      };
      
      initializeData();
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
    
    try {
      const response = await fetch(`/api/items?userId=${user.id}&keyword=${keyword}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to fetch items');
      }
      
      // サーバーから返されたisFavoriteフラグをそのまま使用
      const updatedItems = data.items || [];
      
      setItems(updatedItems);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch items');
    }
  };

  // お気に入りデータ取得
  const fetchFavorites = async () => {
    if (!user) return;
    
    try {
      // お気に入りデータの取得
      const response = await fetch(`/api/favorites/${user.id}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch favorites');
      }
      
      const data = await response.json();
      console.log('サーバーから取得したお気に入りデータ:', data);
      
      // 空の配列の場合は早期リターン
      if (!data.items || !Array.isArray(data.items) || data.items.length === 0) {
        console.log('お気に入りデータが空です');
        setFavorites([]);
        return;
      }

      // データ検証
      const validItems = data.items.filter((item: any) => {
        const isValid = item && 
                        item.id && 
                        item.projectName && 
                        item.type && 
                        item.contentSummary && 
                        item.createdUser && 
                        item.createdUser.name;
        
        if (!isValid) {
          console.warn('不完全なお気に入りアイテムをスキップします:', item);
        }
        
        return isValid;
      });
      
      console.log('検証済みお気に入りアイテム数:', validItems.length);
      
      // お気に入りフラグを追加
      const favoritesWithFlag = validItems.map((item: BacklogItem) => ({
        ...item,
        isFavorite: true
      }));
      
      // サーバーから取得したお気に入りデータをセット
      setFavorites(favoritesWithFlag);
      
      // 全アイテムのお気に入り状態も更新（すでにアイテムが読み込まれている場合）
      if (items.length > 0) {
        const updatedItems = items.map(item => {
          const isFavorite = favoritesWithFlag.some((fav: BacklogItem) => fav.id === item.id);
          return { ...item, isFavorite };
        });
        setItems(updatedItems);
      }
    } catch (err) {
      console.error('Failed to fetch favorites:', err);
      // お気に入りの取得に失敗した場合は空の配列をセット
      setFavorites([]);
    }
  };

  // お気に入り追加/削除
  const toggleFavorite = async (itemId: string, isFavorite: boolean) => {
    if (!user) return;
    
    try {
      const method = isFavorite ? 'DELETE' : 'POST';
      
      // 一時的にローディング状態を表示
      const tempUpdatedItems = items.map(item => 
        item.id === itemId ? { ...item, isUpdating: true } : item
      );
      setItems(tempUpdatedItems);
      
      // APIコールを実行し、完了するまで待機
      const response = await fetch(`/api/favorites/${user.id}/${itemId}`, { method });
      
      // APIコールが成功した場合のみ状態を更新
      if (response.ok) {
        // 現在の項目の状態だけを更新
        const updatedItems = items.map(item => 
          item.id === itemId ? { ...item, isFavorite: !isFavorite, isUpdating: false } : item
        );
        setItems(updatedItems);
        
        // お気に入りリストを直接更新（APIコールなし）
        if (isFavorite) {
          // お気に入りから削除する場合
          setFavorites(prev => prev.filter(item => item.id !== itemId));
        } else {
          // お気に入りに追加する場合
          const itemToAdd = items.find(item => item.id === itemId);
          if (itemToAdd) {
            setFavorites(prev => [...prev, {...itemToAdd, isFavorite: true}]);
          }
        }
      } else {
        // API呼び出しが失敗した場合は元の状態に戻す
        const revertedItems = items.map(item => 
          item.id === itemId ? { ...item, isUpdating: false } : item
        );
        setItems(revertedItems);
        
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to update favorite');
      }
    } catch (err) {
      console.error('Failed to toggle favorite:', err);
      setError(err instanceof Error ? err.message : 'お気に入りの更新に失敗しました');
      
      // エラー表示を一定時間後に消す
      setTimeout(() => {
        setError(null);
      }, 3000);
    }
  };

  // お気に入りアイテム - 確実にisFavoriteがtrueのアイテムのみを表示
  // nullチェックを追加して安全に処理
  const favoriteItems = favorites ? favorites.filter(item => item && item.isFavorite === true) : [];

  // AI分析を実行
  const analyzeWithAI = async (item: BacklogItem) => {
    try {
      // データの検証
      if (!item || !item.contentSummary || !item.projectName || !item.type || !item.createdUser) {
        console.error('アイテムのデータが不完全です:', item);
        setError('アイテムのデータが不完全なため、AI分析を実行できません');
        return;
      }

      setSelectedItem(item);
      setAiLoading(true);
      setShowAiModal(true);
      
      // 実際のAPIコール - OpenAIを使用
      const response = await fetch('/api/ai/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          itemId: item.id, 
          content: item.contentSummary,
          projectName: item.projectName,
          type: item.type,
          createdUserName: item.createdUser.name
        })
      });
      
      if (!response.ok) {
        throw new Error('AI分析に失敗しました');
      }
      
      const data = await response.json();
      console.log('AI分析レスポンス:', data);
      
      // レスポンスの検証
      if (!data || !data.analysis) {
        throw new Error('AIレスポンスが不正な形式です');
      }
      
      // APIレスポンスのnull/undefined対策
      const analysis = {
        summary: data.analysis.summary || '要約が生成されませんでした',
        keyPoints: Array.isArray(data.analysis.keyPoints) ? data.analysis.keyPoints : [],
        nextActions: Array.isArray(data.analysis.nextActions) ? data.analysis.nextActions : []
      };
      
      // 実際のレスポンスを使用
      setAiAnalysis(analysis);
    } catch (err) {
      console.error('AI分析中にエラーが発生しました:', err);
      setError('AI分析に失敗しました');
      setShowAiModal(false);
    } finally {
      setAiLoading(false);
    }
  };
  
  // AIモーダルを閉じる
  const closeAiModal = () => {
    setShowAiModal(false);
    setSelectedItem(null);
    setAiAnalysis(null);
  };

  return (
    <div className="home-page">
      <div className="container">
        {/* <h1>Backlog 更新情報検索</h1> */}
        
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
        
        { (
            <div>
              <p><strong>お気に入り数:</strong> {favorites.length}<strong>   全アイテム数:</strong> {items.length}</p>
            </div>
        )}
        
        {loading && <div className="loading">Loading...</div>}
        
        {error && <div className="error">{error}</div>}
        
        {!loading && !error && (
          <div className="content-sections">
            {/* お気に入りセクション */}
            <div className="section favorites-section">
              <h2 className="favorites-section-title">お気に入りのリスト</h2>
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
                      <th>AI分析</th>
                    </tr>
                  </thead>
                  <tbody>
                    {favoriteItems.map((item) => (
                      <tr key={item.id}>
                        <td className="favorite-col">
                          <input
                            type="checkbox"
                            checked={item.isFavorite}
                            onChange={() => !item.isUpdating && toggleFavorite(item.id, item.isFavorite)}
                            className={`favorite-checkbox ${item.isUpdating ? 'updating' : ''}`}
                            disabled={item.isUpdating}
                          />
                        </td>
                        <td>{item.id}</td>
                        <td>{item.projectName}</td>
                        <td>{item.type}</td>
                        <td>{item.contentSummary}</td>
                        <td>{item.createdUser.name}</td>
                        <td>{new Date(item.created).toLocaleString()}</td>
                        <td>
                          <button 
                            className="ai-button" 
                            onClick={() => analyzeWithAI(item)}
                          >
                            AI分析
                          </button>
                        </td>
                      </tr>
                    ))}
                    {favoriteItems.length === 0 && (
                      <tr>
                        <td colSpan={8} className="no-data">お気に入りに追加された項目がありません</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
            
            {/* 全ての更新情報セクション */}
            <div className="section all-items-section">
              <h2>全ての更新情報</h2>
              <div className="results">
                <table className="items-table">
                  <thead>
                    <tr>
                      <th className="favorite-col">お気に入り</th>
                      <th>ID</th>
                      <th className="favorite-col">プロジェクト</th>
                      <th>種別</th>
                      <th>内容</th>
                      <th>作成者</th>
                      <th>作成日</th>
                      <th>AI分析</th>
                    </tr>
                  </thead>
                  <tbody>
                    {items.map((item) => (
                      <tr key={item.id}>
                        <td className="favorite-col">
                          <input
                            type="checkbox"
                            checked={item.isFavorite}
                            onChange={() => !item.isUpdating && toggleFavorite(item.id, item.isFavorite)}
                            className={`favorite-checkbox ${item.isUpdating ? 'updating' : ''}`}
                            disabled={item.isUpdating}
                          />
                        </td>
                        <td>{item.id}</td>
                        <td>{item.projectName}</td>
                        <td>{item.type}</td>
                        <td>{item.contentSummary}</td>
                        <td>{item.createdUser.name}</td>
                        <td>{new Date(item.created).toLocaleString()}</td>
                        <td>
                          <button 
                            className="ai-button" 
                            onClick={() => analyzeWithAI(item)}
                          >
                            AI分析
                          </button>
                        </td>
                      </tr>
                    ))}
                    {items.length === 0 && (
                      <tr>
                        <td colSpan={8} className="no-data">表示する項目がありません</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}
        
        {/* AI分析モーダル */}
        {showAiModal && (
          <div className="ai-modal-overlay">
            <div className="ai-modal">
              <button className="close-button" onClick={closeAiModal}>✕</button>
              <h3>AI分析結果</h3>
              
              {aiLoading && (
                <div className="ai-loading">
                  <p>AIが分析中です...</p>
                </div>
              )}
              
              {!aiLoading && selectedItem && aiAnalysis && (
                <div className="ai-content">
                  <div className="ai-item-info">
                    <p><strong>項目:</strong> {selectedItem.contentSummary}</p>
                    <p><strong>プロジェクト:</strong> {selectedItem.projectName}</p>
                    <p><strong>作成者:</strong> {selectedItem.createdUser.name}</p>
                  </div>
                  
                  <div className="ai-analysis">
                    <div className="ai-section">
                      <h4>要約</h4>
                      <p>{aiAnalysis.summary || "要約が生成されませんでした"}</p>
                    </div>
                    
                    <div className="ai-section">
                      <h4>重要ポイント</h4>
                      <ul>
                        {aiAnalysis.keyPoints && aiAnalysis.keyPoints.length > 0 ? (
                          aiAnalysis.keyPoints.map((point, index) => (
                            <li key={index}>{point}</li>
                          ))
                        ) : (
                          <li>重要ポイントが見つかりませんでした</li>
                        )}
                      </ul>
                    </div>
                    
                    <div className="ai-section">
                      <h4>推奨アクション</h4>
                      <ul>
                        {aiAnalysis.nextActions && aiAnalysis.nextActions.length > 0 ? (
                          aiAnalysis.nextActions.map((action, index) => (
                            <li key={index}>{action}</li>
                          ))
                        ) : (
                          <li>推奨アクションが見つかりませんでした</li>
                        )}
                      </ul>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default HomePage; 