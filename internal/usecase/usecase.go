package usecase

import (
	"context"
	"github.com/akhmettolegen/cr-booking/internal/entity"
)

type Reservation interface {
	Create(ctx context.Context, reservation entity.Reservation) error
	GetByRoomId(ctx context.Context, roomId string) ([]entity.Reservation, error)
}
