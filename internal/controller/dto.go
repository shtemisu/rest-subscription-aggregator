package controller

import (
	"fmt"
	"net/url"
	"subAggregator/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SubInfoRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubInfoResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func NewSubInfoResponse(s domain.SubcriptionInfo) SubInfoResponse {
	resp := SubInfoResponse{
		ID:          s.ID.String(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   s.StartDate.Format(monthYearLayout),
	}
	if s.EndDate != nil {
		end := s.EndDate.Format(monthYearLayout)
		resp.EndDate = &end
	}
	return resp
}

const monthYearLayout = "01-2006"

func (r *SubInfoRequest) ToDomain() (domain.SubcriptionInfo, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return domain.SubcriptionInfo{}, fmt.Errorf("invalid user_id: %w", err)
	}

	start, err := time.Parse(monthYearLayout, r.StartDate)
	if err != nil {
		return domain.SubcriptionInfo{}, fmt.Errorf("invalid start_date, expected MM-YYYY: %w", err)
	}
	if r.ServiceName == "" {
		return domain.SubcriptionInfo{}, fmt.Errorf("invalid service_name, it must not be empty")
	}
	if r.Price < 0 {
		return domain.SubcriptionInfo{}, fmt.Errorf("invalid price, it must not be lower than 0")
	}
	sub := domain.SubcriptionInfo{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   start,
	}

	if r.EndDate != nil && *r.EndDate != "" {
		end, err := time.Parse(monthYearLayout, *r.EndDate)
		if err != nil {
			return domain.SubcriptionInfo{}, fmt.Errorf("invalid end_date, expected MM-YYYY: %w", err)
		}
		sub.EndDate = &end
	}

	return sub, nil
}

// parseFilter собирает domain.SubsFilter из query-параметров.
// Поддерживает user_id, service_name и период from/to (даты в формате MM-YYYY).
// Любое отсутствующее поле остаётся nil — фильтрация по нему не применяется.
func parseFilter(q url.Values) (domain.SubsFilter, error) {
	var f domain.SubsFilter

	if v := q.Get("user_id"); v != "" {
		userID, err := uuid.Parse(v)
		if err != nil {
			return f, fmt.Errorf("invalid user_id, expected UUID")
		}
		f.UserID = &userID
	}
	if v := q.Get("service_name"); v != "" {
		f.ServiceName = &v
	}
	if v := q.Get("from"); v != "" {
		t, err := time.Parse(monthYearLayout, v)
		if err != nil {
			return f, fmt.Errorf("invalid 'from', expected MM-YYYY")
		}
		f.PeriodStart = &t
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse(monthYearLayout, v)
		if err != nil {
			return f, fmt.Errorf("invalid 'to', expected MM-YYYY")
		}
		f.PeriodEnd = &t
	}

	return f, nil
}
