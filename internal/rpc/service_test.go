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

	// Create a request to pass to our handler.
	reqBody := map[string]interface{}{
		"method": "ExecutionService.Execute",
		"params": []map[string]string{{"executionID": "123", "code": "fmt.Println(\"hello world\")"}},
		"id":     "1",
	}
	jsonReq, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(router)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
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

func TestEchoService_Echo(t *testing.T) {
	router := NewRouter()

	// Create a request to pass to our handler.
	reqBody := map[string]interface{}{
		"method": "EchoService.Echo",
		"params": []map[string]string{{"message": "hello"}},
		"id":     "1",
	}
	jsonReq, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(router)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	result := resp["result"].(map[string]interface{})
	if result["message"] != "hello" {
		t.Errorf("handler returned unexpected message: got %v want %v",
			result["message"], "hello")
	}
}
