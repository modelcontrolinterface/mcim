package rpc

import "time"

type Action string

const (
	ActionRun        Action = "run"
	ActionSave       Action = "save"
	ActionSaveOnRun  Action = "save_on_run"
	ActionSaveAndRun Action = "save_and_run"
)

type ExecuteRequest struct {
	ExecutionID string         `json:"executionID"`
	Code        string         `json:"code"`
	Environment *string        `json:"env,omitempty"`
	Filename    *string        `json:"filename,omitempty"`
	Action      *Action        `json:"action,omitempty"`
	Timeout     *time.Duration `json:"timeout,omitempty"`
}

type ExecuteResponse struct {
	ExecutionID string  `json:"executionID"`
	Output      *string `json:"output,omitempty"`
	StdErr      *string `json:"stderr,omitempty"`
	StdOut      *string `json:"stdout,omitempty"`
}
