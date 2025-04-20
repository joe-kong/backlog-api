package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// BacklogAuthService はBacklogのOAuth認証サービス実装
type BacklogAuthService struct {
	oauthConfig *oauth2.Config
	spaceURL    string
}

// NewBacklogAuthService はBacklogAuthServiceのインスタンスを生成
func NewBacklogAuthService(config model.OAuthConfig, spaceURL string) *BacklogAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		Scopes: config.Scopes,
	}

	return &BacklogAuthService{
		oauthConfig: oauthConfig,
		spaceURL:    spaceURL,
	}
}

// GetAuthorizationURL は認可URLを取得
func (s *BacklogAuthService) GetAuthorizationURL() string {
	return s.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken は認可コードからトークンを取得
func (s *BacklogAuthService) ExchangeCodeForToken(code string) (*model.AuthToken, error) {
	ctx := context.Background()
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// ドメインモデルに変換
	return &model.AuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

// RefreshToken はリフレッシュトークンから新しいトークンを取得
func (s *BacklogAuthService) RefreshToken(refreshToken string) (*model.AuthToken, error) {
	ctx := context.Background()

	// リフレッシュトークンから新しいトークンを取得
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	source := s.oauthConfig.TokenSource(ctx, token)
	newToken, err := source.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// ドメインモデルに変換
	return &model.AuthToken{
		AccessToken:  newToken.AccessToken,
		TokenType:    newToken.TokenType,
		RefreshToken: newToken.RefreshToken,
		ExpiresAt:    newToken.Expiry,
	}, nil
}

// GetBacklogUser はアクセストークンを使用してBacklogユーザー情報を取得
func (s *BacklogAuthService) GetBacklogUser(accessToken string) (*model.User, error) {
	url := fmt.Sprintf("%s/api/v2/users/myself", s.spaceURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userResp struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		RoleType    int    `json:"roleType"`
		Lang        string `json:"lang"`
		MailAddress string `json:"mailAddress"`
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, err
	}

	// ドメインモデルに変換
	return &model.User{
		ID:          fmt.Sprintf("%d", userResp.ID),
		Name:        userResp.Name,
		RoleType:    userResp.RoleType,
		Lang:        userResp.Lang,
		MailAddress: userResp.MailAddress,
	}, nil
}
