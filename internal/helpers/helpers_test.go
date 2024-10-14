// AnhCao 2024
package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncodeResponse(t *testing.T) {
	// Define a sample struct or data type for parameter "v"
	type MyData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testData := MyData{Name: "Alice", Age: 30}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the EncodeResponse function with the recorder, status code, and data
	err := EncodeResponse(recorder, http.StatusOK, testData)

	// Check for errors
	if err != nil {
		t.Errorf("EncodeResponse returned unexpected error: %v", err)
	}

	// Assert the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Assert the content type header
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type header 'application/json', got '%s'", contentType)
	}

	// Assert the response body
	var decodedData MyData
	err = json.NewDecoder(recorder.Body).Decode(&decodedData)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
		return
	}

	if decodedData != testData {
		t.Errorf("Expected data %v in response body, got %v", testData, decodedData)
	}
}
