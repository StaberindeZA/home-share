package googleauth

import (
	"cal-api/internal/utils"

	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type GoogleAuthController struct {
	logic GoogleAuthLogic
}

func (c GoogleAuthController) VerifyIdToken(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData VerifyIdTokenDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	email, err := c.logic.ValidateIdToken(requestData.IdToken)
	if err != nil {
		if errors.Is(err, ErrClientIDNotFound) {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, err := utils.GenerateBackendJWT(email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := GoogleAuthTokenDTO{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewGoogleAuthController(logic GoogleAuthLogic) GoogleAuthController {
	return GoogleAuthController{
		logic: logic,
	}
}
