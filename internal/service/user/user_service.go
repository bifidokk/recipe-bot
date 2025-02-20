package user

import (
	"context"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/repository"
	"github.com/bifidokk/recipe-bot/internal/service"
)

type userService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository,
) service.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u userService) GetUser(ID int64) (*entity.User, error) {
	ctx := context.Background()
	user, err := u.userRepository.FindByTelegramID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) CreateUser(_ *entity.User) error {
	panic("implement me")
}
