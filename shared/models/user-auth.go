package models

type UserAuth struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRegisterInfoRoot struct {
	UserName       string `json:"user_name" validate:"required,min=2"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	OrganizationID *int   `json:"organization_id,omitempty"`
	Role           string `json:"role" validate:"required"`
}

type UserRegisterInfo struct {
	UserName string `json:"user_name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserCredentialsUpdateRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	NewEmail    string `json:"new_email" validate:"email"`
	NewPassword string `json:"new_password" validate:"min=6"`
}

type UserTokenResponse struct {
	Token      string `json:"token"`
	RefreshKey string `json:"refresh_key"`
	LocalId    int    `json:"localId"`
}

type AuthInfo struct {
	UserID         int
	Role           string
	OrganizationID *int
}
