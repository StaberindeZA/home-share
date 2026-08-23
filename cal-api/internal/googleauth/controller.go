// Package googleauth handles google auth to internal jwt conversion
package googleauth

import (
	"errors"
	"net/http"

	"cal-api/internal/utils"
)

type GoogleAuthController struct {
	logic GoogleAuthLogic
}

func (c GoogleAuthController) VerifyIDToken(w http.ResponseWriter, r *http.Request) {
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}

	requestData, ok := utils.DecodeBodyJSON[VerifyIDTokenDTO](w, r)
	if !ok {
		return
	}

	email, err := c.logic.ValidateIDToken(requestData.IDToken)
	if err != nil {
		if errors.Is(err, ErrClientIDNotFound) {
			utils.RecordErrorServer(w, err, "VerifyIDToken.ValidateIDToken")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, err := utils.GenerateBackendJWT(email)
	if err != nil {
		utils.RecordErrorServer(w, err, "VerifyIDToken.GenerateBackendJWT")
		return
	}

	data := GoogleAuthTokenDTO{
		Token: tokenString,
	}

	utils.SendPayloadJSON(w, data)
}

func NewGoogleAuthController(logic GoogleAuthLogic) GoogleAuthController {
	return GoogleAuthController{
		logic: logic,
	}
}
