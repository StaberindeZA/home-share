package auth

type UserInfoDTO struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type VerifyRoleDTO struct {
	HomeSlug string `json:"homeSlug"`
	Role     string `json:"role"`
}
