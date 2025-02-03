package apple

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	json "github.com/json-iterator/go"
)

func (r *Client) ReceiveToken(ctx context.Context, gen Generate) (*AuthCode, error) {
	clientSecret, err := r.tokenGenerator(r.cfg.TokenConfig)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	u, err := prepareLink(r.cfg.URL, url.Values{
		"client_id":     []string{r.cfg.TokenConfig.ClientID},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{gen.Code},
		"client_secret": []string{clientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("prepare url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
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

func (r *Client) RefreshToken(ctx context.Context, refresh Refresh) (*AuthCode, error) {
	clientSecret, err := r.tokenGenerator(r.cfg.TokenConfig)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	u, err := prepareLink(r.cfg.URL, url.Values{
		"client_id":     []string{r.cfg.TokenConfig.ClientID},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refresh.RefreshToken},
		"client_secret": []string{clientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("prepare url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
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

	errMessage := new(ErrorMessage)

	err = json.Unmarshal(body, errMessage)
	if err != nil {
		return fmt.Errorf("decode error message: %w", err)
	}

	if errMessage.ErrorDescription == "The code has expired or has been revoked." {
		return ErrTokenRevoked
	}

	return fmt.Errorf("%w: %s", ErrBadRequest, string(body))
}

func decode(r io.ReadCloser) (*AuthCode, error) {
	defer func() { _ = r.Close() }()

	authCode := new(AuthCode)

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

	u.Path = path.Join(u.Path, "/auth/token")

	u.RawQuery = params.Encode()

	return u, nil
}
