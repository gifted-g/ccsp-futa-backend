package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"ccsp-futa-alumni/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

type HubType struct {
	clients map[string]map[*Client]bool // userID -> clients
	mu      sync.RWMutex
}

var Hub = &HubType{
	clients: make(map[string]map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade", err)
		return
	}
	client := &Client{UserID: userID, Conn: conn}
	Hub.addClient(client)
	go client.readPump()
}

func (h *HubType) addClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c.UserID]; !ok {
		h.clients[c.UserID] = make(map[*Client]bool)
	}
	h.clients[c.UserID][c] = true
}

func (h *HubType) removeClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[c.UserID]; ok {
		if _, exists := conns[c]; exists {
			delete(conns, c)
			c.Conn.Close()
		}
		if len(conns) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (c *Client) readPump() {
	defer Hub.removeClient(c)
	for {
		// Keep connection alive; we don't expect messages from client in this simple hub
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			// disconnected
			return
		}
	}
}

// BroadcastMessage sends message to all members that are connected.
// For simplicity we send to every connected user — in real app you'd lookup channel members and send to them only.
func (h *HubType) BroadcastMessage(msg models.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	payload, _ := json.Marshal(msg)
	for userID, conns := range h.clients {
		_ = userID // In production, filter by actual channel membership
		for c := range conns {
			if err := c.Conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Println("ws write err:", err)
				go h.removeClient(c)
			}
		}
	}
}
