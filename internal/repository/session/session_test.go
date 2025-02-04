package session_test

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
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
				cfg: session.Config{
					SessionDuration: 3600,
				},
				loginInfo: model.LoginInfo{
					UserID:    "test-apple-id",
					Token:     "token-123",
					TokenType: "apple_sign_in",
				},
			},
			updateArgs: args{
				cfg: session.Config{
					SessionDuration: 7200,
				},
				loginInfo: model.LoginInfo{
					UserID:    "test-apple-id",
					Token:     "token-456",
					TokenType: "apple_sign_in",
				},
			},
			insertCheck: func() {
				var user session.User

				err := s.db.Get(&user, "SELECT * FROM users WHERE user_id = ?", "test-apple-id")
				s.Require().NoError(err)
				s.Require().Equal(3600, user.SessionDuration)
				s.Require().Equal("apple_sign_in", user.AuthType)

				var userToken session.UserToken
				err = s.db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.UserID)
				s.Require().NoError(err)
				s.Require().Equal("token-123", userToken.Token)
			},
			updateCheck: func() {
				var user session.User

				err := s.db.Get(&user, "SELECT * FROM users WHERE user_id = ?", "test-apple-id")
				s.Require().NoError(err)
				s.Require().Equal(7200, user.SessionDuration)
				s.Require().Equal("apple_sign_in", user.AuthType)

				var userToken session.UserToken
				err = s.db.Get(&userToken, "SELECT * FROM user_tokens WHERE user_id = ?", user.UserID)
				s.Require().NoError(err)
				s.Require().Equal("token-456", userToken.Token)
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
