package dynamodb

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// NewDynamoDBClient はDynamoDBクライアントのインスタンスを生成
func NewDynamoDBClient(region string) (*dynamodb.Client, error) {
	// 環境変数から認証情報を取得
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN") // 必須ではない。

	var cfg aws.Config
	var err error

	// 認証情報が環境変数で提供されている場合は、それを使用
	if accessKey != "" && secretKey != "" {
		// 静的な認証情報を使用
		credProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credProvider),
		)
	} else {
		// 認証情報が提供されない場合は、デフォルトの認証情報プロバイダーチェーンを使用
		log.Println("AWS認証情報が見つからないため、デフォルトの認証情報プロバイダーチェーンを使用します。")
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
	}

	if err != nil {
		return nil, err
	}

	// DynamoDBクライアント生成
	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}

// CreateFavoriteTable はお気に入りテーブルを作成
func CreateFavoriteTable(client *dynamodb.Client) error {
	// テーブルが既に存在するか確認
	existing, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	// テーブルが既に存在する場合は作成をスキップ
	for _, tableName := range existing.TableNames {
		if tableName == FavoriteTableName {
			log.Printf("Table %s already exists\n", FavoriteTableName)
			return nil
		}
	}

	// テーブル作成リクエスト
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(FavoriteTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("userId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(IndexNameUserID),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("userId"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	// テーブル作成
	_, err = client.CreateTable(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Printf("Created table %s\n", FavoriteTableName)
	return nil
}
