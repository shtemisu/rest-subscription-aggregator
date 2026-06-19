package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"subAggregator/internal/domain"
	"subAggregator/pkg/response"

	"github.com/google/uuid"
)

type SubsHandler struct {
	SubsService domain.SubcriptionAggregatorService
}

func NewSubsHandler(ss domain.SubcriptionAggregatorService) *SubsHandler {
	return &SubsHandler{
		SubsService: ss,
	}
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
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusCreated, fmt.Sprintf("new sub with id:%s was created", id.String()))
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
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, "subscription updated")
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
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, fmt.Sprintf("sub with id: %s was deleted", id.String()))
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
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]int{"total": sum})
}
