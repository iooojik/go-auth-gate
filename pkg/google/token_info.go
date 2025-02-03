package google

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	json "github.com/json-iterator/go"
)

func (r *Client) TokenInfo(ctx context.Context, token string) (*TokenInfo, error) {
	u, err := prepareLink(r.cfg.URL, url.Values{
		"id_token": []string{token},
	})
	if err != nil {
		return nil, fmt.Errorf("prepare url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	return decode(resp.Body)
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest && resp.StatusCode >= http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	errMessage := new(ErrResponse)

	err = json.Unmarshal(body, errMessage)
	if err != nil {
		return fmt.Errorf("decode error message: %w", err)
	}

	if errMessage.Error == "invalid_token" {
		return ErrTokenRevokedOrInvalid
	}

	return fmt.Errorf("%w: %s", ErrBadRequest, string(body))
}

func decode(r io.ReadCloser) (*TokenInfo, error) {
	defer func() { _ = r.Close() }()

	authCode := new(TokenInfo)

	err := json.NewDecoder(r).Decode(authCode)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return authCode, nil
}

// prepareLink appends params and path to link.
func prepareLink(link string, params url.Values) (*url.URL, error) {
	if link == "" {
		return nil, ErrEmptyLink
	}

	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, "/tokeninfo")

	u.RawQuery = params.Encode()

	return u, nil
}
