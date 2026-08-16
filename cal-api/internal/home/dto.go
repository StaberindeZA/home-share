package home

type HomeDTO struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type HomeCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type HomeMatesDTO struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type CreateHomeMatesDTO struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
