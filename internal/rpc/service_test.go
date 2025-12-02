package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecutionService_Execute(t *testing.T) {
	router := NewRouter()

	reqBody := map[string]interface{}{
		"method": "ExecutionService.Execute",
		"params": []map[string]string{{"executionID": "123", "code": "fmt.Println(\"hello world\")"}},
		"id":     "1",
	}
	jsonReq, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.Handler(router)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	result := resp["result"].(map[string]interface{})
	if result["executionID"] != "123" {
		t.Errorf("handler returned unexpected executionID: got %v want %v",
			result["executionID"], "123")
	}
	expectedOutput := "Executed: fmt.Println(\"hello world\")"
	if result["output"] != expectedOutput {
		t.Errorf("handler returned unexpected output: got %v want %v",
			result["output"], expectedOutput)
	}
}
