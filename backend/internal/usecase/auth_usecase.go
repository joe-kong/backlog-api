package usecase

import (
	"errors"
	"time"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// ErrInvalidToken は無効なトークンエラー
var ErrInvalidToken = errors.New("invalid token")

// AuthUseCase は認証に関するユースケース
type AuthUseCase struct {
	authService    model.AuthService
	authRepository model.AuthRepository
}

// NewAuthUseCase は認証ユースケースのインスタンスを生成
func NewAuthUseCase(authService model.AuthService, authRepository model.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		authService:    authService,
		authRepository: authRepository,
	}
}

// GetAuthorizationURL は認可URLを取得
func (u *AuthUseCase) GetAuthorizationURL() string {
	return u.authService.GetAuthorizationURL()
}

// AuthorizeCallback は認可コードからトークンを取得し保存する
func (u *AuthUseCase) AuthorizeCallback(code string) (*model.AuthToken, *model.User, error) {
	// コードからトークンを取得
	token, err := u.authService.ExchangeCodeForToken(code)
	if err != nil {
		return nil, nil, err
	}

	// ユーザー情報を取得
	user, err := u.authService.GetBacklogUser(token.AccessToken)
	if err != nil {
		return nil, nil, err
	}

	// ユーザーIDをトークンに関連付け
	token.UserID = user.ID

	// トークンを保存
	err = u.authRepository.SaveToken(token)
	if err != nil {
		return nil, nil, err
	}

	return token, user, nil
}

// GetValidToken はユーザーの有効なトークンを取得
func (u *AuthUseCase) GetValidToken(userID string) (*model.AuthToken, error) {
	token, err := u.authRepository.GetTokenByUserID(userID)
	if err != nil {
		return nil, err
	}

	// トークンが有効期限切れかどうかチェック
	if token.ExpiresAt.Before(time.Now()) {
		// リフレッシュトークンを使用して新しいトークンを取得
		newToken, err := u.authService.RefreshToken(token.RefreshToken)
		if err != nil {
			return nil, err
		}

		// ユーザーIDを設定
		newToken.UserID = userID

		// 新しいトークンを保存
		err = u.authRepository.SaveToken(newToken)
		if err != nil {
			return nil, err
		}

		return newToken, nil
	}

	return token, nil
}

// Logout はユーザーのログアウト処理
func (u *AuthUseCase) Logout(userID string) error {
	return u.authRepository.DeleteToken(userID)
}
