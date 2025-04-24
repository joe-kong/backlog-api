package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/auth"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/backlog"
	dynamodb_repo "nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/persistence/dynamodb"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/persistence/memory"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/usecase"
)

func main() {

	// .env ファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: 環境変数ファイルが見つかっていませんでした: %v", err)
	}

	spaceURL := getEnv("BACKLOG_SPACE_URL", "")
	clientID := getEnv("BACKLOG_CLIENT_ID", "")
	clientSecret := getEnv("BACKLOG_CLIENT_SECRET", "")
	redirectURI := getEnv("OAUTH_REDIRECT_URI", "")
	authURL := getEnv("BACKLOG_AUTH_URL", "")
	tokenURL := getEnv("BACKLOG_TOKEN_URL", "")
	port := getEnv("PORT", "8081")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	// DynamoDB設定
	useDynamoDB := getEnv("USE_DYNAMODB", "false") == "true"
	dynamoDBRegion := getEnv("DYNAMODB_REGION", "ap-northeast-1")

	// OpenAI APIキーの取得
	openaiAPIKey := getEnv("OPENAI_API_KEY", "")
	if openaiAPIKey == "" {
		log.Println("Warning: OPENAI_API_KEY が見つかっていません、ダミーデータをレスオンするようになります。")
	}

	authRepo := memory.NewAuthRepository()
	var favoriteRepo model.FavoriteRepository

	// リポジトリの初期化（DynamoDBとメモリから選択）
	if useDynamoDB {
		log.Println("Using DynamoDB for favorite repository")

		var dynamoClient *dynamodb.Client
		var err error

		dynamoClient, err = dynamodb_repo.NewDynamoDBClient(dynamoDBRegion)

		if err != nil {
			log.Fatalf("Failed to create DynamoDB client: %v", err)
		}

		// DynamoDBテーブルの作成
		if err := dynamodb_repo.CreateFavoriteTable(dynamoClient); err != nil {
			log.Fatalf("Failed to create DynamoDB table: %v", err)
		}

		favoriteRepo = dynamodb_repo.NewFavoriteRepository(dynamoClient)
	} else {
		log.Println("Using in-memory favorite repository")
		favoriteRepo = memory.NewFavoriteRepository()
	}

	// OAuth設定
	oauthConfig := model.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		Scopes:       []string{"read"},
	}

	// サービスの初期化
	authService := auth.NewBacklogAuthService(oauthConfig, spaceURL)
	backlogClient := backlog.NewBacklogClient(spaceURL, clientID, clientSecret)
	backlogItemService := backlog.NewBacklogItemService(backlogClient, authRepo)

	// ユースケースの初期化
	authUseCase := usecase.NewAuthUseCase(authService, authRepo)
	backlogItemUseCase := usecase.NewBacklogItemUseCase(backlogItemService, favoriteRepo, authUseCase)

	r := gin.Default()
	// CORSミドルウェア
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 認証関連のエンドポイント https://github.com/gin-gonic/gin
	r.GET("/api/auth/url", func(c *gin.Context) {
		authURL := authUseCase.GetAuthorizationURL()
		c.JSON(http.StatusOK, gin.H{
			"url": authURL,
		})
	})

	// ヘルスチェックエンドポイント（AWS ALB用）
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    getEnv("APP_ENV", "development"),
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/api/auth/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
			return
		}

		// ユーザー認証とトークン取得（バックエンドの処理）
		token, user, err := authUseCase.AuthorizeCallback(code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// トークンとユーザー情報をURLエンコードしてフロントエンドに渡す
		tokenJSON, _ := json.Marshal(token)
		userJSON, _ := json.Marshal(user)

		// URLセーフなBase64エンコード
		tokenBase64 := base64.URLEncoding.EncodeToString(tokenJSON)
		userBase64 := base64.URLEncoding.EncodeToString(userJSON)

		// フロントエンドのコールバックページにリダイレクト
		redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&user=%s", frontendURL, tokenBase64, userBase64)
		c.Redirect(http.StatusFound, redirectURL)
	})

	r.GET("/api/auth/logout/:userId", func(c *gin.Context) {
		userID := c.Param("userId")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
			return
		}

		err := authUseCase.Logout(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Backlog更新情報関連のエンドポイント
	r.GET("/api/items", func(c *gin.Context) {
		userID := c.Query("userId")
		keyword := c.Query("keyword")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
			return
		}

		items, err := backlogItemUseCase.SearchItems(userID, keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.GET("/api/favorites/:userId", func(c *gin.Context) {
		userID := c.Param("userId")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
			return
		}

		favorites, err := backlogItemUseCase.GetFavorites(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": favorites})
	})

	r.POST("/api/favorites/:userId/:itemId", func(c *gin.Context) {
		userID := c.Param("userId")
		itemID := c.Param("itemId")

		if userID == "" || itemID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID and item ID are required"})
			return
		}

		err := backlogItemUseCase.AddFavorite(userID, itemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.DELETE("/api/favorites/:userId/:itemId", func(c *gin.Context) {
		userID := c.Param("userId")
		itemID := c.Param("itemId")

		if userID == "" || itemID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID and item ID are required"})
			return
		}

		err := backlogItemUseCase.RemoveFavorite(userID, itemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// AI分析APIエンドポイント
	r.POST("/api/ai/analyze", handleAIAnalyze)

	// フロントエンド用の静的ファイル配信
	r.StaticFS("/static", http.Dir("../../frontend/build/static"))
	r.StaticFile("/", "../../frontend/build/index.html")
	r.StaticFile("/favicon.ico", "../../frontend/build/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("../../frontend/build/index.html")
	})

	// サーバー起動
	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// handleAIAnalyze は更新情報のAI分析を行うハンドラー
func handleAIAnalyze(c *gin.Context) {
	var request struct {
		ItemID          string `json:"itemId"`
		Content         string `json:"content"`
		ProjectName     string `json:"projectName"`
		Type            string `json:"type"`
		CreatedUserName string `json:"createdUserName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// OpenAI APIキーの取得
	openaiAPIKey := getEnv("OPENAI_API_KEY", "")

	// APIキーが設定されていない場合はモックデータを返す
	if openaiAPIKey == "" {
		// モックデータ
		response := gin.H{
			"analysis": gin.H{
				"summary": fmt.Sprintf("「%s」の要約：この項目は重要な更新を含んでいます。", request.Content),
				"keyPoints": []string{
					"プロジェクトスケジュールに影響する可能性があります",
					"共同作業者との連携が必要です",
					"優先度は中程度と判断されます",
				},
				"nextActions": []string{
					"チームメンバーへの共有",
					"関連ドキュメントの更新",
					"進捗の定期的な確認",
				},
			},
		}

		// AI処理を模倣するために少し遅延を入れる
		time.Sleep(1 * time.Second)

		c.JSON(http.StatusOK, response)
		return
	}

	// OpenAI APIリクエストの作成
	openaiRequest := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "あなたはバックログ更新情報を分析するAIアシスタントです。提供された更新情報について、要約、重要ポイント、次のアクションを日本語で提案してください。",
			},
			{
				"role": "user",
				"content": fmt.Sprintf("以下のバックログ更新情報を分析してください:\n\n項目: %s\nプロジェクト: %s\n種別: %s\n作成者: %s",
					request.Content,
					request.ProjectName,
					request.Type,
					request.CreatedUserName),
			},
		},
		"temperature": 0.7,
	}

	// リクエストデータをJSONに変換
	requestJSON, err := json.Marshal(openaiRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI request"})
		return
	}

	// OpenAI APIへのリクエスト送信
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiAPIKey))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to OpenAI API"})
		return
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
		return
	}

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("OpenAI API returned error: %s", string(body)),
		})
		return
	}

	// OpenAI APIレスポンスのパース
	var openaiResponse map[string]interface{}
	if err := json.Unmarshal(body, &openaiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// レスポンスからAI生成テキストを抽出
	choices, ok := openaiResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response format from OpenAI API"})
		return
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid choice format in OpenAI response"})
		return
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid message format in OpenAI response"})
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid content format in OpenAI response"})
		return
	}

	// AIの応答を解析して構造化（簡易的な実装）
	lines := splitTextIntoSections(content)

	summary := "この項目の要約情報が生成されました。"
	var keyPoints []string
	var nextActions []string

	currentSection := ""
	for _, line := range lines {
		if line == "要約:" || line == "要約：" {
			currentSection = "summary"
			continue
		} else if line == "重要ポイント:" || line == "重要ポイント：" {
			currentSection = "keyPoints"
			continue
		} else if line == "次のアクション:" || line == "次のアクション：" ||
			line == "推奨アクション:" || line == "推奨アクション：" {
			currentSection = "nextActions"
			continue
		}

		if line == "" {
			continue
		}

		// 行の先頭の箇条書き記号を削除
		cleanLine := trimBulletPoint(line)

		switch currentSection {
		case "summary":
			summary = cleanLine
		case "keyPoints":
			keyPoints = append(keyPoints, cleanLine)
		case "nextActions":
			nextActions = append(nextActions, cleanLine)
		}
	}

	// AIの分析結果を返す
	aiAnalysis := gin.H{
		"summary":     summary,
		"keyPoints":   keyPoints,
		"nextActions": nextActions,
	}

	// keyPointsとnextActionsがnilの場合は空の配列を設定
	if keyPoints == nil {
		aiAnalysis["keyPoints"] = []string{}
	}
	if nextActions == nil {
		aiAnalysis["nextActions"] = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"analysis": aiAnalysis})
}

// splitTextIntoSections はAIから返されたテキストを行に分割します
func splitTextIntoSections(text string) []string {
	// 改行で分割
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			lines = append(lines, trimmedLine)
		}
	}
	return lines
}

// trimBulletPoint は行の先頭の箇条書き記号を削除します
func trimBulletPoint(line string) string {
	// 一般的な箇条書き記号パターンを削除
	bulletPatterns := []string{"- ", "• ", "* ", "・", "1. ", "2. ", "3. ", "4. ", "5. ", "①", "②", "③", "④", "⑤"}

	trimmedLine := strings.TrimSpace(line)
	for _, pattern := range bulletPatterns {
		if strings.HasPrefix(trimmedLine, pattern) {
			return strings.TrimSpace(trimmedLine[len(pattern):])
		}
	}

	return trimmedLine
}
