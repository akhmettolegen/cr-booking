package v1

import (
	"encoding/json"
	"fmt"
	"github.com/akhmettolegen/cr-booking/internal/usecase"
	"github.com/akhmettolegen/cr-booking/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type ReservationHandler struct {
	ReservationUsecase usecase.Reservation
	logger             logger.Interface
}

func NewReservationHandler(logger logger.Interface, reservationUsecase usecase.Reservation) *ReservationHandler {
	return &ReservationHandler{
		ReservationUsecase: reservationUsecase,
		logger:             logger,
	}
}

func (h *ReservationHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", h.create)
	})

	return r
}

func (h *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req reservationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err, "")
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Error(err, "")
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	err := h.ReservationUsecase.Create(ctx, toReservation(req))
	if err != nil {
		h.logger.Error(err, "")
		_ = render.Render(w, r, errResponse(http.StatusInternalServerError, "internal server error"))
		return
	}
}

func (h *ReservationHandler) GetByRoomId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roomId := chi.URLParam(r, "roomId")
	if roomId == "" {
		h.logger.Error(fmt.Errorf("roomId is empty"))
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	result, err := h.ReservationUsecase.GetByRoomId(ctx, roomId)
	if err != nil {
		_ = render.Render(w, r, errResponse(http.StatusInternalServerError, "internal server error"))
		return
	}

	render.JSON(w, r, toReservationsResponse(result))
}
