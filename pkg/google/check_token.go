package google

import (
	"context"
	"fmt"
)

func (r *Client) CheckToken(ctx context.Context, token string) (bool, error) {
	tokenInfo, err := r.TokenInfo(ctx, token)
	if err != nil {
		return false, fmt.Errorf("fetch token info: %w", err)
	}

	if tokenInfo.Aud != r.cfg.AppID {
		return false, nil
	}

	return true, nil
}
