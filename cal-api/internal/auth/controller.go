// Package auth is used for Authentication and Authorization logic
package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"cal-api/internal/home"
	"cal-api/internal/homemate"
	"cal-api/internal/user"
)

type AuthController struct {
	userLogic     user.UserLogic
	homeLogic     home.HomeLogic
	homemateLogic homemate.HomeMateLogic
}

func NewAuthController(userLogic user.UserLogic, homeLogic home.HomeLogic, homemateLogic homemate.HomeMateLogic) AuthController {
	return AuthController{
		userLogic:     userLogic,
		homeLogic:     homeLogic,
		homemateLogic: homemateLogic,
	}
}

func (ac AuthController) UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(user.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	data := UserInfoDTO{
		Email: user.Email,
		Name:  user.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac AuthController) VerifyRole(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	u, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	var payload VerifyRoleDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	h, err := ac.homeLogic.Read(payload.HomeSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hm, err := ac.homemateLogic.ReadForHomeAndMate(h.ID, u.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO: Improve this logic to support multiple roles
	if hm.Role.String() != payload.Role {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	w.Write([]byte("ok"))
}
