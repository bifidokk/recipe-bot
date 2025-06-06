package user

import (
	"context"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/repository/user"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/rs/zerolog/log"
)

const recipeLimitForNewUser = 5

type userService struct {
	userRepository *user.Repository
}

func NewUserService(
	userRepository *user.Repository,
) service.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u userService) GetUserByTelegramID(ID int64) (*entity.User, error) {
	ctx := context.Background()
	usr, err := u.userRepository.FindByTelegramID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (u userService) getUserByID(ID int) (*entity.User, error) {
	ctx := context.Background()

	usr, err := u.userRepository.FindByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (u userService) CreateUser(user *entity.User) (*entity.User, error) {
	log.Info().Msg("creating user")

	user.RecipeLimit = recipeLimitForNewUser

	ctx := context.Background()
	userID, err := u.userRepository.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return u.getUserByID(userID)
}

func (u userService) DecreaseUserLimit(user *entity.User) error {
	user.RecipeLimit--

	if user.RecipeLimit < 0 {
		return nil
	}

	ctx := context.Background()
	err := u.userRepository.UpdateUser(ctx, user)

	if err != nil {
		return err
	}

	return nil
}
