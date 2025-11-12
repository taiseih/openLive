package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// UserUseCase はユーザー関連のユースケースを提供します
type UserUseCase interface {
	CreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, displayName, photoURL string) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase は新しいUserUseCaseを作成します
func NewUserUseCase(userRepo domain.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error) {
	// 既存ユーザーチェック
	existingUser, err := u.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	user := domain.NewUser(firebaseUID, email, displayName, photoURL)
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *userUseCase) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	user, err := u.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id uuid.UUID, displayName, photoURL string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	user.DisplayName = displayName
	user.PhotoURL = photoURL

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

