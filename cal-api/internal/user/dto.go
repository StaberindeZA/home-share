package user

type UpdateUserDTO struct {
	Name string `json:"name"`
}

type UserProfileDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
