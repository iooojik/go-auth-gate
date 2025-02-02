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
	clientSecret, err := r.tokenGenerator(r.cfg.tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	u, err := prepareLink(r.cfg.URL, url.Values{
		"client_id":     []string{r.cfg.tokenConfig.ClientID},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{gen.Code},
		"client_secret": []string{clientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("prepare url: %w", err)
	}

	u.Path = path.Join(u.Path, "/auth/token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := r.cfg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	return decode(resp.Body)
}

func (r *Client) RefreshToken(ctx context.Context, refresh Refresh) (*AuthCode, error) {
	clientSecret, err := r.tokenGenerator(r.cfg.tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	u, err := prepareLink(r.cfg.URL, url.Values{
		"client_id":     []string{r.cfg.tokenConfig.ClientID},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refresh.RefreshToken},
		"client_secret": []string{clientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("prepare url: %w", err)
	}

	u.Path = path.Join(u.Path, "/auth/token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := r.cfg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	return decode(resp.Body)
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}

		return fmt.Errorf("%w: %s", ErrBadRequest, string(body))
	}

	return nil
}

func decode(r io.Reader) (*AuthCode, error) {
	authCode := new(AuthCode)

	err := json.NewDecoder(r).Decode(authCode)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return authCode, nil
}

func prepareLink(link string, params url.Values) (*url.URL, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, "/auth/token")

	u.RawQuery = params.Encode()

	return u, nil
}
