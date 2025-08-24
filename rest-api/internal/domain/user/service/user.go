package service

import (
	"context"
	"errors"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/repository"
	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repoAuth *r.UserAuthRepository
	repoInfo *r.UserInfoRepository
}

func NewUserService(repoAuth *r.UserAuthRepository, repoInfo *r.UserInfoRepository) *UserService {
	return &UserService{repoAuth: repoAuth, repoInfo: repoInfo}
}

func (s *UserService) RegisterUser(ctx context.Context, userRegisterInfo *models.UserRegisterInfo, organizationId *int, role string) (*models.UserTokenResponse, error) {
	var err error

	if _, err := s.repoAuth.FindByEmail(ctx, userRegisterInfo.Email); err == nil {
		return nil, fmt.Errorf("Email is already taken.")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	userRegisterInfo.Password, err = auth.EncryptThePassword(userRegisterInfo.Password)
	if err != nil {
		return nil, err
	}

	newUser := models.UserAuth{
		Email:    userRegisterInfo.Email,
		Password: userRegisterInfo.Password,
		Role:     role,
	}

	if err := s.repoAuth.Add(ctx, &newUser); err != nil {
		return nil, err
	}

	userInfo := models.UserInfo{
		ID:              newUser.ID,
		UserName:        userRegisterInfo.UserName,
		OrganizationId:  organizationId,
		TotalKilometers: 0,
	}

	if err := s.repoInfo.Add(ctx, &userInfo); err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(&newUser, &userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:   token,
		LocalId: newUser.ID,
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, userCredentials *models.UserCredentials) (*models.UserTokenResponse, error) {
	var err error
	var userAuth *models.UserAuth

	if userAuth, err = s.repoAuth.FindByEmail(ctx, userCredentials.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User does not exist.")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(userCredentials.Password)); err != nil {
		return nil, fmt.Errorf("Invalid login credentials.")
	}

	var userInfo *models.UserInfo
	if userInfo, err = s.repoInfo.GetByID(ctx, userAuth.ID); err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(userAuth, userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:   token,
		LocalId: userAuth.ID,
	}, nil
}

func (s *UserService) UpdateLoginCredentials(ctx context.Context, credUpdateRequest *models.UserCredentialsUpdateRequest) (*models.UserTokenResponse, error) {
	var err error
	var userAuth *models.UserAuth

	if userAuth, err = s.repoAuth.FindByEmail(ctx, credUpdateRequest.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User does not exist.")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(credUpdateRequest.Password)); err != nil {
		return nil, fmt.Errorf("Invalid login credentials.")
	}

	if _, err := s.repoAuth.FindByEmail(ctx, credUpdateRequest.NewEmail); err == nil {
		return nil, fmt.Errorf("Email is already taken.")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if credUpdateRequest.NewEmail != "" {
		userAuth.Email = credUpdateRequest.NewEmail
	}
	if credUpdateRequest.NewPassword != "" {
		var err error
		userAuth.Password, err = auth.EncryptThePassword(credUpdateRequest.NewPassword)
		if err != nil {
			return nil, err
		}
	}

	if err := s.repoAuth.Update(ctx, userAuth); err != nil {
		return nil, err
	}

	var userInfo *models.UserInfo
	if userInfo, err = s.repoInfo.GetByID(ctx, userAuth.ID); err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(userAuth, userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:   token,
		LocalId: userAuth.ID,
	}, nil
}

func (s *UserService) GetAllUsersInfo(ctx context.Context, authInfo models.AuthInfo) ([]models.UserInfo, error) {
	if authInfo.Role != "admin" {
		return nil, fmt.Errorf("Error: User is unauthorized to see other users.")
	}

	return s.repoInfo.FindByOrganizationId(ctx, *authInfo.OrganizationID)
}

func (s *UserService) GetMyUserInfo(ctx context.Context, authInfo models.AuthInfo) (*models.UserInfo, error) {
	return s.repoInfo.GetByID(ctx, authInfo.UserID)
}

func (s *UserService) DeleteUser(ctx context.Context, authInfo models.AuthInfo, id int) error {
	var userInfo *models.UserInfo
	var err error

	if userInfo, err = s.repoInfo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Not found!")
		}
		return err
	}

	isOwner := id == authInfo.UserID
	isOrgAdmin := authInfo.OrganizationID != nil && userInfo.OrganizationId != nil && *userInfo.OrganizationId == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isOwner && !isOrgAdmin && authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to delete the user.")
	}

	if err := s.repoAuth.DeleteById(ctx, id); err != nil {
		return err
	}
	return s.repoInfo.DeleteById(ctx, id)
}
