package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/repository"
	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repoAuth    *r.UserAuthRepository
	repoInfo    *r.UserInfoRepository
	repoRefresh *r.RefreshKeyRepository
}

func NewUserService(repoAuth *r.UserAuthRepository, repoInfo *r.UserInfoRepository, repoRefresh *r.RefreshKeyRepository) *UserService {
	return &UserService{repoAuth: repoAuth, repoInfo: repoInfo, repoRefresh: repoRefresh}
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
		NumberOfRides:   0,
	}

	if err := s.repoInfo.Add(ctx, &userInfo); err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(&newUser, &userInfo)
	if err != nil {
		return nil, err
	}

	refreshKey, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	s.repoRefresh.Add(ctx, &models.RefreshInfo{
		UserID:     userInfo.ID,
		RefreshKey: refreshKey,
		ExpiresAt:  time.Now().Add(time.Hour * 24 * 30),
	})

	return &models.UserTokenResponse{
		Token:      token,
		RefreshKey: refreshKey,
		LocalId:    newUser.ID,
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

	// delete old
	err = s.repoRefresh.DeleteByUserID(ctx, userInfo.ID)
	if err != nil {
		return nil, err
	}

	refreshKey, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	s.repoRefresh.Add(ctx, &models.RefreshInfo{
		UserID:     userInfo.ID,
		RefreshKey: refreshKey,
		ExpiresAt:  time.Now().Add(time.Hour * 24 * 30),
	})

	return &models.UserTokenResponse{
		Token:      token,
		RefreshKey: refreshKey,
		LocalId:    userAuth.ID,
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

	if err := s.repoRefresh.DeleteByUserID(ctx, userInfo.ID); err != nil {
		return nil, err
	}

	refreshKey, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	s.repoRefresh.Add(ctx, &models.RefreshInfo{
		UserID:     userInfo.ID,
		RefreshKey: refreshKey,
		ExpiresAt:  time.Now().Add(time.Hour * 24 * 30),
	})

	return &models.UserTokenResponse{
		Token:      token,
		RefreshKey: refreshKey,
		LocalId:    userAuth.ID,
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

func (s *UserService) GetUserInfoById(ctx context.Context, authInfo models.AuthInfo, id int) (*models.UserInfo, error) {
	userInfo, err := s.repoInfo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isOwner := id == authInfo.UserID
	isOrgAdmin := authInfo.OrganizationID != nil && userInfo.OrganizationId != nil && *userInfo.OrganizationId == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isOwner && !isOrgAdmin && authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to get user's info.")
	}

	return userInfo, nil
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

	if err := s.repoRefresh.DeleteByUserID(ctx, id); err != nil {
		return err
	}

	if err := s.repoAuth.DeleteById(ctx, id); err != nil {
		return err
	}

	return s.repoInfo.DeleteById(ctx, id)
}

func (s *UserService) GetRefreshToken(ctx context.Context, refreshRequest models.TokenRefreshRequest) (*models.UserTokenResponse, error) {
	var refreshInfo *models.RefreshInfo
	var err error

	if refreshInfo, err = s.repoRefresh.FindByMatchingKey(ctx, refreshRequest.RefreshKey); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Invalid key")
		}
		return nil, err
	}

	if refreshInfo.UserID != refreshRequest.UserID {
		return nil, fmt.Errorf("Invalid key")
	}

	if refreshInfo.ExpiresAt.Before(time.Now()) {
		err = s.repoRefresh.DeleteById(ctx, refreshInfo.ID)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Key expired. Please relog to proceed.")
	}

	var newKey string
	if newKey, err = auth.GenerateRefreshToken(); err != nil {
		return nil, fmt.Errorf("Error while creating new key")
	}

	refreshInfo.RefreshKey = newKey
	refreshInfo.ExpiresAt = time.Now().Add(time.Hour * 24 * 30)

	if err = s.repoRefresh.Update(ctx, refreshInfo); err != nil {
		return nil, err
	}

	var userInfo *models.UserInfo
	if userInfo, err = s.repoInfo.GetByID(ctx, refreshInfo.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found!")
		}
		return nil, err
	}

	var userAuth *models.UserAuth
	if userAuth, err = s.repoAuth.GetByID(ctx, refreshInfo.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Not found!")
		}
		return nil, err
	}

	token, err := auth.CreateToken(userAuth, userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:      token,
		RefreshKey: refreshInfo.RefreshKey,
		LocalId:    refreshInfo.UserID,
	}, nil
}
