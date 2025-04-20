package model

import (
	"time"
)

// Favorite はユーザーのお気に入り情報を表すドメインモデル
type Favorite struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ItemID    string    `json:"itemId"`
	CreatedAt time.Time `json:"createdAt"`
}

// FavoriteRepository はお気に入り情報の永続化を担当するリポジトリのインターフェース
type FavoriteRepository interface {
	FindByUserID(userID string) ([]*Favorite, error)
	Save(favorite *Favorite) error
	Delete(userID string, itemID string) error
	Exists(userID string, itemID string) (bool, error)
}
