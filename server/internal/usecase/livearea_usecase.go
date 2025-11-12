package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// LiveAreaUseCase はライブエリア関連のユースケースを提供します
type LiveAreaUseCase interface {
	CreateLiveArea(ctx context.Context, ownerID uuid.UUID, name, description string, isPublic bool) (*domain.LiveArea, error)
	GetLiveArea(ctx context.Context, id uuid.UUID) (*domain.LiveArea, error)
	GetPublicLiveAreas(ctx context.Context, limit, offset int) ([]*domain.LiveArea, error)
	GetUserLiveAreas(ctx context.Context, ownerID uuid.UUID) ([]*domain.LiveArea, error)
	UpdateLiveArea(ctx context.Context, id, ownerID uuid.UUID, name, description string, isPublic bool) (*domain.LiveArea, error)
	DeleteLiveArea(ctx context.Context, id, ownerID uuid.UUID) error
	InviteMember(ctx context.Context, liveAreaID, ownerID, memberID uuid.UUID, role domain.LiveAreaRole) error
	RemoveMember(ctx context.Context, liveAreaID, ownerID, memberID uuid.UUID) error
	GetMembers(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveAreaMember, error)
}

type liveAreaUseCase struct {
	liveAreaRepo       domain.LiveAreaRepository
	liveAreaMemberRepo domain.LiveAreaMemberRepository
}

// NewLiveAreaUseCase は新しいLiveAreaUseCaseを作成します
func NewLiveAreaUseCase(
	liveAreaRepo domain.LiveAreaRepository,
	liveAreaMemberRepo domain.LiveAreaMemberRepository,
) LiveAreaUseCase {
	return &liveAreaUseCase{
		liveAreaRepo:       liveAreaRepo,
		liveAreaMemberRepo: liveAreaMemberRepo,
	}
}

func (u *liveAreaUseCase) CreateLiveArea(ctx context.Context, ownerID uuid.UUID, name, description string, isPublic bool) (*domain.LiveArea, error) {
	liveArea := domain.NewLiveArea(ownerID, name, description, isPublic)
	if err := u.liveAreaRepo.Create(ctx, liveArea); err != nil {
		return nil, fmt.Errorf("failed to create live area: %w", err)
	}

	// オーナーをメンバーとして追加
	member := domain.NewLiveAreaMember(liveArea.ID, ownerID, domain.RoleOwner)
	if err := u.liveAreaMemberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	return liveArea, nil
}

func (u *liveAreaUseCase) GetLiveArea(ctx context.Context, id uuid.UUID) (*domain.LiveArea, error) {
	liveArea, err := u.liveAreaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get live area: %w", err)
	}
	return liveArea, nil
}

func (u *liveAreaUseCase) GetPublicLiveAreas(ctx context.Context, limit, offset int) ([]*domain.LiveArea, error) {
	areas, err := u.liveAreaRepo.FindPublicAreas(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public live areas: %w", err)
	}
	return areas, nil
}

func (u *liveAreaUseCase) GetUserLiveAreas(ctx context.Context, ownerID uuid.UUID) ([]*domain.LiveArea, error) {
	areas, err := u.liveAreaRepo.FindByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user live areas: %w", err)
	}
	return areas, nil
}

func (u *liveAreaUseCase) UpdateLiveArea(ctx context.Context, id, ownerID uuid.UUID, name, description string, isPublic bool) (*domain.LiveArea, error) {
	liveArea, err := u.liveAreaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find live area: %w", err)
	}

	// オーナーチェック
	if liveArea.OwnerID != ownerID {
		return nil, fmt.Errorf("unauthorized: not the owner")
	}

	liveArea.Name = name
	liveArea.Description = description
	liveArea.IsPublic = isPublic

	if err := u.liveAreaRepo.Update(ctx, liveArea); err != nil {
		return nil, fmt.Errorf("failed to update live area: %w", err)
	}

	return liveArea, nil
}

func (u *liveAreaUseCase) DeleteLiveArea(ctx context.Context, id, ownerID uuid.UUID) error {
	liveArea, err := u.liveAreaRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find live area: %w", err)
	}

	// オーナーチェック
	if liveArea.OwnerID != ownerID {
		return fmt.Errorf("unauthorized: not the owner")
	}

	if err := u.liveAreaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete live area: %w", err)
	}

	return nil
}

func (u *liveAreaUseCase) InviteMember(ctx context.Context, liveAreaID, ownerID, memberID uuid.UUID, role domain.LiveAreaRole) error {
	liveArea, err := u.liveAreaRepo.FindByID(ctx, liveAreaID)
	if err != nil {
		return fmt.Errorf("failed to find live area: %w", err)
	}

	// オーナーチェック
	if liveArea.OwnerID != ownerID {
		return fmt.Errorf("unauthorized: not the owner")
	}

	// 既存メンバーチェック
	existingMember, _ := u.liveAreaMemberRepo.FindByLiveAreaAndUser(ctx, liveAreaID, memberID)
	if existingMember != nil {
		return fmt.Errorf("user is already a member")
	}

	member := domain.NewLiveAreaMember(liveAreaID, memberID, role)
	if err := u.liveAreaMemberRepo.Create(ctx, member); err != nil {
		return fmt.Errorf("failed to invite member: %w", err)
	}

	return nil
}

func (u *liveAreaUseCase) RemoveMember(ctx context.Context, liveAreaID, ownerID, memberID uuid.UUID) error {
	liveArea, err := u.liveAreaRepo.FindByID(ctx, liveAreaID)
	if err != nil {
		return fmt.Errorf("failed to find live area: %w", err)
	}

	// オーナーチェック
	if liveArea.OwnerID != ownerID {
		return fmt.Errorf("unauthorized: not the owner")
	}

	// オーナー自身は削除できない
	if memberID == ownerID {
		return fmt.Errorf("cannot remove the owner")
	}

	if err := u.liveAreaMemberRepo.Delete(ctx, liveAreaID, memberID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

func (u *liveAreaUseCase) GetMembers(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveAreaMember, error) {
	members, err := u.liveAreaMemberRepo.FindByLiveAreaID(ctx, liveAreaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	return members, nil
}

