// Package auth is used for Authentication and Authorization logic
package auth

import (
	"net/http"

	"cal-api/internal/home"
	"cal-api/internal/homemate"
	"cal-api/internal/user"
	"cal-api/internal/utils"
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
	user, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	data := UserInfoDTO{
		Email: user.Email,
		Name:  user.Name,
	}

	utils.SendPayloadJSON(w, data)
}

func (ac AuthController) VerifyRole(w http.ResponseWriter, r *http.Request) {
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}

	u, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	payload, ok := utils.DecodeBodyJSON[VerifyRoleDTO](w, r)
	if !ok {
		return
	}

	h, err := ac.homeLogic.Read(payload.HomeSlug)
	if err != nil {
		utils.RecordErrorServer(w, err, "VerifyRole.Read")
		return
	}

	hm, err := ac.homemateLogic.ReadForHomeAndMate(h.ID, u.Id)
	if err != nil {
		utils.RecordErrorServer(w, err, "VerifyRole.ReadForHomeAndMate")
		return
	}

	// TODO: Improve this logic to support multiple roles
	if hm.Role.String() != payload.Role {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	w.Write([]byte("ok"))
}
