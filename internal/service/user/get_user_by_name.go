package user

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *userService) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	user, err := s.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}
