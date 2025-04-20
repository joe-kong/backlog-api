package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/domain/model"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/auth"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/backlog"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/infrastructure/persistence/memory"
	"nulab-exam.backlog.jp/KOU/app/backend/internal/usecase"
)

func main() {
	// 環境変数から設定を読み込む（またはデフォルト値を使用）
	spaceURL := getEnv("BACKLOG_SPACE_URL", "https://nulab-exam.backlog.jp/projects/KOU ")
	clientID := getEnv("BACKLOG_CLIENT_ID", "QgcVk8WlUb4aJZ8GbNrja1ATXXDFA60y")
	clientSecret := getEnv("BACKLOG_CLIENT_SECRET", "6TmFwzdeYE0TKi8mmXmWJ3d14NmL1SqwTgIbu4Ud1ZFwo8x3raCCGIHhzEmPqk7c")
	redirectURI := getEnv("OAUTH_REDIRECT_URI", "http://localhost:8080/api/auth/callback")
	port := getEnv("PORT", "8080")

	// リポジトリの初期化
	authRepo := memory.NewAuthRepository()
	favoriteRepo := memory.NewFavoriteRepository()

	// OAuth設定
	oauthConfig := model.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		AuthURL:      spaceURL + "/OAuth2/authorize",
		TokenURL:     spaceURL + "/OAuth2/token",
		Scopes:       []string{"read"},
	}

	// サービスの初期化
	authService := auth.NewBacklogAuthService(oauthConfig, spaceURL)
	backlogClient := backlog.NewBacklogClient(spaceURL, clientID, clientSecret)
	backlogItemService := backlog.NewBacklogItemService(backlogClient, authRepo)

	// ユースケースの初期化
	authUseCase := usecase.NewAuthUseCase(authService, authRepo)
	backlogItemUseCase := usecase.NewBacklogItemUseCase(backlogItemService, favoriteRepo, authUseCase)

	// Ginルーターの初期化
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

	// 認証関連のエンドポイント
	r.GET("/api/auth/url", func(c *gin.Context) {
		authURL := authUseCase.GetAuthorizationURL()
		c.JSON(http.StatusOK, gin.H{
			"url": authURL,
		})
	})

	r.GET("/api/auth/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
			return
		}

		token, user, err := authUseCase.AuthorizeCallback(code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  user,
		})
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

	// フロントエンド用の静的ファイル配信
	r.StaticFS("/static", http.Dir("../../frontend/build/static"))
	r.StaticFile("/", "../../frontend/build/index.html")
	r.StaticFile("/favicon.ico", "../../frontend/build/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("../../frontend/build/index.html")
	})

	// 以前の認証エンドポイントも下位互換性のために残す
	r.GET("/auth/url", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/auth/url")
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/api/auth/callback?code=%s", c.Query("code")))
	})

	r.GET("/auth/logout/:userId", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/api/auth/logout/%s", c.Param("userId")))
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
