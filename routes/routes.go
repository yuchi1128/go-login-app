package routes

import (
	"net/http"
	"strings"
	"time"

	"your-github-username/go-login-app/handlers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_very_secret_key_that_should_be_in_env")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorizationヘッダーが必要です"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorizationヘッダーの形式が無効です (Bearer <token>)"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 署名方法の検証
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンシグネチャです"})
				c.Abort()
				return
			}
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "不正な形式のトークンです"})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンが期限切れか、まだ有効ではありません"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンを検証できませんでした: " + err.Error()})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンを検証できませんでした: " + err.Error()})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// トークンが有効なら、claimsからユーザー情報を取得してコンテキストにセット
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		// 有効期限のチェック (ParseWithClaims内でも行われるが念のため)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンが期限切れです"})
			c.Abort()
			return
		}


		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	apiRoutes := router.Group("/api")
	apiRoutes.Use(AuthMiddleware())
	{
		apiRoutes.GET("/profile", handlers.Profile)
	}

	// ヘルスチェック用エンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})


	return router
}