package websocket

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Conn     interface{}
	Send     chan []byte
	Username string
}

// Manager provides real-time message broadcasting to connected WebSocket clients
type Manager interface {
	// Broadcast sends a message to all connected clients
	// Returns error if broadcast buffer is full
	Broadcast(msg []byte) error

	// ClientCount returns the number of currently connected clients
	ClientCount() int

	// Start launches the manager's internal goroutines
	Start()

	// Shutdown gracefully closes all connections and stops manager
	Shutdown()

	// Register adds a new client to the manager
	Register(client *Client)

	// Unregister removes a client from the manager
	Unregister(client *Client)
}
