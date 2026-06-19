package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"subAggregator/internal/domain"
)

type SubsRepo struct {
	pool *pgxpool.Pool
}

func NewSubsRepo(pool *pgxpool.Pool) *SubsRepo {
	return &SubsRepo{
		pool: pool,
	}
}

func (r *SubsRepo) Create(ctx context.Context, sub domain.SubcriptionInfo) (uuid.UUID, error) {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create subscription: %w", err)
	}

	return id, nil
}

func (r *SubsRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SubcriptionInfo, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1`

	var sub domain.SubcriptionInfo
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get subscription by id: %w", err)
	}

	return &sub, nil
}

func (r *SubsRepo) Update(ctx context.Context, sub domain.SubcriptionInfo) error {
	const query = `
		UPDATE subscriptions
		SET service_name = $2,
		    price        = $3,
		    user_id      = $4,
		    start_date   = $5,
		    end_date     = $6,
		    updated_at   = now()
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	)
	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *SubsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM subscriptions WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *SubsRepo) List(ctx context.Context, filter domain.SubsFilter) ([]domain.SubcriptionInfo, error) {
	query, args := buildFilter(`
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions`, filter)
	query += " ORDER BY start_date DESC, id"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make([]domain.SubcriptionInfo, 0)
	for rows.Next() {
		var sub domain.SubcriptionInfo
		if err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate,
		); err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions: %w", err)
	}

	return subs, nil
}

func (r *SubsRepo) SumPrice(ctx context.Context, filter domain.SubsFilter) (int, error) {
	query, args := buildFilter(`
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions`, filter)

	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("sum subscriptions price: %w", err)
	}

	return total, nil
}

func buildFilter(base string, f domain.SubsFilter) (string, []any) {
	conds := make([]string, 0, 4)
	args := make([]any, 0, 4)

	if f.UserID != nil {
		args = append(args, *f.UserID)
		conds = append(conds, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if f.ServiceName != nil {
		args = append(args, *f.ServiceName)
		conds = append(conds, fmt.Sprintf("service_name = $%d", len(args)))
	}
	if f.PeriodEnd != nil {
		args = append(args, *f.PeriodEnd)
		conds = append(conds, fmt.Sprintf("start_date <= $%d", len(args)))
	}
	if f.PeriodStart != nil {
		args = append(args, *f.PeriodStart)
		conds = append(conds, fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", len(args)))
	}

	if len(conds) == 0 {
		return base, args
	}

	return base + " WHERE " + strings.Join(conds, " AND "), args
}
