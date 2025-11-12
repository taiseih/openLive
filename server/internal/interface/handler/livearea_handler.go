package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/dto"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/middleware"
	"github.com/taiseihayashizaki/openLive/server/internal/usecase"
)

// LiveAreaHandler はライブエリア関連のハンドラーです
type LiveAreaHandler struct {
	liveAreaUseCase usecase.LiveAreaUseCase
}

// NewLiveAreaHandler は新しいLiveAreaHandlerを作成します
func NewLiveAreaHandler(liveAreaUseCase usecase.LiveAreaUseCase) *LiveAreaHandler {
	return &LiveAreaHandler{
		liveAreaUseCase: liveAreaUseCase,
	}
}

// Create はライブエリアを作成します
// POST /api/v1/liveareas
func (h *LiveAreaHandler) Create(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req dto.LiveAreaCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	liveArea, err := h.liveAreaUseCase.CreateLiveArea(c.Request().Context(), userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create live area")
	}

	return c.JSON(http.StatusCreated, dto.ToLiveAreaResponse(liveArea))
}

// GetAll は公開ライブエリア一覧を取得します
// GET /api/v1/liveareas
func (h *LiveAreaHandler) GetAll(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	areas, err := h.liveAreaUseCase.GetPublicLiveAreas(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get live areas")
	}

	responses := make([]*dto.LiveAreaResponse, len(areas))
	for i, area := range areas {
		responses[i] = dto.ToLiveAreaResponse(area)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetByID はライブエリア詳細を取得します
// GET /api/v1/liveareas/:id
func (h *LiveAreaHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	liveArea, err := h.liveAreaUseCase.GetLiveArea(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "live area not found")
	}

	return c.JSON(http.StatusOK, dto.ToLiveAreaResponse(liveArea))
}

// Update はライブエリアを更新します
// PUT /api/v1/liveareas/:id
func (h *LiveAreaHandler) Update(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req dto.LiveAreaUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	liveArea, err := h.liveAreaUseCase.UpdateLiveArea(c.Request().Context(), id, userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToLiveAreaResponse(liveArea))
}

// Delete はライブエリアを削除します
// DELETE /api/v1/liveareas/:id
func (h *LiveAreaHandler) Delete(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.liveAreaUseCase.DeleteLiveArea(c.Request().Context(), id, userID); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// InviteMember はメンバーを招待します
// POST /api/v1/liveareas/:id/invite
func (h *LiveAreaHandler) InviteMember(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req dto.InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	role := domain.LiveAreaRole(req.Role)
	if err := h.liveAreaUseCase.InviteMember(c.Request().Context(), id, userID, req.UserID, role); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveMember はメンバーを削除します
// DELETE /api/v1/liveareas/:id/members/:userId
func (h *LiveAreaHandler) RemoveMember(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid live area id")
	}

	memberID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	if err := h.liveAreaUseCase.RemoveMember(c.Request().Context(), id, userID, memberID); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMembers はメンバー一覧を取得します
// GET /api/v1/liveareas/:id/members
func (h *LiveAreaHandler) GetMembers(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	members, err := h.liveAreaUseCase.GetMembers(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get members")
	}

	responses := make([]*dto.LiveAreaMemberResponse, len(members))
	for i, member := range members {
		responses[i] = dto.ToLiveAreaMemberResponse(member)
	}

	return c.JSON(http.StatusOK, responses)
}

