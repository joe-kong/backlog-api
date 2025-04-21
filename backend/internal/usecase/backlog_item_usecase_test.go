package usecase_test

import (
	"testing"
	"time"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/persistence/memory"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/usecase"
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
				ProjectName:    "プロジェクトA",
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

// MockAuthUseCase はAuthUseCaseのモック実装
type MockAuthUseCase struct{}

func (m *MockAuthUseCase) GetValidToken(userID string) (*model.AuthToken, error) {
	return &model.AuthToken{
		AccessToken:  "mock-token",
		TokenType:    "Bearer",
		RefreshToken: "mock-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
		UserID:       userID,
	}, nil
}

// TestAuthUseCase は認証ユースケースのテスト用インターフェース
type TestAuthUseCase interface {
	GetValidToken(userID string) (*model.AuthToken, error)
}

func TestBacklogItemUseCase_SearchItems(t *testing.T) {
	// テスト用のリポジトリとサービスを初期化
	mockBacklogService := NewMockBacklogItemService()
	favoriteRepo := memory.NewFavoriteRepository()
	mockAuthUseCase := &MockAuthUseCase{}

	// テスト対象のユースケースを初期化
	usecase := usecase.NewBacklogItemUseCase(mockBacklogService, favoriteRepo, mockAuthUseCase)

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
			results, err := usecase.SearchItems(tc.userID, tc.keyword)
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
	mockAuthUseCase := &MockAuthUseCase{}

	// テスト対象のユースケースを初期化
	usecase := usecase.NewBacklogItemUseCase(mockBacklogService, favoriteRepo, mockAuthUseCase)

	// テスト実行
	userID := "user1"
	itemID := "1"

	// お気に入り追加
	err := usecase.AddFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to add favorite: %v", err)
	}

	// 同じアイテムを再度追加すると重複エラーが発生するはず
	err = usecase.AddFavorite(userID, itemID)
	if err == nil {
		t.Error("Expected duplicate error, but got nil")
	}

	// お気に入り削除
	err = usecase.RemoveFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to remove favorite: %v", err)
	}

	// 再度お気に入り追加ができることを確認
	err = usecase.AddFavorite(userID, itemID)
	if err != nil {
		t.Fatalf("Failed to add favorite after removing: %v", err)
	}
}
