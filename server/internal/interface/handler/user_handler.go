package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/dto"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/middleware"
	"github.com/taiseihayashizaki/openLive/server/internal/usecase"
)

// UserHandler はユーザー関連のハンドラーです
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler は新しいUserHandlerを作成します
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetMe は自分のユーザー情報を取得します
// GET /api/v1/users/me
func (h *UserHandler) GetMe(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	user, err := h.userUseCase.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// UpdateMe は自分のユーザー情報を更新します
// PUT /api/v1/users/me
func (h *UserHandler) UpdateMe(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req dto.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.userUseCase.UpdateUser(c.Request().Context(), userID, req.DisplayName, req.PhotoURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
	}

	return c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

