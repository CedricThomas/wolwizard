package base

import (
	"log"
	"sync"
	"sync/atomic"

	ws "github.com/fasthttp/websocket"

	"github.com/CedricThomas/console/internal/service/websocket"
)

const (
	defaultBroadcastBufferSize = 100
)

var (
	ErrBroadcastBufferFull = &ErrBroadcastFull{
		msg: "broadcast buffer full",
	}
)

type ErrBroadcastFull struct {
	msg string
}

func (e *ErrBroadcastFull) Error() string {
	return e.msg
}

type manager struct {
	clients    map[string]*websocket.Client
	broadcast  chan []byte
	register   chan *websocket.Client
	unregister chan *websocket.Client
	mu         sync.RWMutex
	running    atomic.Bool
}

func New() *manager {
	m := &manager{
		clients:    make(map[string]*websocket.Client),
		broadcast:  make(chan []byte, defaultBroadcastBufferSize),
		register:   make(chan *websocket.Client),
		unregister: make(chan *websocket.Client),
	}
	return m
}

func (m *manager) Start() {
	log.Printf("WebSocket manager started")
	m.running.Store(true)
	go m.loop()
}

func (m *manager) loop() {
	for m.running.Load() {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			clientCount := len(m.clients)
			m.mu.Unlock()
			log.Printf("WebSocket: client registered - %s, total clients: %d", client.ID, clientCount)
			go m.writer(client)

		case client := <-m.unregister:
			m.mu.Lock()
			delete(m.clients, client.ID)
			m.mu.Unlock()
			log.Printf("WebSocket: client unregistered - %s", client.ID)
			close(client.Send)

		case msg := <-m.broadcast:
			m.mu.RLock()
			clientCount := len(m.clients)
			m.mu.RUnlock()
			log.Printf("WebSocket: broadcasting to %d clients, message size: %d bytes", clientCount, len(msg))
			m.mu.RLock()
			for _, client := range m.clients {
				select {
				case client.Send <- msg:
				default:
					log.Printf("WebSocket: failed to send to client (buffer full)")
				}
			}
			m.mu.RUnlock()
		}
	}
	log.Printf("WebSocket manager loop stopped")
}

func (m *manager) writer(client *websocket.Client) {
	defer func() {
		log.Printf("WebSocket: writer goroutine stopped for client %s", client.ID)
	}()

	for msg := range client.Send {
		conn := client.Conn.(*ws.Conn)
		err := conn.WriteMessage(ws.TextMessage, msg)
		if err != nil {
			log.Printf("WebSocket: error writing message to client %s: %v", client.ID, err)
			return
		}
	}
}

func (m *manager) Broadcast(msg []byte) error {
	select {
	case m.broadcast <- msg:
		return nil
	default:
		return ErrBroadcastBufferFull
	}
}

func (m *manager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

func (m *manager) Shutdown() {
	if !m.running.Swap(false) {
		return
	}

	log.Printf("WebSocket manager shutting down...")
	close(m.broadcast)

	m.mu.Lock()
	clientCount := len(m.clients)
	for _, client := range m.clients {
		close(client.Send)
	}
	m.clients = make(map[string]*websocket.Client)
	m.mu.Unlock()

	log.Printf("WebSocket manager shutdown complete, %d clients disconnected", clientCount)
}

func (m *manager) Register(client *websocket.Client) {
	m.register <- client
}

func (m *manager) Unregister(client *websocket.Client) {
	m.unregister <- client
}
