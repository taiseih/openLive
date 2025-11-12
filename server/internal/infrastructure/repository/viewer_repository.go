package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

type liveStreamViewerRepository struct {
	db *sql.DB
}

// NewLiveStreamViewerRepository は新しい視聴者リポジトリを作成します
func NewLiveStreamViewerRepository(db *sql.DB) domain.LiveStreamViewerRepository {
	return &liveStreamViewerRepository{db: db}
}

func (r *liveStreamViewerRepository) Create(ctx context.Context, viewer *domain.LiveStreamViewer) error {
	query := `
		INSERT INTO live_stream_viewers (id, live_stream_id, user_id, joined_at, left_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		viewer.ID, viewer.LiveStreamID, viewer.UserID, viewer.JoinedAt, viewer.LeftAt, viewer.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create viewer: %w", err)
	}
	return nil
}

func (r *liveStreamViewerRepository) FindByStreamID(ctx context.Context, streamID uuid.UUID) ([]*domain.LiveStreamViewer, error) {
	query := `
		SELECT id, live_stream_id, user_id, joined_at, left_at, is_active
		FROM live_stream_viewers WHERE live_stream_id = $1 ORDER BY joined_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to find viewers: %w", err)
	}
	defer rows.Close()

	var viewers []*domain.LiveStreamViewer
	for rows.Next() {
		viewer := &domain.LiveStreamViewer{}
		if err := rows.Scan(
			&viewer.ID, &viewer.LiveStreamID, &viewer.UserID, &viewer.JoinedAt, &viewer.LeftAt, &viewer.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to scan viewer: %w", err)
		}
		viewers = append(viewers, viewer)
	}
	return viewers, nil
}

func (r *liveStreamViewerRepository) FindActiveByStreamID(ctx context.Context, streamID uuid.UUID) ([]*domain.LiveStreamViewer, error) {
	query := `
		SELECT id, live_stream_id, user_id, joined_at, left_at, is_active
		FROM live_stream_viewers WHERE live_stream_id = $1 AND is_active = true ORDER BY joined_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to find active viewers: %w", err)
	}
	defer rows.Close()

	var viewers []*domain.LiveStreamViewer
	for rows.Next() {
		viewer := &domain.LiveStreamViewer{}
		if err := rows.Scan(
			&viewer.ID, &viewer.LiveStreamID, &viewer.UserID, &viewer.JoinedAt, &viewer.LeftAt, &viewer.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to scan viewer: %w", err)
		}
		viewers = append(viewers, viewer)
	}
	return viewers, nil
}

func (r *liveStreamViewerRepository) CountActiveByStreamID(ctx context.Context, streamID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM live_stream_viewers WHERE live_stream_id = $1 AND is_active = true`
	var count int
	err := r.db.QueryRowContext(ctx, query, streamID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count viewers: %w", err)
	}
	return count, nil
}

func (r *liveStreamViewerRepository) Update(ctx context.Context, viewer *domain.LiveStreamViewer) error {
	query := `
		UPDATE live_stream_viewers
		SET left_at = $2, is_active = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, viewer.ID, viewer.LeftAt, viewer.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update viewer: %w", err)
	}
	return nil
}

func (r *liveStreamViewerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM live_stream_viewers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete viewer: %w", err)
	}
	return nil
}

