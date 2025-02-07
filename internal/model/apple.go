package model

type Refresh struct {
	RefreshToken string
	UserID       string
}

type Generate struct {
	Code   string
	UserID string
}
