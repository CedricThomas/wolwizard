package handlers

import (
	"context"
	"log"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"

	ws "github.com/CedricThomas/console/internal/service/websocket"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(manager ws.Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		agentID := c.Params("agent_id")
		username := "anonymous"
		if u, ok := c.Locals("username").(string); ok {
			username = u
		}

		err := upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			handle(c.Context(), conn, agentID, username, manager)
		})
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return fiber.ErrUpgradeRequired
		}

		return nil
	}
}

func handle(ctx context.Context, conn *websocket.Conn, agentID, username string, manager ws.Manager) {
	client := &ws.Client{
		ID:       agentID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
	}

	log.Printf("WebSocket: connected %s (user: %s) | total: %d",
		agentID, username, manager.ClientCount()+1)

	manager.Register(client)

	defer func() {
		manager.Unregister(client)
		conn.Close()
		log.Printf("WebSocket: disconnected %s (user: %s) | total: %d",
			agentID, username, manager.ClientCount())
	}()

	conn.SetPingHandler(func(appData string) error {
		return conn.WriteMessage(websocket.TextMessage, []byte(appData))
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("WebSocket error for %s: %v", agentID, err)
			}
			break
		}
	}
}
