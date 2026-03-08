package api

type ShutdownResponse struct {
	Status bool   `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}
