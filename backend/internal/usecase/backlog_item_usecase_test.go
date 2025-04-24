package usecase

import (
	"testing"
	"time"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/persistence/memory"
)

// MockBacklogItemService はBacklogItemServiceのモック実装
type MockBacklogItemService struct {
	items []*model.BacklogItem
}

func NewMockBacklogItemService() *MockBacklogItemService {
	return &MockBacklogItemService{
		items: []*model.BacklogItem{
			{
				ID:             "1",
				ProjectID:      "1",
				ProjectName:    "プロジェクトA",
				Type:           "課題",
				ContentSummary: "ログイン機能の実装",
				CreatedUser: model.User{
					ID:          "1",
					Name:        "山田太郎",
					RoleType:    1,
					Lang:        "ja",
					MailAddress: "yamada@example.com",
				},
				Created: time.Now(),
			},
			{
				ID:             "2",
				ProjectID:      "1",
				ProjectName:    "プロジェクトB",
				Type:           "課題",
				ContentSummary: "検索機能の追加",
				CreatedUser: model.User{
					ID:          "2",
					Name:        "佐藤花子",
					RoleType:    1,
					Lang:        "ja",
					MailAddress: "sato@example.com",
				},
				Created: time.Now(),
			},
		},
	}
}

func (m *MockBacklogItemService) SearchItems(keyword string) ([]*model.BacklogItem, error) {
	if keyword == "" {
		return m.items, nil
	}

	var result []*model.BacklogItem
	for _, item := range m.items {
		if item.ID == keyword || item.ProjectName == keyword || item.Type == keyword ||
			item.ContentSummary == keyword || item.CreatedUser.Name == keyword {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockBacklogItemService) GetFavorites(userID string) ([]*model.BacklogItem, error) {
	return m.items[:1], nil
}

func (m *MockBacklogItemService) AddFavorite(userID string, itemID string) error {
	return nil
}

func (m *MockBacklogItemService) RemoveFavorite(userID string, itemID string) error {
	return nil
}

// MockAuthRepository はAuthRepositoryのモック実装
type MockAuthRepository struct{}

func (m *MockAuthRepository) SaveToken(token *model.AuthToken) error {
	return nil
}

func (m *MockAuthRepository) GetTokenByUserID(userID string) (*model.AuthToken, error) {
	return &model.AuthToken{
		AccessToken:  "test-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
		UserID:       userID,
	}, nil
}

func (m *MockAuthRepository) DeleteToken(userID string) error {
	return nil
}

func (m *MockAuthRepository) GetAllTokens() ([]*model.AuthToken, error) {
	return []*model.AuthToken{
		{
			AccessToken:  "test-token-1",
			TokenType:    "Bearer",
			RefreshToken: "test-refresh-token-1",
			ExpiresAt:    time.Now().Add(time.Hour),
			UserID:       "user1",
		},
	}, nil
}

func (m *MockAuthRepository) GetToken(userID string) (*model.AuthToken, error) {
	return &model.AuthToken{
		AccessToken:  "test-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
		UserID:       userID,
	}, nil
}

// MockAuthService はAuthServiceのモック実装
type MockAuthService struct{}

func (m *MockAuthService) GetAuthorizationURL() string {
	return "https://test.com/auth"
}

// 認証認可用ダミートークン取得
func (m *MockAuthService) ExchangeCodeForToken(code string) (*model.AuthToken, error) {
	return &model.AuthToken{
		AccessToken:  "test-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}, nil
}

// テスト用ダミーリフレッシュトークン更新
func (m *MockAuthService) RefreshToken(refreshToken string) (*model.AuthToken, error) {
	return &model.AuthToken{
		AccessToken:  "refreshed-token",
		TokenType:    "Bearer",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}, nil
}

// テスト用のダミーユーザーを取得
func (m *MockAuthService) GetBacklogUser(accessToken string) (*model.User, error) {
	return &model.User{
		ID:          "test-user",
		Name:        "Test User",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example.com",
	}, nil
}

// テスト用のAuthUseCaseを作成
func createTestAuthUseCase() *AuthUseCase {
	return NewAuthUseCase(
		&MockAuthService{},
		&MockAuthRepository{},
	)
}

// データ検索メソッドをテストする
func TestBacklogItemUseCase_SearchItems(t *testing.T) {
	mockBacklogService := NewMockBacklogItemService()
	favoriteRepo := memory.NewFavoriteRepository()
	mockAuthUseCase := createTestAuthUseCase()

	// テスト対象のユースケースを初期化
	backlogUseCase := NewBacklogItemUseCase(mockBacklogService, favoriteRepo, mockAuthUseCase)

	// テスト実行
	testCases := []struct {
		name     string
		userID   string
		keyword  string
		expected int
	}{
		{
			name:     "全件取得",
			userID:   "user1",
			keyword:  "",
			expected: 2,
		},
		{
			name:     "キーワード検索：プロジェクトA",
			userID:   "user1",
			keyword:  "プロジェクトA",
			expected: 1,
		},
		{
			name:     "キーワード検索：存在しないキーワード",
			userID:   "user1",
			keyword:  "存在しないキーワード",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := backlogUseCase.SearchItems(tc.userID, tc.keyword)
			if err != nil {
				t.Fatalf("Failed to search items: %v", err)
			}

			if len(results) != tc.expected {
				t.Errorf("Expected %d items, got %d", tc.expected, len(results))
			}
		})
	}
}

func TestBacklogItemUseCase_AddFavorite(t *testing.T) {
	// テスト用のリポジトリとサービスを初期化
	mockBacklogService := NewMockBacklogItemService()
	favoriteRepo := memory.NewFavoriteRepository()
	authUseCase := createTestAuthUseCase()

	// テスト対象のユースケースを初期化
	backlogUseCase := NewBacklogItemUseCase(mockBacklogService, favoriteRepo, authUseCase)

	// テスト実行
	userID := "user1"
	itemID := "1"

	// お気に入り追加
	err := backlogUseCase.AddFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to add favorite: %v", err)
	}

	// 同じアイテムを再度追加すると重複エラーが発生するはず
	err = backlogUseCase.AddFavorite(userID, itemID)
	if err == nil {
		t.Error("Expected duplicate error, but got nil")
	}

	// お気に入り削除
	err = backlogUseCase.RemoveFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to remove favorite: %v", err)
	}

	// 再度お気に入り追加ができることを確認
	err = backlogUseCase.AddFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to add favorite after removing: %v", err)
	}
}
