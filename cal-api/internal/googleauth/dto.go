package googleauth

type VerifyIDTokenDTO struct {
	IDToken string `json:"idToken"`
}

type GoogleAuthTokenDTO struct {
	Token string `json:"token"`
}
