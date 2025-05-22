package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yuchi1128/go-login-app/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDB はデータベースへの接続を確立します
func ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	// 接続確認
	err = DB.Ping()
	if err != nil {
		log.Fatalf("データベースPingエラー: %v", err)
	}

	log.Println("データベース接続成功！")

	// 接続プールの設定 (任意)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// ユーザテーブルを作成 (存在しない場合のみ)
	createUsersTable()
}

// createUsersTable はusersテーブルを作成します
func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("usersテーブル作成エラー: %v", err)
	}
	log.Println("usersテーブル準備完了 (または既に存在します)")
}

// CreateUser は新しいユーザーをデータベースに登録します
func CreateUser(user *models.User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash)
			   VALUES ($1, $2, $3) RETURNING id`
	var userID int
	err := DB.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("ユーザー作成エラー: %w", err)
	}
	return userID, nil
}

// GetUserByEmailOrUsername はEmailまたはUsernameでユーザーを検索します
func GetUserByEmailOrUsername(identifier string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE email = $1 OR username = $1`
	err := DB.QueryRow(query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		return nil, fmt.Errorf("ユーザー検索エラー (%s): %w", identifier, err)
	}
	return user, nil
}

// GetUserByEmail はEmailでユーザーを検索します（重複チェック用）
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id FROM users WHERE email = $1`
	err := DB.QueryRow(query, email).Scan(&user.ID) // IDだけ取得できれば存在確認は可能
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Email検索エラー: %w", err)
	}
	return user, nil
}

// GetUserByUsername はUsernameでユーザーを検索します（重複チェック用）
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id FROM users WHERE username = $1`
	err := DB.QueryRow(query, username).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Username検索エラー: %w", err)
	}
	return user, nil
}