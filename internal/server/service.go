package server

import (
	"fmt"
	"net/http"
)

type ServerManagerService struct{}

func (s *ServerManagerService) Connect(r *http.Request, args *ConnectRequest, reply *ConnectResponse) error {
	reply.Message = fmt.Sprintf("Successfully connected to server ID: %s", args.ServerID)
	return nil
}

func (s *ServerManagerService) Disconnect(r *http.Request, args *DisconnectRequest, reply *DisconnectResponse) error {
	reply.Message = fmt.Sprintf("Successfully disconnected from server ID: %s", args.ServerID)
	return nil
}
