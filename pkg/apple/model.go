package apple

type Refresh struct {
	RefreshToken string
}

type Generate struct {
	Code string
}

type AuthCode struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type ErrorMessage struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
