package base

import (
	"sync"
	"sync/atomic"

	"github.com/CedricThomas/console/internal/service/websocket"
	ws "github.com/gofiber/contrib/websocket"
)

const (
	defaultBroadcastBufferSize  = 100
	defaultClientSendBufferSize = 256
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
	m.running.Store(true)
	go m.loop()
}

func (m *manager) loop() {
	for m.running.Load() {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			m.mu.Unlock()
			go m.writer(client)

		case client := <-m.unregister:
			m.mu.Lock()
			delete(m.clients, client.ID)
			m.mu.Unlock()
			close(client.Send)

		case msg := <-m.broadcast:
			m.mu.RLock()
			for _, client := range m.clients {
				select {
				case client.Send <- msg:
				default:
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *manager) writer(client *websocket.Client) {
	for msg := range client.Send {
		conn := client.Conn.(*ws.Conn)
		err := conn.WriteMessage(ws.TextMessage, msg)
		if err != nil {
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

	close(m.broadcast)

	m.mu.Lock()
	for _, client := range m.clients {
		close(client.Send)
	}
	m.clients = make(map[string]*websocket.Client)
	m.mu.Unlock()
}

func (m *manager) Register(client *websocket.Client) {
	m.register <- client
}

func (m *manager) Unregister(client *websocket.Client) {
	m.unregister <- client
}
