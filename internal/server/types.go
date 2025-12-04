package server

type ConnectRequest struct {
	ServerID string `json:"server_id"`
}

type ConnectResponse struct {
	Message string `json:"message"`
}

type DisconnectRequest struct {
	ServerID string `json:"server_id"`
}

type DisconnectResponse struct {
	Message string `json:"message"`
}
