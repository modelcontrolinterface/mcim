package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func NewRouter() *mux.Router {
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(ExecutionService), "")

	r := mux.NewRouter()
	r.Handle("/", s)

	return r
}
