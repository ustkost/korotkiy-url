package service

import (
	"context"

	"github.com/ustkost/korotkiy-url/internal/model"
	"github.com/ustkost/korotkiy-url/internal/repository"
)

type ClickList struct {
	Clicks []model.Click `json:"clicks"`
	Total  int64         `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type ClickService struct {
	repo *repository.ClickRepository
}

func NewClickService(repo *repository.ClickRepository) *ClickService {
	return &ClickService{repo: repo}
}

func (s *ClickService) RecordClick(ctx context.Context, linkID int64, referrer string) error {
	click := &model.Click{
		LinkID:   linkID,
		Referrer: referrer,
	}
	return s.repo.Create(ctx, click)
}

func (s *ClickService) ListByLinkID(ctx context.Context, linkID int64, limit, offset int) (*ClickList, error) {
	if limit <= 0 || limit > maxListLimit {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}

	clicks, err := s.repo.ListByLinkID(ctx, linkID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByLinkID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	return &ClickList{
		Clicks: clicks,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
