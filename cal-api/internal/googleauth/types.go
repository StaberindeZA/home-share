package googleauth

type GoogleAuthLogic interface {
	ValidateIDToken(idToken string) (string, error)
}
