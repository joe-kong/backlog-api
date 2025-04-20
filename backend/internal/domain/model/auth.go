package model

import (
	"time"
)

// AuthToken は認証トークン情報を表すドメインモデル
type AuthToken struct {
	AccessToken  string    `json:"accessToken"`
	TokenType    string    `json:"tokenType"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	UserID       string    `json:"userId"`
}

// OAuthConfig はOAuth2.0認証の設定情報
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
	Scopes       []string
}

// AuthService は認証に関するドメインサービスのインターフェース
type AuthService interface {
	GetAuthorizationURL() string
	ExchangeCodeForToken(code string) (*AuthToken, error)
	RefreshToken(refreshToken string) (*AuthToken, error)
	GetBacklogUser(accessToken string) (*User, error)
}

// AuthRepository は認証情報の永続化を担当するリポジトリのインターフェース
type AuthRepository interface {
	SaveToken(token *AuthToken) error
	GetTokenByUserID(userID string) (*AuthToken, error)
	DeleteToken(userID string) error
}
