package user

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	s "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/service"
	common "github.com/czxrny/veh-sense-backend/shared/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type UserInfoHandler struct {
	*s.UserService
}

func NewUserInfoHandler(userService *s.UserService) *UserInfoHandler {
	return &UserInfoHandler{
		UserService: userService,
	}
}

// For admins
func (uh *UserInfoHandler) GetAllUsersInfo(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.UserInfo, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return uh.UserService.GetAllUsersInfo(ctx, authClaims)
	})
}

func (uh *UserInfoHandler) GetMyUserInfo(w http.ResponseWriter, r *http.Request) {
	common.GetSimpleHandler(w, r, func(ctx context.Context) (*models.UserInfo, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return uh.UserService.GetMyUserInfo(ctx, authClaims)
	})
}

// Must be either owner, admin of the user org
func (uh *UserInfoHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return uh.UserService.DeleteUser(ctx, authClaims, id)
	})
}
