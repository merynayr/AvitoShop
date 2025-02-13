package shop

import (
	"context"
)

func (s *shopService) GetMerchPrice(ctx context.Context, item string) (int64, error) {
	user, err := s.shopRepository.GetMerchPrice(ctx, item)
	if err != nil {
		return 0, err
	}

	return user, nil
}
