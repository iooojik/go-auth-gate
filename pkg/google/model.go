package google

type ErrResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TokenInfo struct{}
