package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/ustkost/korotkiy-url/internal/model"
)

type ClickRepository struct {
	pool *pgxpool.Pool
}

func NewClickRepository(pool *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{pool: pool}
}

func (r *ClickRepository) Create(ctx context.Context, click *model.Click) error {
	query := `
		INSERT INTO clicks (link_id, referrer)
		VALUES ($1, $2)
		RETURNING id, timestamp
	`
	err := r.pool.QueryRow(ctx, query, click.LinkID, click.Referrer).
		Scan(&click.ID, &click.Timestamp)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *ClickRepository) ListByLinkID(ctx context.Context, linkID int64, limit, offset int) ([]model.Click, error) {
	query := `
		SELECT id, link_id, timestamp, referrer
		FROM clicks
		WHERE link_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, linkID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []model.Click
	for rows.Next() {
		var click model.Click
		if err := rows.Scan(&click.ID, &click.LinkID, &click.Timestamp, &click.Referrer); err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}
	return clicks, rows.Err()
}
