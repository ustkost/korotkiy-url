package service

import (
	"context"

	"github.com/ustkost/korotkiy-url/internal/model"
	"github.com/ustkost/korotkiy-url/internal/repository"
)

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

func (s *ClickService) ListByLinkID(ctx context.Context, linkID int64, limit, offset int) ([]model.Click, error) {
	return s.repo.ListByLinkID(ctx, linkID, limit, offset)
}
