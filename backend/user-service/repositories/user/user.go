package user

import (
	"context"
	"errors"
	customErr "user-service/common/custom-error"
	errConstant "user-service/constants/custom-error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := ur.db.
		WithContext(ctx).
		Model(&models.User{}).
		Preload("Role").
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}

		return nil, customErr.WrapError(errConstant.ErrSQL)
	}

	return &user, nil
}

func (ur *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := ur.db.
		WithContext(ctx).
		Model(&models.User{}).
		Preload("Role").
		Where("uuid = ?", uuid).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}

		return nil, customErr.WrapError(errConstant.ErrSQL)
	}

	return &user, nil
}

func (ur *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := ur.db.
		WithContext(ctx).
		Model(&models.User{}).
		Preload("Role").
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}

		return nil, customErr.WrapError(errConstant.ErrSQL)
	}

	return &user, nil
}

func (ur *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		RoleID:      req.RoleID,
	}

	err := ur.db.
		WithContext(ctx).
		Model(&models.User{}).
		Create(&user).
		Error
	if err != nil {
		return nil, customErr.WrapError(errConstant.ErrSQL)
	}

	return &user, err
}

func (ur *UserRepository) Update(ctx context.Context, req *dto.UpdateRequest, userUuid string) (*models.User, error) {
	user := models.User{
		UUID:        uuid.MustParse(userUuid),
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    *req.Password,
		PhoneNumber: req.PhoneNumber,
		RoleID:      req.RoleID,
	}

	err := ur.db.
		WithContext(ctx).
		Where("uuid = ?", userUuid).
		Updates(&user).
		Error
	if err != nil {
		return nil, customErr.WrapError(errConstant.ErrSQL)
	}

	return &user, err
}
