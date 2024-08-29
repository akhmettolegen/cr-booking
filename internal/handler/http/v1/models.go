package v1

import (
	"errors"
	"github.com/akhmettolegen/cr-booking/internal/entity"
	"time"
)

type reservationCreateRequest struct {
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func toReservation(in reservationCreateRequest) entity.Reservation {
	return entity.Reservation{
		RoomId:    in.RoomId,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
	}
}

func (r *reservationCreateRequest) validate() error {
	if r.RoomId == "" {
		return errors.New("room_id is required")
	}

	return nil
}

type reservationResponse struct {
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func toReservationResponse(in entity.Reservation) reservationResponse {
	return reservationResponse{
		RoomId:    in.RoomId,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
	}
}

func toReservationsResponse(in []entity.Reservation) []reservationResponse {
	res := make([]reservationResponse, 0, len(in))
	for _, r := range in {
		res = append(res, toReservationResponse(r))
	}

	return res
}
