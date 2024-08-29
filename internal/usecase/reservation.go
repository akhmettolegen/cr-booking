package usecase

import (
	"context"
	"github.com/akhmettolegen/cr-booking/internal/entity"
	"github.com/akhmettolegen/cr-booking/internal/repo"
)

type ReservationUsecase struct {
	repo repo.ReservationRepo
}

func New(repo repo.ReservationRepo) *ReservationUsecase {
	return &ReservationUsecase{repo}
}

func (u *ReservationUsecase) Create(ctx context.Context, reservation entity.Reservation) error {
	return u.repo.Store(ctx, reservation)
}

func (u *ReservationUsecase) GetByRoomId(ctx context.Context, roomId string) ([]entity.Reservation, error) {
	return u.repo.GetByRoomId(ctx, roomId)
}
