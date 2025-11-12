package websocket

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// クライアントへの書き込みタイムアウト
	writeWait = 10 * time.Second

	// クライアントからのPongメッセージ待機時間
	pongWait = 60 * time.Second

	// Pingメッセージ送信間隔
	pingPeriod = (pongWait * 9) / 10

	// 最大メッセージサイズ
	maxMessageSize = 512 * 1024 // 512KB
)

// Client はWebSocket接続の単一のクライアントを表します
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	streamID uuid.UUID
	userID   uuid.UUID
}

// Message はWebSocketメッセージを表します
type Message struct {
	StreamID uuid.UUID
	From     uuid.UUID
	To       uuid.UUID // 空の場合は全員にブロードキャスト
	Data     []byte
}

// NewClient は新しいWebSocketクライアントを作成します
func NewClient(hub *Hub, conn *websocket.Conn, streamID, userID uuid.UUID) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		streamID: streamID,
		userID:   userID,
	}
}

// ReadPump はWebSocket接続からメッセージを読み取ります
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// メッセージをハブにブロードキャスト
		c.hub.Broadcast <- &Message{
			StreamID: c.streamID,
			From:     c.userID,
			Data:     message,
		}
	}
}

// WritePump はWebSocket接続にメッセージを書き込みます
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hubがチャネルを閉じた
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// キューに溜まっているメッセージを追加
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

