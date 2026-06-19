package usecase

import (
	"context"
	"subAggregator/internal/domain"

	"github.com/google/uuid"
)

type SubAggregatorService struct {
	repo domain.SubsInfoRepository
}

func NewSubAggregatorService(repo domain.SubsInfoRepository) *SubAggregatorService {
	return &SubAggregatorService{
		repo: repo,
	}
}

func (s *SubAggregatorService) List(ctx context.Context, filter domain.SubsFilter) ([]domain.SubcriptionInfo, error) {
	subsInfo, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return subsInfo, nil
}
func (s *SubAggregatorService) SumPrice(ctx context.Context, filter domain.SubsFilter) (int, error) {
	sum, err := s.repo.SumPrice(ctx, filter)
	if err != nil {
		return 0, err
	}
	return sum, err
}

func (s *SubAggregatorService) Create(ctx context.Context, sub domain.SubcriptionInfo) (uuid.UUID, error)
func (s *SubAggregatorService) GetByID(ctx context.Context, id uuid.UUID) (*domain.SubcriptionInfo, error)
func (s *SubAggregatorService) Update(ctx context.Context, sub domain.SubcriptionInfo) error
func (s *SubAggregatorService) Delete(ctx context.Context, id uuid.UUID) error
