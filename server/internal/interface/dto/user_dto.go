package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// UserResponse はユーザー情報のレスポンスです
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	PhotoURL    string    `json:"photoUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UserUpdateRequest はユーザー更新リクエストです
type UserUpdateRequest struct {
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoUrl"`
}

// ToUserResponse はドメインモデルをレスポンスDTOに変換します
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

