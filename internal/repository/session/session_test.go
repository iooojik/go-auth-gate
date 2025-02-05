package session_test

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

func (s *RepositoryTestSuite) TestRepository_Login() {
	type args struct {
		cfg       session.Config
		loginInfo model.LoginInfo
	}

	tests := []struct {
		name        string
		insertArgs  args
		updateArgs  args
		insertCheck func()
		updateCheck func()
		wantErr     bool
	}{
		{
			name: "test#1",
			insertArgs: args{
				cfg: session.Config{},
				loginInfo: model.LoginInfo{
					UserID: "test-apple-id",
					AppleTokenInfo: &apple.AuthCode{
						AccessToken:  "access-token-123",
						TokenType:    "Bearer",
						ExpiresIn:    100,
						RefreshToken: "refresh-token-123",
						IDToken:      "id-token-123",
					},
				},
			},
			updateArgs: args{
				cfg: session.Config{},
				loginInfo: model.LoginInfo{
					UserID: "test-apple-id",
					AppleTokenInfo: &apple.AuthCode{
						AccessToken:  "access-token-456",
						TokenType:    "Bearer",
						ExpiresIn:    200,
						RefreshToken: "refresh-token-456",
						IDToken:      "id-token-456",
					},
				},
			},
			insertCheck: func() {
				var user session.User

				err := s.db.Get(&user, "SELECT * FROM users WHERE user_id = ?", "test-apple-id")
				s.Require().NoError(err)
				s.Require().Equal("test-apple-id", user.UserID)

				var userToken session.UserToken

				err = s.db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.UserID)
				s.Require().NoError(err)
				s.Require().Equal("test-apple-id", userToken.UserID)
				s.Require().Equal("access-token-123", userToken.AccessToken)
				s.Require().Equal("Bearer", userToken.TokenType)
				s.Require().Equal(100, userToken.ExpiresIn)
				s.Require().Equal("refresh-token-123", userToken.RefreshToken)
				s.Require().Equal("id-token-123", userToken.IDToken)
			},
			updateCheck: func() {
				var user session.User

				err := s.db.Get(&user, "SELECT * FROM users WHERE user_id = ?", "test-apple-id")
				s.Require().NoError(err)
				s.Require().Equal("test-apple-id", user.UserID)

				var userToken session.UserToken

				err = s.db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.UserID)
				s.Require().NoError(err)
				s.Require().Equal("test-apple-id", userToken.UserID)
				s.Require().Equal("access-token-456", userToken.AccessToken)
				s.Require().Equal("Bearer", userToken.TokenType)
				s.Require().Equal(200, userToken.ExpiresIn)
				s.Require().Equal("refresh-token-456", userToken.RefreshToken)
				s.Require().Equal("id-token-456", userToken.IDToken)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		r := session.New(tt.insertArgs.cfg, s.db)

		// insert
		err := r.Login(context.Background(), tt.insertArgs.loginInfo)
		if !tt.wantErr {
			s.Require().NoError(err)
		} else {
			s.Require().Error(err)
		}

		tt.insertCheck()

		r = session.New(tt.updateArgs.cfg, s.db)

		// update
		err = r.Login(context.Background(), tt.updateArgs.loginInfo)
		if !tt.wantErr {
			s.Require().NoError(err)
		} else {
			s.Require().Error(err)
		}

		tt.updateCheck()
	}
}
