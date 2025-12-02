package rpc

// ExecuteRequest is the request for the ExecutionService.Execute method.
type ExecuteRequest struct {
	ExecutionID string `json:"executionID"`
	Code        string `json:"code"`
}

// ExecuteResponse is the response for the ExecutionService.Execute method.
type ExecuteResponse struct {
	ExecutionID string `json:"executionID"`
	Output      string `json:"output"`
}

// EchoRequest is the request for the EchoService.Echo method.
type EchoRequest struct {
	Message string `json:"message"`
}

// EchoResponse is the response for the EchoService.Echo method.
type EchoResponse struct {
	Message string `json:"message"`
}
