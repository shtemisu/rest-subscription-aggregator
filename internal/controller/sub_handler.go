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

// Create
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

// Read one
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

// Update
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

// Delete
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

// List
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

// Sum
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
