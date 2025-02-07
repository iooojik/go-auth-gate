package google

import (
	"context"
	"fmt"
)

func (r *Client) CheckToken(ctx context.Context, token string) (string, error) {
	tokenInfo, err := r.TokenInfo(ctx, token)
	if err != nil {
		return "", fmt.Errorf("fetch token info: %w", err)
	}

	if tokenInfo.Aud != r.cfg.AppID {
		return "", nil
	}

	return tokenInfo.Sub, nil
}
