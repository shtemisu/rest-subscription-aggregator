package controller

import (
	"fmt"
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
