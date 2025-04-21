package backlog

import (
	"log"
	"time"

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
	// アクセストークンの取得（認証リポジトリからランダムなユーザーのトークンを取得）
	tokens, err := s.authRepository.GetAllTokens()
	if err != nil || len(tokens) == 0 {
		// トークンがない場合はモックデータで代用
		return s.mockBacklogItems(), nil
	}

	// 最初のトークンを使用
	token := tokens[0]

	// Backlog APIを呼び出してアクティビティを検索
	items, err := s.client.SearchActivities(token.AccessToken, keyword, 100)
	log.Println("SearchActivities items:", items)
	if err != nil {
		log.Println("SearchActivities err:", err)
		// APIエラーの場合はモックデータで代用
		return s.mockBacklogItems(), nil
	}

	return items, nil
}

// GetFavorites はユーザーのお気に入りBacklog更新情報を取得
func (s *BacklogItemService) GetFavorites(userID string) ([]*model.BacklogItem, error) {
	// ユーザーIDからトークンを取得
	token, err := s.authRepository.GetToken(userID)
	if err != nil {
		// トークンがない場合はモックデータで代用
		items := s.mockBacklogItems()
		return items[:2], nil
	}

	// Backlog APIを呼び出して全アクティビティを取得
	items, err := s.client.GetActivities(token.AccessToken, 50)
	log.Println("GetActivities items:", items)
	if err != nil {
		log.Println("GetActivities err:", err)
		// APIエラーの場合はモックデータで代用
		mockItems := s.mockBacklogItems()
		return mockItems[:2], nil
	}

	// if len(items) > 3 {
	// 	return items[:3], nil
	// }
	return items, nil
}

// AddFavorite はBacklog更新情報をお気に入りに追加
func (s *BacklogItemService) AddFavorite(userID string, itemID string) error {
	// お気に入り追加の実装
	// データベースに保存する
	return nil
}

// RemoveFavorite はBacklog更新情報をお気に入りから削除
func (s *BacklogItemService) RemoveFavorite(userID string, itemID string) error {
	// お気に入り削除の実装
	// データベースから削除する
	return nil
}

// mockBacklogItems はモックのBacklog更新情報を生成
func (s *BacklogItemService) mockBacklogItems() []*model.BacklogItem {
	// 現在時刻を基準に日付を設定
	now := time.Now()

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
			Created: now.AddDate(0, 0, -5), // 5日前
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
			Created: now.AddDate(0, 0, -3), // 3日前
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
			Created: now.AddDate(0, 0, -2), // 2日前
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
			Created: now.AddDate(0, 0, -1), // 1日前
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
			Created: now, // 現在
		},
	}
	return items
}
