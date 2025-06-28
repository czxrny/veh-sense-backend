package models

type UserAuth struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRegisterInfo struct {
	UserName       string `json:"username" validate:"required,min=2"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	OrganizationId *int   `json:"organization_id,omitempty"`
}

type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserTokenResponse struct {
	Token   string `json:"token"`
	LocalId int    `json:"localId"`
}
