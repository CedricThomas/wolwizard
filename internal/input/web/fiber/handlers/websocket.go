package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"

	ws "github.com/CedricThomas/console/internal/service/websocket"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(manager ws.Manager, authCtrl controller.Auth) fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Params("sessionID")

		err := upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			handle(c.Context(), conn, sessionID, authCtrl, manager)
		})
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return fiber.ErrUpgradeRequired
		}

		return nil
	}
}

type authMessage struct {
	Token string `json:"token"`
}

func handle(ctx context.Context, conn *websocket.Conn, sessionID string, authCtrl controller.Auth, manager ws.Manager) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	msgType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("WebSocket: failed to read auth for %s: %v", sessionID, err)
		conn.Close()
		return
	}

	if msgType != websocket.TextMessage {
		log.Printf("WebSocket: expected text message for auth, got type %d", msgType)
		conn.Close()
		return
	}

	var authMsg authMessage
	if err := json.Unmarshal(message, &authMsg); err != nil {
		log.Printf("WebSocket: invalid auth message for %s: %v", sessionID, err)
		conn.Close()
		return
	}

	if authMsg.Token == "" {
		log.Printf("WebSocket: missing token for %s", sessionID)
		conn.Close()
		return
	}

	username, err := authCtrl.ValidateToken(ctx, authMsg.Token)
	if err != nil {
		log.Printf("WebSocket: invalid token for %s: %v", sessionID, err)
		conn.Close()
		return
	}

	log.Printf("WebSocket: auth successful for %s (user: %s)", sessionID, username)

	client := &ws.Client{
		ID:       sessionID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
	}

	log.Printf("WebSocket: connected %s (user: %s) | total: %d",
		sessionID, username, manager.ClientCount()+1)

	manager.Register(client)

	conn.SetReadDeadline(time.Time{})

	defer func() {
		manager.Unregister(client)
		conn.Close()
		log.Printf("WebSocket: disconnected %s (user: %s) | total: %d",
			sessionID, username, manager.ClientCount())
	}()

	conn.SetPingHandler(func(appData string) error {
		return conn.WriteMessage(websocket.TextMessage, []byte(appData))
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("WebSocket error for %s: %v", sessionID, err)
			}
			break
		}
	}
}
