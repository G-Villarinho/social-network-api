package repository

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	di          *internal.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(di *internal.Di) (domain.UserRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := internal.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		di:          di,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user domain.User) error {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrUserNotFound
		}
		return err
	}

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user domain.User) error {
	if err := u.db.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) DeleteUser(ctx context.Context, ID uuid.UUID) error {
	if err := u.db.WithContext(ctx).Where("id = ?", ID).Delete(&domain.User{}).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) CheckUsername(ctx context.Context, username string) (bool, error) {
	var exists bool

	err := u.db.WithContext(ctx).
		Raw("SELECT EXISTS(SELECT 1 FROM User WHERE username = ?)", username).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}
