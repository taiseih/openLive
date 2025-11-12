package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// LiveAreaResponse はライブエリアのレスポンスです
type LiveAreaResponse struct {
	ID           uuid.UUID `json:"id"`
	OwnerID      uuid.UUID `json:"ownerId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsPublic     bool      `json:"isPublic"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// LiveAreaCreateRequest はライブエリア作成リクエストです
type LiveAreaCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

// LiveAreaUpdateRequest はライブエリア更新リクエストです
type LiveAreaUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

// InviteMemberRequest はメンバー招待リクエストです
type InviteMemberRequest struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
	Role   string    `json:"role" validate:"required,oneof=member moderator"`
}

// LiveAreaMemberResponse はライブエリアメンバーのレスポンスです
type LiveAreaMemberResponse struct {
	ID         uuid.UUID `json:"id"`
	LiveAreaID uuid.UUID `json:"liveAreaId"`
	UserID     uuid.UUID `json:"userId"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joinedAt"`
}

// ToLiveAreaResponse はドメインモデルをレスポンスDTOに変換します
func ToLiveAreaResponse(area *domain.LiveArea) *LiveAreaResponse {
	return &LiveAreaResponse{
		ID:           area.ID,
		OwnerID:      area.OwnerID,
		Name:         area.Name,
		Description:  area.Description,
		IsPublic:     area.IsPublic,
		ThumbnailURL: area.ThumbnailURL,
		CreatedAt:    area.CreatedAt,
		UpdatedAt:    area.UpdatedAt,
	}
}

// ToLiveAreaMemberResponse はドメインモデルをレスポンスDTOに変換します
func ToLiveAreaMemberResponse(member *domain.LiveAreaMember) *LiveAreaMemberResponse {
	return &LiveAreaMemberResponse{
		ID:         member.ID,
		LiveAreaID: member.LiveAreaID,
		UserID:     member.UserID,
		Role:       string(member.Role),
		JoinedAt:   member.JoinedAt,
	}
}

