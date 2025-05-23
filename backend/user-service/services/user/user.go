package user

import (
	"context"
	"time"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/custom-error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
	IsUsernameExist(context.Context, string) bool
	IsEmailExist(context.Context, string) bool
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (us *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := us.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Code,
	}

	return &data, nil
}

func (us *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Email:       userLogin.Email,
		Username:    userLogin.Username,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}

	return &data, nil
}

func (us *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := us.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	expiryTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Code,
	}
	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-service",
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiryTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (us *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if us.IsUsernameExist(ctx, req.Username) {
		return nil, errConstant.ErrUsernameExist
	}

	if us.IsEmailExist(ctx, req.Email) {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := us.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPass),
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role.Code,
		},
	}

	return response, nil
}

func (us *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password   string
		hashedPass []byte
		user       *models.User
		err        error
		data       dto.UserResponse
	)

	user, err = us.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	isExist := us.IsUsernameExist(ctx, req.Username)
	if isExist && user.Username != req.Username {
		return nil, errConstant.ErrUsernameExist
	}

	isExist = us.IsEmailExist(ctx, req.Email)
	if isExist && user.Email != req.Email {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			logrus.Infof("Password and ConfirmPassword did not match: %s - %s", *req.Password, *req.ConfirmPassword)
			return nil, errConstant.ErrPasswordDoesNotMatch
		}

		hashedPass, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
	}

	password = string(hashedPass)
	user, err = us.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    &password,
		PhoneNumber: req.PhoneNumber,
	}, uuid)
	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Email:       user.Email,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Code,
	}

	return &data, nil
}

func (us *UserService) IsUsernameExist(ctx context.Context, username string) bool {
	user, _ := us.repository.GetUser().FindByUsername(ctx, username)
	if user != nil {
		return true
	}

	return false
}

func (us *UserService) IsEmailExist(ctx context.Context, email string) bool {
	user, _ := us.repository.GetUser().FindByEmail(ctx, email)
	if user != nil {
		return true
	}

	return false
}
