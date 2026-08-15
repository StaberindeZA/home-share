package user

import (
	"encoding/json"
	"log"
	"net/http"

	"cal-api/internal/utils"
)

type UserController struct {
	logic UserLogic
}

func NewUserController(logic UserLogic) UserController {
	return UserController{
		logic: logic,
	}
}

func (c UserController) ReadLoggedInUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetLoggedInUser[User](r)
	if err != nil {
		log.Printf("Error in GetLoggedInUser: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := UserProfileDTO{
		Name:  user.Name,
		Email: user.Email,
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c UserController) UpdateLoggedInUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	user, err := utils.GetLoggedInUser[User](r)
	if err != nil {
		log.Printf("Error in GetLoggedInUser: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var updateUser UpdateUserDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updateUser)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	message, err := c.logic.UpdateByEmail(user.Email, updateUser.Name)
	if err != nil {
		log.Printf("Error in update user logic: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(message))
}
