package backlog

import (
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// BacklogItemService はBacklogItemServiceのインフラ層実装
type BacklogItemService struct {
	client         *BacklogClient
	authRepository model.AuthRepository
}

// NewBacklogItemService はBacklogItemServiceのインスタンスを生成
func NewBacklogItemService(client *BacklogClient, authRepository model.AuthRepository) *BacklogItemService {
	return &BacklogItemService{
		client:         client,
		authRepository: authRepository,
	}
}

// SearchItems はキーワードでBacklog更新情報を検索
func (s *BacklogItemService) SearchItems(keyword string) ([]*model.BacklogItem, error) {
	// この実装ではユーザーIDは使用しないため、モックデータを返す
	// 実際のアプリケーションではユーザーIDからトークンを取得し、Backlog APIを呼び出す
	if keyword == "" {
		return s.mockBacklogItems(), nil
	}

	// キーワードでフィルタリング
	var filtered []*model.BacklogItem
	for _, item := range s.mockBacklogItems() {
		if contains(item.ID, keyword) ||
			contains(item.ProjectName, keyword) ||
			contains(item.Type, keyword) ||
			contains(item.ContentSummary, keyword) ||
			contains(item.CreatedUser.Name, keyword) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

// GetFavorites はユーザーのお気に入りBacklog更新情報を取得
func (s *BacklogItemService) GetFavorites(userID string) ([]*model.BacklogItem, error) {
	// 開発環境ではモックデータを返す
	// モックのお気に入り情報を生成
	items := s.mockBacklogItems()
	// 最初の2つをお気に入りとする
	return items[:2], nil
}

// AddFavorite はBacklog更新情報をお気に入りに追加
func (s *BacklogItemService) AddFavorite(userID string, itemID string) error {
	// お気に入り追加の実装
	// 実際のアプリケーションではデータベースに保存する
	return nil
}

// RemoveFavorite はBacklog更新情報をお気に入りから削除
func (s *BacklogItemService) RemoveFavorite(userID string, itemID string) error {
	// お気に入り削除の実装
	// 実際のアプリケーションではデータベースから削除する
	return nil
}

// mockBacklogItems はモックのBacklog更新情報を生成
func (s *BacklogItemService) mockBacklogItems() []*model.BacklogItem {
	// 日付形式を修正
	items := []*model.BacklogItem{
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
		},
		{
			ID:             "3",
			ProjectID:      "2",
			ProjectName:    "プロジェクトB",
			Type:           "Wiki",
			ContentSummary: "設計ドキュメント",
			CreatedUser: model.User{
				ID:          "3",
				Name:        "鈴木一郎",
				RoleType:    1,
				Lang:        "ja",
				MailAddress: "suzuki@example.com",
			},
		},
		{
			ID:             "4",
			ProjectID:      "2",
			ProjectName:    "プロジェクトB",
			Type:           "Git",
			ContentSummary: "バグ修正のコミット",
			CreatedUser: model.User{
				ID:          "4",
				Name:        "田中次郎",
				RoleType:    1,
				Lang:        "ja",
				MailAddress: "tanaka@example.com",
			},
		},
		{
			ID:             "5",
			ProjectID:      "3",
			ProjectName:    "プロジェクトC",
			Type:           "課題",
			ContentSummary: "UIデザインの改善",
			CreatedUser: model.User{
				ID:          "5",
				Name:        "高橋三郎",
				RoleType:    1,
				Lang:        "ja",
				MailAddress: "takahashi@example.com",
			},
		},
	}
	return items
}

// contains は文字列に部分文字列が含まれているかチェック
func contains(str, substr string) bool {
	return str != "" && substr != "" && str == substr
}
