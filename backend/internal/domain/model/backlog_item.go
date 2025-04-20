package model

import (
	"time"
)

// BacklogItem はBacklogの更新情報を表すドメインモデル
type BacklogItem struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"projectId"`
	ProjectName    string    `json:"projectName"`
	Type           string    `json:"type"`
	ContentSummary string    `json:"contentSummary"`
	CreatedUser    User      `json:"createdUser"`
	Created        time.Time `json:"created"`
}

// User はBacklogのユーザー情報を表す
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	RoleType    int    `json:"roleType"`
	Lang        string `json:"lang"`
	MailAddress string `json:"mailAddress"`
}

// BacklogItemService はBacklogItemに関するドメインサービスのインターフェース
type BacklogItemService interface {
	SearchItems(keyword string) ([]*BacklogItem, error)
	GetFavorites(userID string) ([]*BacklogItem, error)
	AddFavorite(userID string, itemID string) error
	RemoveFavorite(userID string, itemID string) error
}
