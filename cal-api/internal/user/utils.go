package user

import (
	"net/http"
)

func RetrieveUserFromContext(w http.ResponseWriter, r *http.Request) (User, bool) {
	u, ok := r.Context().Value("user").(User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return u, false
	}
	return u, true
}
