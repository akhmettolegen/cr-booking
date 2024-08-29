package repo

import (
	"context"
	"github.com/akhmettolegen/cr-booking/internal/entity"
)

type ReservationRepo interface {
	Store(ctx context.Context, reservation entity.Reservation) error
	GetByRoomId(ctx context.Context, roomId string) ([]entity.Reservation, error)
}
