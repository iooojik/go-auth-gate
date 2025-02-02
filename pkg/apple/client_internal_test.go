package apple

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

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

func Test_prepareLink(t *testing.T) {
	type args struct {
		link   string
		params url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apple.prepareLink(tt.args.link, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
