package user

import (
	"context"

	"github.com/rs/zerolog/log"

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

func (u userService) GetUserByTelegramID(ID int64) (*entity.User, error) {
	ctx := context.Background()
	user, err := u.userRepository.FindByTelegramID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) getUserByID(ID int) (*entity.User, error) {
	ctx := context.Background()

	user, err := u.userRepository.FindByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) CreateUser(user *entity.User) (*entity.User, error) {
	log.Info().Msg("creating user")

	ctx := context.Background()
	userID, err := u.userRepository.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return u.getUserByID(userID)
}
