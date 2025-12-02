package rpc

import (
	"net/http"
)

// ExecutionService is the service for executing code.
type ExecutionService struct{}

// Execute is the method that will be exposed over JSON-RPC.
func (s *ExecutionService) Execute(r *http.Request, args *ExecuteRequest, reply *ExecuteResponse) error {
	// In a real implementation, you would execute the code here.
	// For now, we'll just return a dummy response.
	reply.ExecutionID = args.ExecutionID
	reply.Output = "Executed: " + args.Code
	return nil
}

// EchoService is the service for echoing messages.
type EchoService struct{}

// Echo is the method that will be exposed over JSON-RPC.
func (s *EchoService) Echo(r *http.Request, args *EchoRequest, reply *EchoResponse) error {
	reply.Message = args.Message
	return nil
}
