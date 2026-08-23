package user

import (
	"errors"
	"net/http"

	"cal-api/internal/utils"
)

func RetrieveUserFromContext(w http.ResponseWriter, r *http.Request) (User, bool) {
	u, ok := r.Context().Value("user").(User)
	if !ok {
		utils.RecordErrorServer(w, errors.New("user missing from context"), "RetrieveUserFromContext")
		return u, false
	}
	return u, true
}
