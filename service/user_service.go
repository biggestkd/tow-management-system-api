package service

import (
	"context"
	"fmt"
	"time"

	"tow-management-system-api/model"
	"tow-management-system-api/repository"
)

type UserRepository interface {
	Create(ctx context.Context, item *model.User) error
	Find(ctx context.Context, filter *model.User) ([]*model.User, error)
	Update(ctx context.Context, id string, updateData *model.User) error
}

// UserService has ONLY the required dep per spec.
type UserService struct {
	userRepository repository.Repository[model.User]
}

func NewUserServiceWithMongo(userRepo *repository.UserMongoRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

// CreateUser returns success/failure (bool) per spec.
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("user payload is nil")
	}

	user.CreatedDate = time.Now().UTC().Unix()

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	return nil
}

// UpdateUser returns success/failure to match the pattern.
func (s *UserService) UpdateUser(ctx context.Context, userId *string, user *model.User) error {

	if user == nil {
		return fmt.Errorf("user payload is nil")
	}

	if err := s.userRepository.Update(ctx, *userId, user); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (s *UserService) FindUserById(ctx context.Context, user *model.User) (*model.User, error) {

	if user == nil {
		return nil, fmt.Errorf("user payload is nil")
	}

	users, err := s.userRepository.Find(ctx, user)

	if err != nil || len(users) < 1 {
		return nil, fmt.Errorf("find user failed: %w", err)
	}

	return users[0], nil
}
