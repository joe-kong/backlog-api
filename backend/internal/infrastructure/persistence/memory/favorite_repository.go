package memory

import (
	"sync"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// FavoriteRepository はインメモリお気に入りリポジトリの実装
type FavoriteRepository struct {
	favorites []*model.Favorite
	mu        sync.RWMutex
}

// NewFavoriteRepository はFavoriteRepositoryのインスタンスを生成
func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{
		favorites: make([]*model.Favorite, 0),
	}
}

// FindByUserID はユーザーIDからお気に入りを検索
func (r *FavoriteRepository) FindByUserID(userID string) ([]*model.Favorite, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Favorite, 0)
	for _, fav := range r.favorites {
		if fav.UserID == userID {
			result = append(result, fav)
		}
	}

	return result, nil
}

// Save はお気に入りを保存
func (r *FavoriteRepository) Save(favorite *model.Favorite) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.favorites = append(r.favorites, favorite)
	return nil
}

// Delete はお気に入りを削除
func (r *FavoriteRepository) Delete(userID string, itemID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := make([]*model.Favorite, 0)
	for _, fav := range r.favorites {
		if !(fav.UserID == userID && fav.ItemID == itemID) {
			filtered = append(filtered, fav)
		}
	}

	r.favorites = filtered
	return nil
}

// Exists はお気に入りが存在するかチェック
func (r *FavoriteRepository) Exists(userID string, itemID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, fav := range r.favorites {
		if fav.UserID == userID && fav.ItemID == itemID {
			return true, nil
		}
	}

	return false, nil
}
