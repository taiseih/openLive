package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository はユーザーのリポジトリインターフェースです
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LiveAreaRepository はライブエリアのリポジトリインターフェースです
type LiveAreaRepository interface {
	Create(ctx context.Context, liveArea *LiveArea) error
	FindByID(ctx context.Context, id uuid.UUID) (*LiveArea, error)
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*LiveArea, error)
	FindPublicAreas(ctx context.Context, limit, offset int) ([]*LiveArea, error)
	Update(ctx context.Context, liveArea *LiveArea) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LiveAreaMemberRepository はライブエリアメンバーのリポジトリインターフェースです
type LiveAreaMemberRepository interface {
	Create(ctx context.Context, member *LiveAreaMember) error
	FindByLiveAreaID(ctx context.Context, liveAreaID uuid.UUID) ([]*LiveAreaMember, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*LiveAreaMember, error)
	FindByLiveAreaAndUser(ctx context.Context, liveAreaID, userID uuid.UUID) (*LiveAreaMember, error)
	Delete(ctx context.Context, liveAreaID, userID uuid.UUID) error
}

// LiveStreamRepository はライブ配信のリポジトリインターフェースです
type LiveStreamRepository interface {
	Create(ctx context.Context, stream *LiveStream) error
	FindByID(ctx context.Context, id uuid.UUID) (*LiveStream, error)
	FindByLiveAreaID(ctx context.Context, liveAreaID uuid.UUID) ([]*LiveStream, error)
	FindByStreamerID(ctx context.Context, streamerID uuid.UUID) ([]*LiveStream, error)
	FindLiveStreams(ctx context.Context, limit, offset int) ([]*LiveStream, error)
	Update(ctx context.Context, stream *LiveStream) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LiveStreamViewerRepository は配信視聴者のリポジトリインターフェースです
type LiveStreamViewerRepository interface {
	Create(ctx context.Context, viewer *LiveStreamViewer) error
	FindByStreamID(ctx context.Context, streamID uuid.UUID) ([]*LiveStreamViewer, error)
	FindActiveByStreamID(ctx context.Context, streamID uuid.UUID) ([]*LiveStreamViewer, error)
	CountActiveByStreamID(ctx context.Context, streamID uuid.UUID) (int, error)
	Update(ctx context.Context, viewer *LiveStreamViewer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

