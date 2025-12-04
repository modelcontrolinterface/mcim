package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecutionService_Execute(t *testing.T) {
	reqBody := map[string]any{
		"id":     "1",
		"method": "ExecutionService.Execute",
		"params": []map[string]string{{
			"executionID": "123",
			"code":        "fmt.Println(\"hello world\")",
		}},
	}
	jsonReq, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	router := NewRouter()
	rr := httptest.NewRecorder()

	handler := http.Handler(router)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	result := resp["result"].(map[string]any)

	expectedExecutionID := "123"
	if result["executionID"] != expectedExecutionID {
		t.Errorf("handler returned unexpected executionID: got %v wanted %v", result["executionID"], expectedExecutionID)
	}

	expectedOutput := "fmt.Println(\"hello world\")"
	if result["output"] != expectedOutput {
		t.Errorf("handler returned unexpected output: got %v wanted %v", result["output"], expectedOutput)
	}
}
