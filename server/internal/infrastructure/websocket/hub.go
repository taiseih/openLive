package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// Hub はWebSocket接続を管理します
type Hub struct {
	// 配信IDごとのクライアント管理
	clients map[uuid.UUID]map[*Client]bool

	// クライアントの登録リクエスト
	Register chan *Client

	// クライアントの登録解除リクエスト
	Unregister chan *Client

	// クライアントへのメッセージブロードキャスト
	Broadcast chan *Message

	mu sync.RWMutex
}

// NewHub は新しいHubを作成します
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 256),
	}
}

// Run はHubを起動します
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.streamID] == nil {
		h.clients[client.streamID] = make(map[*Client]bool)
	}
	h.clients[client.streamID][client] = true
	log.Printf("Client registered: streamID=%s, userID=%s, total=%d",
		client.streamID, client.userID, len(h.clients[client.streamID]))
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.streamID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)
			log.Printf("Client unregistered: streamID=%s, userID=%s, remaining=%d",
				client.streamID, client.userID, len(clients))

			// 配信に接続しているクライアントがいなくなったら削除
			if len(clients) == 0 {
				delete(h.clients, client.streamID)
			}
		}
	}
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[message.StreamID]
	if !ok {
		return
	}

	// 特定のユーザーへのメッセージの場合
	if message.To != uuid.Nil {
		for client := range clients {
			if client.userID == message.To {
				select {
				case client.send <- message.Data:
				default:
					// 送信できない場合はクライアントを閉じる
					close(client.send)
					delete(clients, client)
				}
				return
			}
		}
		return
	}

	// 全員へのブロードキャスト
	for client := range clients {
		select {
		case client.send <- message.Data:
		default:
			// 送信できない場合はクライアントを閉じる
			close(client.send)
			delete(clients, client)
		}
	}
}

// GetStreamViewerCount は配信の視聴者数を取得します
func (h *Hub) GetStreamViewerCount(streamID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[streamID]; ok {
		return len(clients)
	}
	return 0
}

// BroadcastToStream は特定の配信にメッセージをブロードキャストします
func (h *Hub) BroadcastToStream(streamID uuid.UUID, data []byte) {
	h.Broadcast <- &Message{
		StreamID: streamID,
		Data:     data,
	}
}

// SendToUser は特定のユーザーにメッセージを送信します
func (h *Hub) SendToUser(streamID, userID uuid.UUID, data []byte) {
	h.Broadcast <- &Message{
		StreamID: streamID,
		To:       userID,
		Data:     data,
	}
}

