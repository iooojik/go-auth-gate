package jwt

type ctxKey int

const UCtxKey = ctxKey(0)

type TokenUser struct {
	ID string `json:"id"`
}
