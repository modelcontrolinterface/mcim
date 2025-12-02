package rpc

import "net/http"

type ExecutionService struct{}

func (s *ExecutionService) Execute(r *http.Request, args *ExecuteRequest, reply *ExecuteResponse) error {
	out := "Executed: " + args.Code

	reply.ExecutionID = args.ExecutionID
	reply.Output = &out
	reply.StdErr = nil
	reply.StdOut = nil

	return nil
}
