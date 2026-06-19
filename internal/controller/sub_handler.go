package controller

import "subAggregator/internal/domain"

type SubsHandler struct {
	SubsService domain.SubcriptionAggregatorService
}

func NewSubsHandler(ss domain.SubcriptionAggregatorService) *SubsHandler {
	return &SubsHandler{
		SubsService: ss,
	}
}
