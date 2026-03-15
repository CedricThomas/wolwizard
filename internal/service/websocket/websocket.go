package websocket

type Manager interface {
	Broadcast(message []byte) error
}
