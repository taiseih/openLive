package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/dto"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/middleware"
	"github.com/taiseihayashizaki/openLive/server/internal/usecase"
)

// LiveStreamHandler はライブ配信関連のハンドラーです
type LiveStreamHandler struct {
	liveStreamUseCase usecase.LiveStreamUseCase
}

// NewLiveStreamHandler は新しいLiveStreamHandlerを作成します
func NewLiveStreamHandler(liveStreamUseCase usecase.LiveStreamUseCase) *LiveStreamHandler {
	return &LiveStreamHandler{
		liveStreamUseCase: liveStreamUseCase,
	}
}

// Create はライブ配信を作成します
// POST /api/v1/liveareas/:id/streams
func (h *LiveStreamHandler) Create(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	liveAreaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid live area id")
	}

	var req dto.LiveStreamCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	stream, err := h.liveStreamUseCase.CreateStream(c.Request().Context(), liveAreaID, userID, req.Title, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusCreated, dto.ToLiveStreamResponse(stream))
}

// GetByLiveArea はライブエリアの配信一覧を取得します
// GET /api/v1/liveareas/:id/streams
func (h *LiveStreamHandler) GetByLiveArea(c echo.Context) error {
	liveAreaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid live area id")
	}

	streams, err := h.liveStreamUseCase.GetStreamsByLiveArea(c.Request().Context(), liveAreaID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get streams")
	}

	responses := make([]*dto.LiveStreamResponse, len(streams))
	for i, stream := range streams {
		responses[i] = dto.ToLiveStreamResponse(stream)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetLiveStreams は配信中の配信一覧を取得します
// GET /api/v1/streams/live
func (h *LiveStreamHandler) GetLiveStreams(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	streams, err := h.liveStreamUseCase.GetLiveStreams(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get live streams")
	}

	responses := make([]*dto.LiveStreamResponse, len(streams))
	for i, stream := range streams {
		responses[i] = dto.ToLiveStreamResponse(stream)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetByID は配信詳細を取得します
// GET /api/v1/streams/:id
func (h *LiveStreamHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	stream, err := h.liveStreamUseCase.GetStream(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "stream not found")
	}

	return c.JSON(http.StatusOK, dto.ToLiveStreamResponse(stream))
}

// Start は配信を開始します
// POST /api/v1/streams/:id/start
func (h *LiveStreamHandler) Start(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	stream, err := h.liveStreamUseCase.StartStream(c.Request().Context(), id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToLiveStreamResponse(stream))
}

// End は配信を終了します
// DELETE /api/v1/streams/:id
func (h *LiveStreamHandler) End(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	stream, err := h.liveStreamUseCase.EndStream(c.Request().Context(), id, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToLiveStreamResponse(stream))
}

// GetViewers は視聴者一覧を取得します
// GET /api/v1/streams/:id/viewers
func (h *LiveStreamHandler) GetViewers(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	viewers, err := h.liveStreamUseCase.GetViewers(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get viewers")
	}

	responses := make([]*dto.LiveStreamViewerResponse, len(viewers))
	for i, viewer := range viewers {
		responses[i] = dto.ToLiveStreamViewerResponse(viewer)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetViewerCount は視聴者数を取得します
// GET /api/v1/streams/:id/viewers/count
func (h *LiveStreamHandler) GetViewerCount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	count, err := h.liveStreamUseCase.GetViewerCount(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get viewer count")
	}

	return c.JSON(http.StatusOK, dto.ViewerCountResponse{Count: count})
}

