package domain

import (
	"time"

	"github.com/google/uuid"
)

// User はユーザーエンティティを表します
type User struct {
	ID          uuid.UUID
	FirebaseUID string
	Email       string
	DisplayName string
	PhotoURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewUser は新しいユーザーを作成します
func NewUser(firebaseUID, email, displayName, photoURL string) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New(),
		FirebaseUID: firebaseUID,
		Email:       email,
		DisplayName: displayName,
		PhotoURL:    photoURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// LiveArea はライブエリア（配信広場）エンティティを表します
type LiveArea struct {
	ID           uuid.UUID
	OwnerID      uuid.UUID
	Name         string
	Description  string
	IsPublic     bool
	ThumbnailURL string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewLiveArea は新しいライブエリアを作成します
func NewLiveArea(ownerID uuid.UUID, name, description string, isPublic bool) *LiveArea {
	now := time.Now()
	return &LiveArea{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// LiveAreaRole はライブエリアにおける役割を表します
type LiveAreaRole string

const (
	RoleOwner     LiveAreaRole = "owner"
	RoleModerator LiveAreaRole = "moderator"
	RoleMember    LiveAreaRole = "member"
)

// LiveAreaMember はライブエリアメンバーエンティティを表します
type LiveAreaMember struct {
	ID         uuid.UUID
	LiveAreaID uuid.UUID
	UserID     uuid.UUID
	Role       LiveAreaRole
	JoinedAt   time.Time
}

// NewLiveAreaMember は新しいライブエリアメンバーを作成します
func NewLiveAreaMember(liveAreaID, userID uuid.UUID, role LiveAreaRole) *LiveAreaMember {
	return &LiveAreaMember{
		ID:         uuid.New(),
		LiveAreaID: liveAreaID,
		UserID:     userID,
		Role:       role,
		JoinedAt:   time.Now(),
	}
}

// StreamStatus は配信ステータスを表します
type StreamStatus string

const (
	StatusScheduled StreamStatus = "scheduled"
	StatusLive      StreamStatus = "live"
	StatusEnded     StreamStatus = "ended"
)

// LiveStream はライブ配信エンティティを表します
type LiveStream struct {
	ID           uuid.UUID
	LiveAreaID   uuid.UUID
	StreamerID   uuid.UUID
	Title        string
	Description  string
	Status       StreamStatus
	ThumbnailURL string
	StartedAt    *time.Time
	EndedAt      *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewLiveStream は新しいライブ配信を作成します
func NewLiveStream(liveAreaID, streamerID uuid.UUID, title, description string) *LiveStream {
	now := time.Now()
	return &LiveStream{
		ID:          uuid.New(),
		LiveAreaID:  liveAreaID,
		StreamerID:  streamerID,
		Title:       title,
		Description: description,
		Status:      StatusScheduled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Start は配信を開始します
func (ls *LiveStream) Start() {
	now := time.Now()
	ls.Status = StatusLive
	ls.StartedAt = &now
}

// End は配信を終了します
func (ls *LiveStream) End() {
	now := time.Now()
	ls.Status = StatusEnded
	ls.EndedAt = &now
}

// IsLive は配信中かどうかを返します
func (ls *LiveStream) IsLive() bool {
	return ls.Status == StatusLive
}

// LiveStreamViewer は配信視聴者エンティティを表します
type LiveStreamViewer struct {
	ID           uuid.UUID
	LiveStreamID uuid.UUID
	UserID       *uuid.UUID
	JoinedAt     time.Time
	LeftAt       *time.Time
	IsActive     bool
}

// NewLiveStreamViewer は新しい視聴者を作成します
func NewLiveStreamViewer(liveStreamID uuid.UUID, userID *uuid.UUID) *LiveStreamViewer {
	return &LiveStreamViewer{
		ID:           uuid.New(),
		LiveStreamID: liveStreamID,
		UserID:       userID,
		JoinedAt:     time.Now(),
		IsActive:     true,
	}
}

// Leave は視聴者が退出したことを記録します
func (v *LiveStreamViewer) Leave() {
	now := time.Now()
	v.LeftAt = &now
	v.IsActive = false
}

