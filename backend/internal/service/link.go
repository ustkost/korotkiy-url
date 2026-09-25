package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"regexp"

	"github.com/ustkost/korotkiy-url/internal/model"
	"github.com/ustkost/korotkiy-url/internal/repository"
)

var (
	ErrInvalidShortCode     = errors.New("short code must be 1-16 characters, alphanumeric or underscore")
	ErrInvalidURL           = errors.New("original URL must be a valid http or https URL")
	ErrCodeGenerationFailed = errors.New("failed to generate a unique short code after several attempts")
	shortCodePattern        = regexp.MustCompile(`^[a-zA-Z0-9_]{1,16}$`)
)

const (
	autoCodeLength   = 8
	maxCreateRetries = 10
	codeCharset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type LinkList struct {
	Links  []model.Link `json:"links"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type LinkService struct {
	repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func validateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return ErrInvalidURL
	}
	return nil
}

func (s *LinkService) Create(ctx context.Context, originalURL, customCode string) (*model.Link, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}
	if customCode != "" {
		return s.createWithCustomCode(ctx, originalURL, customCode)
	}
	return s.createWithGeneratedCode(ctx, originalURL)
}

func (s *LinkService) createWithCustomCode(ctx context.Context, originalURL, customCode string) (*model.Link, error) {
	if !shortCodePattern.MatchString(customCode) {
		return nil, ErrInvalidShortCode
	}

	link := &model.Link{
		ShortCode:   customCode,
		OriginalURL: originalURL,
	}
	if err := s.repo.Create(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func withGeneratedCode(fn func(code string) (*model.Link, error)) (*model.Link, error) {
	for range maxCreateRetries {
		code, err := generateShortCode(autoCodeLength)
		if err != nil {
			return nil, err
		}
		link, err := fn(code)
		if err == nil {
			return link, nil
		}
		if !errors.Is(err, repository.ErrDuplicateCode) {
			return nil, err
		}
	}
	return nil, ErrCodeGenerationFailed
}

func (s *LinkService) createWithGeneratedCode(ctx context.Context, originalURL string) (*model.Link, error) {
	return withGeneratedCode(func(code string) (*model.Link, error) {
		link := &model.Link{
			ShortCode:   code,
			OriginalURL: originalURL,
		}
		if err := s.repo.Create(ctx, link); err != nil {
			return nil, err
		}
		return link, nil
	})
}

func generateShortCode(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeCharset))))
		if err != nil {
			return "", err
		}
		b[i] = codeCharset[n.Int64()]
	}
	return string(b), nil
}

func (s *LinkService) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	return s.repo.GetByShortCode(ctx, shortCode)
}

func (s *LinkService) List(ctx context.Context, limit, offset int) (*LinkList, error) {
	if limit <= 0 || limit > maxListLimit {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}

	links, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &LinkList{
		Links:  links,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LinkService) UpdateOriginalURL(ctx context.Context, id int64, originalURL string) (*model.Link, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}
	return s.repo.UpdateOriginalURL(ctx, id, originalURL)
}

func (s *LinkService) UpdateShortCode(ctx context.Context, id int64, customCode string) (*model.Link, error) {
	if customCode != "" {
		return s.updateWithCustomCode(ctx, id, customCode)
	}
	return s.updateWithGeneratedCode(ctx, id)
}

func (s *LinkService) updateWithCustomCode(ctx context.Context, id int64, customCode string) (*model.Link, error) {
	if !shortCodePattern.MatchString(customCode) {
		return nil, ErrInvalidShortCode
	}
	return s.repo.UpdateShortCode(ctx, id, customCode)
}

func (s *LinkService) updateWithGeneratedCode(ctx context.Context, id int64) (*model.Link, error) {
	return withGeneratedCode(func(code string) (*model.Link, error) {
		return s.repo.UpdateShortCode(ctx, id, code)
	})
}

func (s *LinkService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
