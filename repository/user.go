package repository

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	di          *pkg.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(di *pkg.Di) (domain.UserRepository, error) {
	db, err := pkg.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.Invoke[*redis.Client](di)
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

func (u *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User

	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
