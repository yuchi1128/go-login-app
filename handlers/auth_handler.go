package handlers

import (
	"net/http"
	"time"

	"github.com/yuchi1128/go-login-app/database"
	"github.com/yuchi1128/go-login-app/models"
	"github.com/yuchi1128/go-login-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT署名キー (実際には環境変数などから読み込むべき秘密の値)
var jwtKey = []byte("your_very_secret_key_that_should_be_in_env")

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Register はユーザー登録処理を行います
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力が無効です: " + err.Error()})
		return
	}

	// Emailの重複チェック
	existingUser, _ := database.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に使用されています"})
		return
	}

	// Usernameの重複チェック
	existingUser, _ = database.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "このユーザー名は既に使用されています"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userID, err := database.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました: " + err.Error()})
		return
	}
	user.ID = userID

	c.JSON(http.StatusCreated, gin.H{
		"message": "ユーザー登録が成功しました",
		"user_id": userID,
		"username": user.Username,
		"email": user.Email,
	})
}

// Login はログイン処理を行います
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力が無効です: " + err.Error()})
		return
	}

	user, err := database.GetUserByEmailOrUsername(req.Identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索中にエラーが発生しました"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレス/ユーザー名またはパスワードが正しくありません"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレス/ユーザー名またはパスワードが正しくありません"})
		return
	}

	// JWTトークン生成
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-login-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ログイン成功",
		"token":   tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Profile(c *gin.Context) {
	username, _ := c.Get("username")
	userID, _ := c.Get("user_id")

	if username == nil || userID == nil {
		// これは実際にはミドルウェアが弾くので、ここには到達しない想定
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ようこそ、保護されたエリアへ！",
		"user_id": userID,
		"username": username,
	})
}