package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

type liveAreaMemberRepository struct {
	db *sql.DB
}

// NewLiveAreaMemberRepository は新しいライブエリアメンバーリポジトリを作成します
func NewLiveAreaMemberRepository(db *sql.DB) domain.LiveAreaMemberRepository {
	return &liveAreaMemberRepository{db: db}
}

func (r *liveAreaMemberRepository) Create(ctx context.Context, member *domain.LiveAreaMember) error {
	query := `
		INSERT INTO live_area_members (id, live_area_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.LiveAreaID, member.UserID, member.Role, member.JoinedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create live area member: %w", err)
	}
	return nil
}

func (r *liveAreaMemberRepository) FindByLiveAreaID(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveAreaMember, error) {
	query := `
		SELECT id, live_area_id, user_id, role, joined_at
		FROM live_area_members WHERE live_area_id = $1 ORDER BY joined_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, liveAreaID)
	if err != nil {
		return nil, fmt.Errorf("failed to find members: %w", err)
	}
	defer rows.Close()

	var members []*domain.LiveAreaMember
	for rows.Next() {
		member := &domain.LiveAreaMember{}
		if err := rows.Scan(
			&member.ID, &member.LiveAreaID, &member.UserID, &member.Role, &member.JoinedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *liveAreaMemberRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.LiveAreaMember, error) {
	query := `
		SELECT id, live_area_id, user_id, role, joined_at
		FROM live_area_members WHERE user_id = $1 ORDER BY joined_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find members: %w", err)
	}
	defer rows.Close()

	var members []*domain.LiveAreaMember
	for rows.Next() {
		member := &domain.LiveAreaMember{}
		if err := rows.Scan(
			&member.ID, &member.LiveAreaID, &member.UserID, &member.Role, &member.JoinedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *liveAreaMemberRepository) FindByLiveAreaAndUser(ctx context.Context, liveAreaID, userID uuid.UUID) (*domain.LiveAreaMember, error) {
	query := `
		SELECT id, live_area_id, user_id, role, joined_at
		FROM live_area_members WHERE live_area_id = $1 AND user_id = $2
	`
	member := &domain.LiveAreaMember{}
	err := r.db.QueryRowContext(ctx, query, liveAreaID, userID).Scan(
		&member.ID, &member.LiveAreaID, &member.UserID, &member.Role, &member.JoinedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find member: %w", err)
	}
	return member, nil
}

func (r *liveAreaMemberRepository) Delete(ctx context.Context, liveAreaID, userID uuid.UUID) error {
	query := `DELETE FROM live_area_members WHERE live_area_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, liveAreaID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}

