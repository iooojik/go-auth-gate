package service

import (
	"context"
	"testing"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func TestService_AppleCallback(t *testing.T) {
	type fields struct {
		appleSignIn AppleSignIn
		userRepo    UserRepository
	}
	type args struct {
		ctx          context.Context
		callbackInfo model.AppleCallback
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				appleSignIn: tt.fields.appleSignIn,
				userRepo:    tt.fields.userRepo,
			}
			if err := s.AppleCallback(tt.args.ctx, tt.args.callbackInfo); (err != nil) != tt.wantErr {
				t.Errorf("AppleCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
