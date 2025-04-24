package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// BacklogItemUseCase はBacklog更新情報に関するユースケース
type BacklogItemUseCase struct {
	backlogItemService model.BacklogItemService
	favoriteRepository model.FavoriteRepository
	authUseCase        *AuthUseCase
}

// BacklogItemOutput はBacklogItemの出力用データ
type BacklogItemOutput struct {
	ID             string `json:"id"`
	ProjectID      string `json:"projectId"`
	ProjectName    string `json:"projectName"`
	Type           string `json:"type"`
	ContentSummary string `json:"contentSummary"`
	CreatedUser    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createdUser"`
	Created    time.Time `json:"created"`
	IsFavorite bool      `json:"isFavorite"`
}

// NewBacklogItemUseCase はBacklogItemUseCaseのインスタンスを生成
func NewBacklogItemUseCase(
	backlogItemService model.BacklogItemService,
	favoriteRepository model.FavoriteRepository,
	authUseCase *AuthUseCase,
) *BacklogItemUseCase {
	return &BacklogItemUseCase{
		backlogItemService: backlogItemService,
		favoriteRepository: favoriteRepository,
		authUseCase:        authUseCase,
	}
}

// SearchItems はキーワードでBacklog更新情報を検索
func (u *BacklogItemUseCase) SearchItems(userID, keyword string) ([]*BacklogItemOutput, error) {
	// ユーザーのアクセストークンを取得
	_, err := u.authUseCase.GetValidToken(userID)
	if err != nil {
		return nil, err
	}

	// 更新情報を検索
	items, err := u.backlogItemService.SearchItems(keyword)
	if err != nil {
		return nil, err
	}

	// ユーザーのお気に入り情報を取得
	favorites, err := u.favoriteRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// お気に入りマップを作成
	favoriteMap := make(map[string]bool)
	for _, fav := range favorites {
		favoriteMap[fav.ItemID] = true
	}

	// 出力データを作成
	outputs := make([]*BacklogItemOutput, len(items))

	for i, item := range items {
		output := &BacklogItemOutput{
			ID:             item.ID,
			ProjectID:      item.ProjectID,
			ProjectName:    item.ProjectName,
			Type:           item.Type,
			ContentSummary: item.ContentSummary,
			Created:        item.Created,
			IsFavorite:     favoriteMap[item.ID],
		}
		output.CreatedUser.ID = item.CreatedUser.ID
		output.CreatedUser.Name = item.CreatedUser.Name
		outputs[i] = output
	}

	return outputs, nil
}

// GetFavorites はユーザーのお気に入り情報を取得
func (u *BacklogItemUseCase) GetFavorites(userID string) ([]*BacklogItemOutput, error) {
	_, err := u.authUseCase.GetValidToken(userID)
	if err != nil {
		return nil, err
	}

	favorites, err := u.favoriteRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(favorites) == 0 {
		return []*BacklogItemOutput{}, nil
	}

	favoriteIDs := make(map[string]bool)
	for _, fav := range favorites {
		favoriteIDs[fav.ItemID] = true
	}

	// すべてのアイテムを取得
	allItems, err := u.backlogItemService.SearchItems("")
	if err != nil {
		return nil, err
	}

	// お気に入りに登録されているアイテムだけをフィルタリング
	var favoriteItems []*model.BacklogItem
	for _, item := range allItems {
		if favoriteIDs[item.ID] {
			favoriteItems = append(favoriteItems, item)
		}
	}

	// 出力データを作成
	outputs := make([]*BacklogItemOutput, len(favoriteItems))
	for i, item := range favoriteItems {
		output := &BacklogItemOutput{
			ID:             item.ID,
			ProjectID:      item.ProjectID,
			ProjectName:    item.ProjectName,
			Type:           item.Type,
			ContentSummary: item.ContentSummary,
			Created:        item.Created,
			IsFavorite:     true,
		}
		output.CreatedUser.ID = item.CreatedUser.ID
		output.CreatedUser.Name = item.CreatedUser.Name
		outputs[i] = output
	}

	return outputs, nil
}

// AddFavorite はお気に入りを追加
func (u *BacklogItemUseCase) AddFavorite(userID, itemID string) error {
	// 既に存在するかチェック
	exists, err := u.favoriteRepository.Exists(userID, itemID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("favorite already exists")
	}

	// 新しいお気に入りを作成
	favorite := &model.Favorite{
		ID:        uuid.New().String(),
		UserID:    userID,
		ItemID:    itemID,
		CreatedAt: time.Now(),
	}

	// お気に入りを保存
	return u.favoriteRepository.Save(favorite)
}

// RemoveFavorite はお気に入りを削除
func (u *BacklogItemUseCase) RemoveFavorite(userID, itemID string) error {
	return u.favoriteRepository.Delete(userID, itemID)
}
