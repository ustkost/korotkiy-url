package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/ustkost/korotkiy-url/internal/model"
)

type LinkRepository struct {
	pool *pgxpool.Pool
}

func NewLinkRepository(pool *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{pool: pool}
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	query := `
		INSERT INTO links (short_code, original_url)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, link.ShortCode, link.OriginalURL).
		Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateCode
		}
		return err
	}
	return nil
}

func (r *LinkRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	query := `
		SELECT id, short_code, original_url, created_at
		FROM links
		WHERE short_code = $1
	`
	var link model.Link
	err := r.pool.QueryRow(ctx, query, shortCode).
		Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) List(ctx context.Context, limit, offset int) ([]model.Link, error) {
	query := `
		SELECT id, short_code, original_url, created_at
		FROM links
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (r *LinkRepository) UpdateOriginalURL(ctx context.Context, id int64, originalURL string) error {
	query := `UPDATE links SET original_url = $1 WHERE id = $2`
	tag, err := r.pool.Exec(ctx, query, originalURL, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *LinkRepository) UpdateShortCode(ctx context.Context, id int64, shortCode string) error {
	query := `UPDATE links SET short_code = $1 WHERE id = $2`
	tag, err := r.pool.Exec(ctx, query, shortCode, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return ErrDuplicateCode
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM links WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
