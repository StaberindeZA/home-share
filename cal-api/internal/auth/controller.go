package auth

import (
	"cal-api/internal/user"

	"encoding/json"
	"log"
	"net/http"
)

type AuthController struct {
	userLogic user.UserLogic
}

func NewAuthController(userLogic user.UserLogic) AuthController {
	return AuthController{
		userLogic: userLogic,
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
