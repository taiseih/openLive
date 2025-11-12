package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

type liveAreaRepository struct {
	db *sql.DB
}

// NewLiveAreaRepository は新しいライブエリアリポジトリを作成します
func NewLiveAreaRepository(db *sql.DB) domain.LiveAreaRepository {
	return &liveAreaRepository{db: db}
}

func (r *liveAreaRepository) Create(ctx context.Context, liveArea *domain.LiveArea) error {
	query := `
		INSERT INTO live_areas (id, owner_id, name, description, is_public, thumbnail_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		liveArea.ID, liveArea.OwnerID, liveArea.Name, liveArea.Description,
		liveArea.IsPublic, liveArea.ThumbnailURL, liveArea.CreatedAt, liveArea.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create live area: %w", err)
	}
	return nil
}

func (r *liveAreaRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LiveArea, error) {
	query := `
		SELECT id, owner_id, name, description, is_public, thumbnail_url, created_at, updated_at
		FROM live_areas WHERE id = $1
	`
	area := &domain.LiveArea{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&area.ID, &area.OwnerID, &area.Name, &area.Description,
		&area.IsPublic, &area.ThumbnailURL, &area.CreatedAt, &area.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("live area not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find live area: %w", err)
	}
	return area, nil
}

func (r *liveAreaRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.LiveArea, error) {
	query := `
		SELECT id, owner_id, name, description, is_public, thumbnail_url, created_at, updated_at
		FROM live_areas WHERE owner_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find live areas: %w", err)
	}
	defer rows.Close()

	var areas []*domain.LiveArea
	for rows.Next() {
		area := &domain.LiveArea{}
		if err := rows.Scan(
			&area.ID, &area.OwnerID, &area.Name, &area.Description,
			&area.IsPublic, &area.ThumbnailURL, &area.CreatedAt, &area.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan live area: %w", err)
		}
		areas = append(areas, area)
	}
	return areas, nil
}

func (r *liveAreaRepository) FindPublicAreas(ctx context.Context, limit, offset int) ([]*domain.LiveArea, error) {
	query := `
		SELECT id, owner_id, name, description, is_public, thumbnail_url, created_at, updated_at
		FROM live_areas WHERE is_public = true ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find public areas: %w", err)
	}
	defer rows.Close()

	var areas []*domain.LiveArea
	for rows.Next() {
		area := &domain.LiveArea{}
		if err := rows.Scan(
			&area.ID, &area.OwnerID, &area.Name, &area.Description,
			&area.IsPublic, &area.ThumbnailURL, &area.CreatedAt, &area.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan live area: %w", err)
		}
		areas = append(areas, area)
	}
	return areas, nil
}

func (r *liveAreaRepository) Update(ctx context.Context, liveArea *domain.LiveArea) error {
	query := `
		UPDATE live_areas
		SET name = $2, description = $3, is_public = $4, thumbnail_url = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		liveArea.ID, liveArea.Name, liveArea.Description,
		liveArea.IsPublic, liveArea.ThumbnailURL, liveArea.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update live area: %w", err)
	}
	return nil
}

func (r *liveAreaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM live_areas WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete live area: %w", err)
	}
	return nil
}

