package rpc

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

// NewRouter creates a new RPC router.
func NewRouter() *mux.Router {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")

	s.RegisterService(new(ExecutionService), "")
	s.RegisterService(new(EchoService), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)
	return r
}
