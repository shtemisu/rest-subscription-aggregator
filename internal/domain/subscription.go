package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SubcriptionInfo — запись о подписке пользователя.
type SubcriptionInfo struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

// SubsFilter — фильтр для списка и подсчёта суммы.
// Все поля опциональны: nil означает «не фильтровать».
type SubsFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
}

// SubsInfoRepository — доступ к хранилищу подписок.
type SubsInfoRepository interface {
	Create(ctx context.Context, sub SubcriptionInfo) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SubcriptionInfo, error)
	Update(ctx context.Context, sub SubcriptionInfo) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubsFilter) ([]SubcriptionInfo, error)
	SumPrice(ctx context.Context, filter SubsFilter) (int, error)
}

// SubcriptionAggregatorService — бизнес-логика сервиса.
type SubcriptionAggregatorService interface {
	Create(ctx context.Context, sub SubcriptionInfo) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SubcriptionInfo, error)
	Update(ctx context.Context, sub SubcriptionInfo) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubsFilter) ([]SubcriptionInfo, error)
	SumPrice(ctx context.Context, filter SubsFilter) (int, error)
}
