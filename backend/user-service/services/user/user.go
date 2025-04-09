package user

import (
	"context"
	"time"
	"user-service/config"
	"user-service/domain/dto"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
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
	panic("unimplemented")
}

func (u *UserService) GetUserLogin(context.Context) (*dto.UserResponse, error) {
	panic("unimplemented")
}

func (us *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := us.repository.GetUser().FindByEmail(ctx, req.Username)
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

func (u *UserService) Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	panic("unimplemented")
}

func (u *UserService) Update(context.Context, *dto.UpdateRequest) (*dto.UserResponse, error) {
	panic("unimplemented")
}
