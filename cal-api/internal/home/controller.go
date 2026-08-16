package home

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cal-api/internal/homemate"
	"cal-api/internal/user"
)

type HomeController struct {
	logic         HomeLogic
	homeMateLogic homemate.HomeMateLogic
	userLogic     user.UserLogic
}

func NewHomeController(logic HomeLogic, homeMateLogic homemate.HomeMateLogic, userLogic user.UserLogic) HomeController {
	return HomeController{
		logic:         logic,
		homeMateLogic: homeMateLogic,
		userLogic:     userLogic,
	}
}

func (c HomeController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	currentUser, ok := r.Context().Value("user").(user.User)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var createHome HomeCreateDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createHome)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	slug, err := c.logic.Create(createHome.Name, createHome.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	home, err := c.logic.Read(slug)
	if err != nil {
		log.Printf("Could not find home with slug: %s\n%v", slug, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = c.homeMateLogic.Create(home.ID, currentUser.Id, homemate.Admin)
	if err != nil {
		log.Printf("Could not create home mate: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
}

func (c HomeController) Read(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}
	home, err := c.logic.Read(slug)
	if err != nil {
		log.Println("Error during Read", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	data := HomeDTO{
		Name:        home.name,
		Slug:        home.slug,
		Description: home.description,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c HomeController) List(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("user").(user.User)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	homes, err := c.logic.List(currentUser.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var data []HomeDTO
	for _, home := range homes {
		data = append(data, HomeDTO{
			Name:        home.name,
			Slug:        home.slug,
			Description: home.description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c HomeController) Delete(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}

	_, _, err := c.VerifyUserForAdminAction(slug, r)
	if err != nil {
		log.Printf("Error during Authorization: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = c.logic.Delete(slug)
	if err != nil {
		log.Printf("Error during Delete: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func (c HomeController) ReadHomeMates(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}

	mates, err := c.logic.ReadMates(slug)
	if err != nil {
		log.Println("Error during RetrieveMates", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	var data []HomeMatesDTO
	for _, mate := range mates {
		data = append(data, HomeMatesDTO{
			ID:    mate.ID,
			Email: mate.Email,
			Name:  mate.Name,
			Role:  mate.Role.String(),
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

func (c HomeController) CreateHomeMate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}

	_, h, err := c.VerifyUserForAdminAction(slug, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var createHomeMates CreateHomeMatesDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&createHomeMates)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	u, err := c.userLogic.FindOrCreate(0, createHomeMates.Name, createHomeMates.Email)
	if err != nil {
		log.Printf("Could not find or create mate: %s\n%v", createHomeMates.Email, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = c.homeMateLogic.Create(h.ID, u.Id, homemate.Mate)
	if err != nil {
		log.Printf("Could not find or create home mate: %s\n%v", createHomeMates.Email, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func (c HomeController) DeleteHomeMate(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing email query param", http.StatusBadRequest)
	}

	_, h, err := c.VerifyUserForAdminAction(slug, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	u, err := c.userLogic.FindByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Could not find mate", http.StatusBadRequest)
		} else {
			log.Printf("Error finding home mate: %s\n%v", email, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	err = c.homeMateLogic.Delete(h.ID, u.Id)
	if err != nil {
		log.Printf("Error deleting home mate: %s - %s", slug, email)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Write([]byte("ok"))
}

func (c HomeController) VerifyUserForAdminAction(slug string, r *http.Request) (user.User, Home, error) {
	var u user.User
	var h Home
	var ok bool
	h, err := c.logic.Read(slug)
	if err != nil {
		return u, h, err
	}

	u, ok = r.Context().Value("user").(user.User)
	if !ok {
		return u, h, errors.New("user not in context")
	}

	hm, err := c.homeMateLogic.ReadForHomeAndMate(h.ID, u.Id)
	if err == sql.ErrNoRows {
		return u, h, errors.New("user is not authorized for this action")
	}

	if err != nil {
		return u, h, err
	}

	if hm.Role != homemate.Admin {
		return u, h, errors.New("user is not authorized for this action")
	}

	return u, h, nil
}
