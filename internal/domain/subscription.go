package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SubscriptionInfo — запись о подписке пользователя.
type SubscriptionInfo struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

// SubsFilter — фильтр для списка и подсчёта суммы.
// Поля-указатели опциональны: nil означает «не фильтровать».
// Limit/Offset используются только в List для пагинации.
type SubsFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	Limit       int
	Offset      int
}

// SubsInfoRepository — доступ к хранилищу подписок.
type SubsInfoRepository interface {
	Create(ctx context.Context, sub SubscriptionInfo) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SubscriptionInfo, error)
	Update(ctx context.Context, sub SubscriptionInfo) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubsFilter) ([]SubscriptionInfo, error)
	SumPrice(ctx context.Context, filter SubsFilter) (int, error)
}

// SubscriptionAggregatorService — бизнес-логика сервиса.
type SubscriptionAggregatorService interface {
	Create(ctx context.Context, sub SubscriptionInfo) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SubscriptionInfo, error)
	Update(ctx context.Context, sub SubscriptionInfo) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubsFilter) ([]SubscriptionInfo, error)
	SumPrice(ctx context.Context, filter SubsFilter) (int, error)
}
