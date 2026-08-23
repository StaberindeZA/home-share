package home

import (
	"database/sql"
	"errors"
	"net/http"

	"cal-api/internal/homemate"
	"cal-api/internal/user"
	"cal-api/internal/utils"
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
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}

	currentUser, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	createHome, ok := utils.DecodeBodyJSON[HomeCreateDTO](w, r)
	if !ok {
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
		utils.RecordErrorServer(w, err, "HomeController.Create.Read")
		return
	}

	err = c.homeMateLogic.Create(home.ID, currentUser.Id, homemate.Admin)
	if err != nil {
		utils.RecordErrorServer(w, err, "HomeController.Create.Create")
		return
	}
	w.Write([]byte("ok"))
}

func (c HomeController) Read(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
		return
	}
	home, err := c.logic.Read(slug)
	if err != nil {
		utils.RecordErrorServer(w, err, "HomeController.Read.Read")
		return
	}

	data := HomeDTO{
		Name:        home.name,
		Slug:        home.slug,
		Description: home.description,
	}

	utils.SendPayloadJSON(w, data)
}

func (c HomeController) List(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
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

	utils.SendPayloadJSON(w, data)
}

func (c HomeController) Delete(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
		return
	}

	_, _, err := c.VerifyUserForAdminAction(slug, w, r)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.Delete.VerifyUserForAdminAction")
		return
	}

	err = c.logic.Delete(slug)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.Delete.Delete")
		return
	}

	w.Write([]byte("ok"))
}

func (c HomeController) ReadHomeMates(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
		return
	}

	mates, err := c.logic.ReadMates(slug)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.ReadHomeMates.ReadMates")
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
	utils.SendPayloadJSON(w, data)
}

func (c HomeController) CreateHomeMate(w http.ResponseWriter, r *http.Request) {
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
		return
	}

	_, h, err := c.VerifyUserForAdminAction(slug, w, r)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.CreateHomeMates.VerifyUserForAdminAction")
		return
	}

	createHomeMates, ok := utils.DecodeBodyJSON[CreateHomeMatesDTO](w, r)
	if !ok {
		return
	}

	u, err := c.userLogic.FindOrCreate(0, createHomeMates.Name, createHomeMates.Email)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.CreateHomeMates.FindOrCreate")
		return
	}

	err = c.homeMateLogic.Create(h.ID, u.Id, homemate.Mate)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.CreateHomeMates.Create")
		return
	}

	w.Write([]byte("ok"))
}

func (c HomeController) DeleteHomeMate(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug path param", http.StatusBadRequest)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing email query param", http.StatusBadRequest)
	}

	_, h, err := c.VerifyUserForAdminAction(slug, w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	u, err := c.userLogic.FindByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Could not find mate", http.StatusBadRequest)
		} else {
			utils.RecordErrorServer(w, err, "Home.DeleteHomeMate.FindByEmail")
			return
		}
	}

	err = c.homeMateLogic.Delete(h.ID, u.Id)
	if err != nil {
		utils.RecordErrorServer(w, err, "Home.DeleteHomeMate.Delete")
	}

	w.Write([]byte("ok"))
}

func (c HomeController) VerifyUserForAdminAction(slug string, w http.ResponseWriter, r *http.Request) (user.User, Home, error) {
	var u user.User
	var h Home
	var ok bool
	h, err := c.logic.Read(slug)
	if err != nil {
		return u, h, err
	}

	u, ok = user.RetrieveUserFromContext(w, r)
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
