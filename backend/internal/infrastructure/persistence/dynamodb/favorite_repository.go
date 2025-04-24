package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
)

const (
	// FavoriteTableName はDynamoDBのテーブル名
	FavoriteTableName = "Favorites"
	// IndexNameUserID はユーザーIDによる検索用のグローバルセカンダリインデックス名
	IndexNameUserID = "UserID-index"
)

// FavoriteItem はDynamoDBに保存するためのお気に入りアイテム構造体
type FavoriteItem struct {
	ID        string    `dynamodbav:"id"`
	UserID    string    `dynamodbav:"userId"`
	ItemID    string    `dynamodbav:"itemId"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
}

// FavoriteRepository はDynamoDBを使ったお気に入りリポジトリの実装
type FavoriteRepository struct {
	client *dynamodb.Client
}

// NewFavoriteRepository はFavoriteRepositoryのインスタンスを生成
func NewFavoriteRepository(client *dynamodb.Client) *FavoriteRepository {
	return &FavoriteRepository{
		client: client,
	}
}

// FindByUserID はユーザーIDからお気に入りを検索
func (r *FavoriteRepository) FindByUserID(userID string) ([]*model.Favorite, error) {
	// GSIを使用してユーザーIDでクエリ
	input := &dynamodb.QueryInput{
		TableName:              aws.String(FavoriteTableName),
		IndexName:              aws.String(IndexNameUserID),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userID},
		},
	}

	output, err := r.client.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorites: %w", err)
	}

	if output.Count == 0 {
		return []*model.Favorite{}, nil
	}

	// DynamoDB項目をドメインモデルに変換
	var favoriteItems []FavoriteItem
	err = attributevalue.UnmarshalListOfMaps(output.Items, &favoriteItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal favorites: %w", err)
	}

	// ドメインモデルに変換
	favorites := make([]*model.Favorite, len(favoriteItems))
	for i, item := range favoriteItems {
		favorites[i] = &model.Favorite{
			ID:        item.ID,
			UserID:    item.UserID,
			ItemID:    item.ItemID,
			CreatedAt: item.CreatedAt,
		}
	}

	return favorites, nil
}

// Save はお気に入りを保存
func (r *FavoriteRepository) Save(favorite *model.Favorite) error {
	// ドメインモデルをDynamoDB項目に変換
	item := FavoriteItem{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		ItemID:    favorite.ItemID,
		CreatedAt: favorite.CreatedAt,
	}

	// 項目をマーシャリング
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal favorite: %w", err)
	}

	// DynamoDBに保存
	input := &dynamodb.PutItemInput{
		TableName: aws.String(FavoriteTableName),
		Item:      av,
	}

	_, err = r.client.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to save favorite: %w", err)
	}

	return nil
}

// Delete はお気に入りを削除
func (r *FavoriteRepository) Delete(userID string, itemID string) error {
	// 最初に一致するアイテムを探してIDを取得
	existingItems, err := r.findByUserIDAndItemID(userID, itemID)
	if err != nil {
		return fmt.Errorf("failed to find favorite for deletion: %w", err)
	}

	if len(existingItems) == 0 {
		// 削除するお気に入りが見つからない場合はエラーを返さずに成功とする
		return nil
	}

	// 一致する最初のアイテムを削除（通常は1つだけのはず）
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(FavoriteTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: existingItems[0].ID},
		},
	}

	_, err = r.client.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete favorite: %w", err)
	}

	return nil
}

// Exists はお気に入りが存在するかチェック
func (r *FavoriteRepository) Exists(userID string, itemID string) (bool, error) {
	existingItems, err := r.findByUserIDAndItemID(userID, itemID)
	if err != nil {
		return false, err
	}

	return len(existingItems) > 0, nil
}

// findByUserIDAndItemID はユーザーIDとアイテムIDの両方に一致するお気に入りを検索する内部メソッド
func (r *FavoriteRepository) findByUserIDAndItemID(userID string, itemID string) ([]*model.Favorite, error) {
	// ユーザーIDでクエリを実行
	input := &dynamodb.QueryInput{
		TableName:              aws.String(FavoriteTableName),
		IndexName:              aws.String(IndexNameUserID),
		KeyConditionExpression: aws.String("userId = :userId"),
		FilterExpression:       aws.String("itemId = :itemId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userID},
			":itemId": &types.AttributeValueMemberS{Value: itemID},
		},
	}

	// クエリ実行
	output, err := r.client.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorites: %w", err)
	}

	// 結果がない場合は空のスライスを返す
	if output.Count == 0 {
		return []*model.Favorite{}, nil
	}

	// DynamoDB項目をドメインモデルに変換
	var favoriteItems []FavoriteItem
	err = attributevalue.UnmarshalListOfMaps(output.Items, &favoriteItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal favorites: %w", err)
	}

	// ドメインモデルに変換
	favorites := make([]*model.Favorite, len(favoriteItems))
	for i, item := range favoriteItems {
		favorites[i] = &model.Favorite{
			ID:        item.ID,
			UserID:    item.UserID,
			ItemID:    item.ItemID,
			CreatedAt: item.CreatedAt,
		}
	}

	return favorites, nil
}
