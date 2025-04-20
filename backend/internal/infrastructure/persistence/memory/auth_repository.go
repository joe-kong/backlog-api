package memory

import (
	"errors"
	"sync"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// AuthRepository はインメモリ認証リポジトリの実装
type AuthRepository struct {
	tokens map[string]*model.AuthToken
	mu     sync.RWMutex
}

// NewAuthRepository はAuthRepositoryのインスタンスを生成
func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		tokens: make(map[string]*model.AuthToken),
	}
}

// SaveToken はトークンを保存
func (r *AuthRepository) SaveToken(token *model.AuthToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[token.UserID] = token
	return nil
}

// GetTokenByUserID はユーザーIDからトークンを取得
func (r *AuthRepository) GetTokenByUserID(userID string) (*model.AuthToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, exists := r.tokens[userID]
	if !exists {
		return nil, errors.New("token not found")
	}

	return token, nil
}

// DeleteToken はトークンを削除
func (r *AuthRepository) DeleteToken(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tokens, userID)
	return nil
}
