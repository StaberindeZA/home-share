// Package home, also commonly known as a workspace is a
// grouping of related users. Users can be members of
// many Homes
package home

import (
	"encoding/json"
	"log"
	"net/http"
)

type HomeController struct {
	logic HomeLogic
}

func NewHomeController(logic HomeLogic) HomeController {
	return HomeController{
		logic: logic,
	}
}

func (c HomeController) ListHomeMates(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}

	mates, err := c.logic.RetrieveMates(slug)
	if err != nil {
		log.Println("Error during RetrieveMates", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	var data []HomeMatesDTO
	for _, mate := range mates {
		data = append(data, HomeMatesDTO{
			ID:    mate.Id,
			Email: mate.Email,
			Name:  mate.Name,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error while encoding HomeMatesDTO", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
