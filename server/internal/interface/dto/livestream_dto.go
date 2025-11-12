package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// LiveStreamResponse はライブ配信のレスポンスです
type LiveStreamResponse struct {
	ID           uuid.UUID  `json:"id"`
	LiveAreaID   uuid.UUID  `json:"liveAreaId"`
	StreamerID   uuid.UUID  `json:"streamerId"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	ThumbnailURL string     `json:"thumbnailUrl"`
	StartedAt    *time.Time `json:"startedAt"`
	EndedAt      *time.Time `json:"endedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// LiveStreamCreateRequest はライブ配信作成リクエストです
type LiveStreamCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

// LiveStreamViewerResponse は視聴者のレスポンスです
type LiveStreamViewerResponse struct {
	ID           uuid.UUID  `json:"id"`
	LiveStreamID uuid.UUID  `json:"liveStreamId"`
	UserID       *uuid.UUID `json:"userId"`
	JoinedAt     time.Time  `json:"joinedAt"`
	IsActive     bool       `json:"isActive"`
}

// ViewerCountResponse は視聴者数のレスポンスです
type ViewerCountResponse struct {
	Count int `json:"count"`
}

// ToLiveStreamResponse はドメインモデルをレスポンスDTOに変換します
func ToLiveStreamResponse(stream *domain.LiveStream) *LiveStreamResponse {
	return &LiveStreamResponse{
		ID:           stream.ID,
		LiveAreaID:   stream.LiveAreaID,
		StreamerID:   stream.StreamerID,
		Title:        stream.Title,
		Description:  stream.Description,
		Status:       string(stream.Status),
		ThumbnailURL: stream.ThumbnailURL,
		StartedAt:    stream.StartedAt,
		EndedAt:      stream.EndedAt,
		CreatedAt:    stream.CreatedAt,
		UpdatedAt:    stream.UpdatedAt,
	}
}

// ToLiveStreamViewerResponse はドメインモデルをレスポンスDTOに変換します
func ToLiveStreamViewerResponse(viewer *domain.LiveStreamViewer) *LiveStreamViewerResponse {
	return &LiveStreamViewerResponse{
		ID:           viewer.ID,
		LiveStreamID: viewer.LiveStreamID,
		UserID:       viewer.UserID,
		JoinedAt:     viewer.JoinedAt,
		IsActive:     viewer.IsActive,
	}
}

