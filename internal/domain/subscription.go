package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubcriptionInfo struct {
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
}

type SubsInfoRepository interface {
	GetSubsInfoByUserId(userID uuid.UUID)
}

type SubcriptionAggregatorService interface {
	GetSubsInfoByUserId(userID uuid.UUID) (*SubcriptionInfo, error)
	GetSubsInfoList() (*[]SubcriptionInfo, error)
}
