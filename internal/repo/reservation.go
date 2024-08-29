package repo

import (
	"context"
	"github.com/akhmettolegen/cr-booking/internal/entity"
	"github.com/akhmettolegen/cr-booking/pkg/postgres"
)

type ReservationStorage struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *ReservationStorage {
	return &ReservationStorage{pg}
}

func (r *ReservationStorage) GetByRoomId(ctx context.Context, roomId string) ([]entity.Reservation, error) {
	sql, _, err := r.Builder.
		Select("id, room_id, start_time, end_time").
		From("reservations").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reservations []entity.Reservation
	for rows.Next() {
		r := entity.Reservation{}
		err = rows.Scan(&r.Id, &r.RoomId, &r.StartTime, &r.EndTime)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, r)
	}

	return reservations, nil
}

func (r *ReservationStorage) Store(ctx context.Context, reservation entity.Reservation) error {
	sql, args, err := r.Builder.
		Insert("reservations").
		Columns("room_id, start_time, end_time").
		Values(reservation.RoomId, reservation.StartTime, reservation.EndTime).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil

}
