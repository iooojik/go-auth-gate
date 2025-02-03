package user

import (
	"context"
	"testing"

	"github.com/iooojik/go-auth-gate/internal/model"
)

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
			r := &Repository{}
			if err := r.Login(tt.args.ctx, tt.args.loginInfo); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
