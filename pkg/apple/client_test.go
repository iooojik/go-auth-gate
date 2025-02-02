package apple_test

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/iooojik/go-auth-gate/pkg/apple"
)

func TestClient_RefreshToken(t *testing.T) {
	type fields struct {
		cfg            Config
		tokenGenerator ClientSecretGenerator
	}
	type args struct {
		ctx     context.Context
		refresh Refresh
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AuthCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Client{
				cfg:            tt.fields.cfg,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, err := r.RefreshToken(tt.args.ctx, tt.args.refresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ReceiveToken(t *testing.T) {
	type fields struct {
		cfg            apple.Config
		tokenGenerator apple.ClientSecretGenerator
	}
	type args struct {
		ctx context.Context
		gen apple.Generate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AuthCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &apple.Client{
				cfg:            tt.fields.cfg,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, err := r.ReceiveToken(tt.args.ctx, tt.args.gen)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceiveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReceiveToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkResponse(t *testing.T) {
	type args struct {
		resp *http.Response
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
			if err := apple.checkResponse(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("checkResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
