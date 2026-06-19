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

func (s *SubAggregatorService) Create(ctx context.Context, sub domain.SubcriptionInfo) (uuid.UUID, error) {
	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (s *SubAggregatorService) GetByID(ctx context.Context, id uuid.UUID) (*domain.SubcriptionInfo, error) {
	info, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (s *SubAggregatorService) Update(ctx context.Context, sub domain.SubcriptionInfo) error {
	return s.repo.Update(ctx, sub)
}
func (s *SubAggregatorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
