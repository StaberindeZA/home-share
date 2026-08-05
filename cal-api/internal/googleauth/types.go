package googleauth

type GoogleAuthLogic interface {
	ValidateIdToken(idToken string) (string, error)
}
