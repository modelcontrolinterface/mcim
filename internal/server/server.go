package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func NewRouter() *mux.Router {
	s := rpc.NewServer()
	serverManagerService := new(ServerManagerService)

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(serverManagerService, "")

	r := mux.NewRouter()
	r.Handle("/", s)

	return r
}
