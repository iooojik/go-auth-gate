package google

import (
	"context"
	"fmt"
	"slices"
)

func (r *Client) CheckToken(ctx context.Context, token string) (string, error) {
	tokenInfo, err := r.TokenInfo(ctx, token)
	if err != nil {
		return "", fmt.Errorf("fetch token info: %w", err)
	}

	if !slices.Contains(r.cfg.AppID, tokenInfo.Aud) {
		return "", nil
	}

	return tokenInfo.Sub, nil
}
