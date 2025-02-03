package session_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
)

// Тест на добавление пользователя и обновление токена
func TestUserAuthFlow(t *testing.T) {
	// teardown := setupTestDB(t)
	// defer teardown()
	//
	// appleID := "test-apple-id"
	// token1 := "token-123"
	// token2 := "token-456"
	// sessionDuration := 3600
	//
	//
	// Проверяем, что пользователь создан
	// var user User
	// err = db.Get(&user, "SELECT * FROM users WHERE apple_id = ?", appleID)
	// assert.NoError(t, err)
	// assert.Equal(t, appleID, user.AppleID)
	// assert.Equal(t, sessionDuration, user.SessionDuration)
	// assert.Equal(t, "apple_sign_in", user.AuthType)
	//
	// // Проверяем, что токен создан
	// var userToken UserToken
	// err = db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.ID)
	// assert.NoError(t, err)
	// assert.Equal(t, token1, userToken.Token)
	//
	// // 2. Второй запрос: тот же пользователь, но с новым токеном
	// _, err = db.Exec(query, appleID, sessionDuration, appleID, token2, appleID)
	// assert.NoError(t, err, "Error updating token")
	//
	// // Проверяем, что токен обновился
	// err = db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.ID)
	// assert.NoError(t, err)
	// assert.Equal(t, token2, userToken.Token, "Token should be updated")
}

func TestRepository_Login(t *testing.T) {
	type args struct {
		ctx       context.Context
		loginInfo model.LoginInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &session.Repository{}
			if err := r.Login(tt.args.ctx, tt.args.loginInfo); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
