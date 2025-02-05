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

				err = s.db.Get(&userToken, "SELECT * FROM apple_tokens WHERE user_id = ?", user.UserID)
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

				err = s.db.Get(&userToken, "SELECT * FROM apple_tokens WHERE user_id = ?", user.UserID)
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

func (s *RepositoryTestSuite) TestRepository_FetchAll() {
	type fields struct {
		cfg session.Config
	}

	type args struct {
		authType model.TokenType
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(ctx context.Context, r *session.Repository)
		want    []model.Refresh
		wantErr bool
	}{
		{
			name: "test#1",
			fields: fields{
				cfg: session.Config{},
			},
			args: args{
				authType: model.GoogleSignInAuth,
			},
			setup:   func(_ context.Context, _ *session.Repository) {},
			want:    []model.Refresh{},
			wantErr: true,
		},
		{
			name: "test#2",
			fields: fields{
				cfg: session.Config{},
			},
			args: args{
				authType: model.AppleID,
			},
			setup: func(ctx context.Context, r *session.Repository) {
				tx, err := s.db.Begin()
				s.Require().NoError(err)

				err = r.InsertAppleIDToken(ctx, tx, apple.AuthCode{
					AccessToken:  "access-token-123",
					TokenType:    "token-type-123",
					ExpiresIn:    33,
					RefreshToken: "refresh-token-123",
					IDToken:      "id-token-123",
				}, "user_1")
				s.Require().NoError(err)

				err = r.InsertAppleIDToken(ctx, tx, apple.AuthCode{
					AccessToken:  "access-token-456",
					TokenType:    "token-type-456",
					ExpiresIn:    44,
					RefreshToken: "refresh-token-456",
					IDToken:      "id-token-456",
				}, "user_2")
				s.Require().NoError(err)

				err = tx.Commit()
				s.Require().NoError(err)
			},
			want: []model.Refresh{
				{
					RefreshToken: "refresh-token-123",
					UserID:       "user_1",
				},
				{
					RefreshToken: "refresh-token-456",
					UserID:       "user_2",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			r := session.New(tt.fields.cfg, s.db)

			ctx := context.Background()

			tt.setup(ctx, r)

			got, err := r.FetchAll(ctx, tt.args.authType)
			if tt.wantErr {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)

			var tokens []model.Refresh

			for t, err := range got {
				s.Require().NoError(err)

				tokens = append(tokens, t)
			}

			s.Require().Equal(tt.want, tokens)
		})
	}
}
