package user

import (
	"context"
	"fmt"
	"net/http"

	s "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/service"
	common "github.com/czxrny/veh-sense-backend/shared/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type UserAuthHandler struct {
	*s.UserService
}

func NewUserAuthHandler(userService *s.UserService) *UserAuthHandler {
	return &UserAuthHandler{
		UserService: userService,
	}
}

func (uh *UserAuthHandler) RegisterPrivateUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) (*models.UserTokenResponse, error) {
		return uh.UserService.RegisterUser(ctx, userRegisterInfo, nil, "user")
	})
}

func (uh *UserAuthHandler) RegisterCorporateUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandlerSilent(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok || authClaims.Role != "admin" {
			return fmt.Errorf("Error: to create an organization user, login as an admin and pass the JWT!")
		}

		_, err := uh.UserService.RegisterUser(ctx, userRegisterInfo, authClaims.OrganizationID, "user")
		return err
	})
}

func (uh *UserAuthHandler) RegisterUserRoot(w http.ResponseWriter, r *http.Request) {
	common.PostHandlerSilent(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfoRoot) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok || authClaims.Role != "root" {
			return fmt.Errorf("Error: to create a custom user, login as a root and pass the JWT!")
		}

		userInfo := models.UserRegisterInfo{
			UserName: userRegisterInfo.UserName,
			Email:    userRegisterInfo.Email,
			Password: userRegisterInfo.Password,
		}

		_, err := uh.UserService.RegisterUser(ctx, &userInfo, userRegisterInfo.OrganizationID, userRegisterInfo.Role)
		return err
	})
}

func (uh *UserAuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userCredentials *models.UserCredentials) (*models.UserTokenResponse, error) {
		return uh.UserService.LoginUser(ctx, userCredentials)
	})
}

// Requires the user to login and pass updated information
func (uh *UserAuthHandler) UpdateLoginCredentials(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, credUpdateRequest *models.UserCredentialsUpdateRequest) (*models.UserTokenResponse, error) {
		if credUpdateRequest.NewEmail == "" || credUpdateRequest.NewPassword == "" {
			return nil, fmt.Errorf("User should pass email or password, or both to update")
		}

		return uh.UserService.UpdateLoginCredentials(ctx, credUpdateRequest)
	})
}

func (uh *UserAuthHandler) RefreshByKey(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, request *models.TokenRefreshRequest) (*models.UserTokenResponse, error) {
		return uh.UserService.GetRefreshToken(ctx, *request)
	})
}
