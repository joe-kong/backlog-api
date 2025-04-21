package backlog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

// BacklogClient はBacklog APIクライアント
type BacklogClient struct {
	spaceURL     string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

// NewBacklogClient はBacklogClientのインスタンスを生成
func NewBacklogClient(spaceURL, clientID, clientSecret string) *BacklogClient {
	return &BacklogClient{
		spaceURL:     spaceURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetActivities はBacklogのアクティビティ（更新情報）を取得
func (c *BacklogClient) GetActivities(token string, count int) ([]*model.BacklogItem, error) {
	apiURL := fmt.Sprintf("%s/api/v2/space/activities", c.spaceURL)

	// クエリパラメータの設定
	params := url.Values{}
	params.Add("count", fmt.Sprintf("%d", count))

	// リクエスト作成
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// 認証ヘッダーの設定
	req.Header.Add("Authorization", "Bearer "+token)

	// リクエスト実行
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get activities, status: %d, response: %s", resp.StatusCode, string(body))
	}

	// レスポンスボディの読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// JSONデコード
	var activities []struct {
		ID      int `json:"id"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Type    int `json:"type"`
		Content struct {
			Summary string `json:"summary"`
		} `json:"content"`
		CreatedUser struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			RoleType    int    `json:"roleType"`
			Lang        string `json:"lang"`
			MailAddress string `json:"mailAddress"`
		} `json:"createdUser"`
		Created string `json:"created"`
	}

	if err := json.Unmarshal(body, &activities); err != nil {
		return nil, err
	}

	// ドメインモデルへの変換
	items := make([]*model.BacklogItem, 0, len(activities))
	for _, activity := range activities {
		createdTime, _ := time.Parse(time.RFC3339, activity.Created)

		// typeの値を文字列に変換
		typeStr := convertTypeToString(activity.Type)

		item := &model.BacklogItem{
			ID:             fmt.Sprintf("%d", activity.ID),
			ProjectID:      fmt.Sprintf("%d", activity.Project.ID),
			ProjectName:    activity.Project.Name,
			Type:           typeStr,
			ContentSummary: activity.Content.Summary,
			CreatedUser: model.User{
				ID:          fmt.Sprintf("%d", activity.CreatedUser.ID),
				Name:        activity.CreatedUser.Name,
				RoleType:    activity.CreatedUser.RoleType,
				Lang:        activity.CreatedUser.Lang,
				MailAddress: activity.CreatedUser.MailAddress,
			},
			Created: createdTime,
		}

		items = append(items, item)
	}

	return items, nil
}

// SearchActivities はキーワードでアクティビティを検索
func (c *BacklogClient) SearchActivities(token, keyword string, count int) ([]*model.BacklogItem, error) {
	// 全てのアクティビティを取得
	activities, err := c.GetActivities(token, count)
	if err != nil {
		return nil, err
	}

	// キーワードが空の場合は全て返す
	if keyword == "" {
		return activities, nil
	}

	// キーワードでフィルタリング
	var filtered []*model.BacklogItem
	for _, activity := range activities {
		if strings.Contains(activity.ID, keyword) ||
			strings.Contains(activity.ProjectName, keyword) ||
			strings.Contains(activity.Type, keyword) ||
			strings.Contains(activity.ContentSummary, keyword) ||
			strings.Contains(activity.CreatedUser.Name, keyword) {
			filtered = append(filtered, activity)
		}
	}

	return filtered, nil
}

// convertTypeToString はBacklog APIのtype値を文字列に変換
func convertTypeToString(typeCode int) string {
	// Backlog API公式ドキュメントに基づく種別コードと文字列のマッピング
	switch typeCode {
	case 1:
		return "課題の追加"
	case 2:
		return "課題の更新"
	case 3:
		return "課題にコメント"
	case 4:
		return "課題の削除"
	case 5:
		return "Wikiを追加"
	case 6:
		return "Wikiを更新"
	case 7:
		return "Wikiを削除"
	case 8:
		return "共有ファイルを追加"
	case 9:
		return "共有ファイルを更新"
	case 10:
		return "共有ファイルを削除"
	case 11:
		return "Subversionコミット"
	case 12:
		return "GITプッシュ"
	case 13:
		return "GITリポジトリ作成"
	case 14:
		return "課題をまとめて更新"
	case 15:
		return "ユーザーがプロジェクトに参加"
	case 16:
		return "ユーザーがプロジェクトから脱退"
	case 17:
		return "コメントにお知らせを追加"
	case 18:
		return "プルリクエストの追加"
	case 19:
		return "プルリクエストの更新"
	case 20:
		return "プルリクエストにコメント"
	case 21:
		return "プルリクエストの削除"
	case 22:
		return "マイルストーンの追加"
	case 23:
		return "マイルストーンの更新"
	case 24:
		return "マイルストーンの削除"
	case 25:
		return "グループがプロジェクトに参加"
	case 26:
		return "グループがプロジェクトから脱退"
	default:
		return fmt.Sprintf("種別(%d)", typeCode)
	}
}
