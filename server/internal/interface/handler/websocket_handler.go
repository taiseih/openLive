package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	wsHub "github.com/taiseihayashizaki/openLive/server/internal/infrastructure/websocket"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 本番環境では適切なオリジンチェックを実装してください
		return true
	},
}

// WebSocketHandler はWebSocket接続を処理します
type WebSocketHandler struct {
	hub *wsHub.Hub
}

// NewWebSocketHandler は新しいWebSocketHandlerを作成します
func NewWebSocketHandler(hub *wsHub.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket はWebSocket接続を処理します
// WS /api/v1/ws/streams/:id
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	userID, ok := c.Get(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid stream id")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return err
	}

	client := wsHub.NewClient(h.hub, conn, streamID, userID)
	h.hub.Register <- client

	// goroutineで読み書きを開始
	go client.WritePump()
	go client.ReadPump()

	return nil
}

