package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"subAggregator/internal/domain"
	"subAggregator/pkg/response"

	"github.com/google/uuid"
)

type SubsHandler struct {
	SubsService domain.SubscriptionAggregatorService
	log         *slog.Logger
}

func NewSubsHandler(ss domain.SubscriptionAggregatorService, log *slog.Logger) *SubsHandler {
	return &SubsHandler{SubsService: ss, log: log}
}

// Create godoc
// @Summary      Создать подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input  body      SubInfoRequest  true  "Данные подписки"
// @Success      201    {object}  map[string]string
// @Failure      400    {object}  response.ErrorResponse
// @Failure      500    {object}  response.ErrorResponse
// @Router       /subscriptions [post]
func (h *SubsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var subReq SubInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}
	sub, err := subReq.ToDomain()
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.SubsService.Create(r.Context(), sub)
	if err != nil {
		h.log.Error("failed to create subscription", "error", err, "user_id", sub.UserID)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Info("sub was created", "id", id.String())
	response.WriteJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

// GetById godoc
// @Summary      Получить подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки (UUID)"
// @Success      200  {object}  SubInfoResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *SubsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id, expected UUID")
		return
	}
	sub, err := h.SubsService.GetByID(r.Context(), id)

	if errors.Is(err, domain.ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		h.log.Error("failed to get subscription", "error", err)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, NewSubInfoResponse(*sub))
}

// Update godoc
// @Summary      Обновить подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id     path      string          true  "ID подписки (UUID)"
// @Param        input  body      SubInfoRequest  true  "Новые данные подписки"
// @Success      200    {object}  SubInfoResponse
// @Failure      400    {object}  response.ErrorResponse
// @Failure      404    {object}  response.ErrorResponse
// @Failure      500    {object}  response.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *SubsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id, expected UUID")
		return
	}

	var subReq SubInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	sub, err := subReq.ToDomain()
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub.ID = id

	err = h.SubsService.Update(r.Context(), sub)
	if errors.Is(err, domain.ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		h.log.Error("failed to update subscription", "error", err, "id", sub.ID)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Info("sub was updated", "id", id.String())
	response.WriteJSON(w, http.StatusOK, NewSubInfoResponse(sub))
}

// Delete godoc
// @Summary      Удалить подписку
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *SubsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id, expected UUID")
		return
	}
	err1 := h.SubsService.Delete(r.Context(), id)
	if errors.Is(err1, domain.ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	if err1 != nil {
		h.log.Error("failed to delete subscription", "error", err1, "id", id)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

// List godoc
// @Summary      Список подписок
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "Фильтр по пользователю (UUID)"
// @Param        service_name  query     string  false  "Фильтр по названию сервиса"
// @Param        from          query     string  false  "Начало периода (MM-YYYY)"
// @Param        to            query     string  false  "Конец периода (MM-YYYY)"
// @Param        limit         query     int     false  "Размер страницы (по умолчанию 50, максимум 100)"
// @Param        offset        query     int     false  "Смещение"
// @Success      200  {array}   SubInfoResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /subscriptions [get]
func (h *SubsHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r.URL.Query())
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	subs, err := h.SubsService.List(r.Context(), filter)
	if err != nil {
		h.log.Error("failed to get list", "error", err)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]SubInfoResponse, 0, len(subs))
	for _, s := range subs {
		resp = append(resp, NewSubInfoResponse(s))
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

// SumUsingFilter godoc
// @Summary      Суммарная стоимость подписок за период
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "Фильтр по пользователю (UUID)"
// @Param        service_name  query     string  false  "Фильтр по названию сервиса"
// @Param        from          query     string  false  "Начало периода (MM-YYYY)"
// @Param        to            query     string  false  "Конец периода (MM-YYYY)"
// @Success      200  {object}  map[string]int
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /subscriptions/summary [get]
func (h *SubsHandler) SumUsingFilter(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r.URL.Query())
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	sum, err := h.SubsService.SumPrice(r.Context(), filter)
	if err != nil {
		h.log.Error("failed to get sum using filter", "error", err)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]int{"total": sum})
}
