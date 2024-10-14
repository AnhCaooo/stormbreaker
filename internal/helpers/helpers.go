// AnhCao 2024
package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// encode response bodies
func EncodeResponse[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// TODO: add unit test
// decode the request bodies
func DecodeRequest[T any](r *http.Request) (v T, err error) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// TODO: add unit test
// decode the response bodies
func DecodeResponse[T any](r *http.Response) (v T, err error) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// get current directory (from root to this repo only)
func GetCurrentDir() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %s", err.Error())
	}
	return
}
