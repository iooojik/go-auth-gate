package jwt_test

import (
	"reflect"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		key    string
		domain string
	}
	tests := []struct {
		name string
		args args
		want TokenGenerator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(tt.args.key, tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		headerToken string
		secret      string
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenClaims
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateToken(tt.args.headerToken, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
