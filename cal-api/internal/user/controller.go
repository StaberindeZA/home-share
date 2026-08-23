package user

import (
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
	user, ok := RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	data := UserProfileDTO{
		Name:  user.Name,
		Email: user.Email,
	}
	utils.SendPayloadJSON(w, data)
}

func (c UserController) UpdateLoggedInUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	user, ok := RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	updateUser, ok := utils.DecodeBodyJSON[UpdateUserDTO](w, r)
	if !ok {
		return
	}

	message, err := c.logic.UpdateByEmail(user.Email, updateUser.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(message))
}
