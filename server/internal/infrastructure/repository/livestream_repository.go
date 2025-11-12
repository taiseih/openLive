package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

type liveStreamRepository struct {
	db *sql.DB
}

// NewLiveStreamRepository は新しいライブ配信リポジトリを作成します
func NewLiveStreamRepository(db *sql.DB) domain.LiveStreamRepository {
	return &liveStreamRepository{db: db}
}

func (r *liveStreamRepository) Create(ctx context.Context, stream *domain.LiveStream) error {
	query := `
		INSERT INTO live_streams (id, live_area_id, streamer_id, title, description, status, thumbnail_url, started_at, ended_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		stream.ID, stream.LiveAreaID, stream.StreamerID, stream.Title, stream.Description,
		stream.Status, stream.ThumbnailURL, stream.StartedAt, stream.EndedAt,
		stream.CreatedAt, stream.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create live stream: %w", err)
	}
	return nil
}

func (r *liveStreamRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LiveStream, error) {
	query := `
		SELECT id, live_area_id, streamer_id, title, description, status, thumbnail_url, started_at, ended_at, created_at, updated_at
		FROM live_streams WHERE id = $1
	`
	stream := &domain.LiveStream{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stream.ID, &stream.LiveAreaID, &stream.StreamerID, &stream.Title, &stream.Description,
		&stream.Status, &stream.ThumbnailURL, &stream.StartedAt, &stream.EndedAt,
		&stream.CreatedAt, &stream.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("live stream not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find live stream: %w", err)
	}
	return stream, nil
}

func (r *liveStreamRepository) FindByLiveAreaID(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveStream, error) {
	query := `
		SELECT id, live_area_id, streamer_id, title, description, status, thumbnail_url, started_at, ended_at, created_at, updated_at
		FROM live_streams WHERE live_area_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, liveAreaID)
	if err != nil {
		return nil, fmt.Errorf("failed to find live streams: %w", err)
	}
	defer rows.Close()

	var streams []*domain.LiveStream
	for rows.Next() {
		stream := &domain.LiveStream{}
		if err := rows.Scan(
			&stream.ID, &stream.LiveAreaID, &stream.StreamerID, &stream.Title, &stream.Description,
			&stream.Status, &stream.ThumbnailURL, &stream.StartedAt, &stream.EndedAt,
			&stream.CreatedAt, &stream.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan live stream: %w", err)
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func (r *liveStreamRepository) FindByStreamerID(ctx context.Context, streamerID uuid.UUID) ([]*domain.LiveStream, error) {
	query := `
		SELECT id, live_area_id, streamer_id, title, description, status, thumbnail_url, started_at, ended_at, created_at, updated_at
		FROM live_streams WHERE streamer_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, streamerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find live streams: %w", err)
	}
	defer rows.Close()

	var streams []*domain.LiveStream
	for rows.Next() {
		stream := &domain.LiveStream{}
		if err := rows.Scan(
			&stream.ID, &stream.LiveAreaID, &stream.StreamerID, &stream.Title, &stream.Description,
			&stream.Status, &stream.ThumbnailURL, &stream.StartedAt, &stream.EndedAt,
			&stream.CreatedAt, &stream.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan live stream: %w", err)
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func (r *liveStreamRepository) FindLiveStreams(ctx context.Context, limit, offset int) ([]*domain.LiveStream, error) {
	query := `
		SELECT id, live_area_id, streamer_id, title, description, status, thumbnail_url, started_at, ended_at, created_at, updated_at
		FROM live_streams WHERE status = 'live' ORDER BY started_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find live streams: %w", err)
	}
	defer rows.Close()

	var streams []*domain.LiveStream
	for rows.Next() {
		stream := &domain.LiveStream{}
		if err := rows.Scan(
			&stream.ID, &stream.LiveAreaID, &stream.StreamerID, &stream.Title, &stream.Description,
			&stream.Status, &stream.ThumbnailURL, &stream.StartedAt, &stream.EndedAt,
			&stream.CreatedAt, &stream.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan live stream: %w", err)
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func (r *liveStreamRepository) Update(ctx context.Context, stream *domain.LiveStream) error {
	query := `
		UPDATE live_streams
		SET title = $2, description = $3, status = $4, thumbnail_url = $5, started_at = $6, ended_at = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		stream.ID, stream.Title, stream.Description, stream.Status,
		stream.ThumbnailURL, stream.StartedAt, stream.EndedAt, stream.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update live stream: %w", err)
	}
	return nil
}

func (r *liveStreamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM live_streams WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete live stream: %w", err)
	}
	return nil
}

