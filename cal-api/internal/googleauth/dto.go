package googleauth

type VerifyIdTokenDTO struct {
	IdToken string `json:"idToken"`
}

type GoogleAuthTokenDTO struct {
	Token string `json:"token"`
}
