// AnhCao 2024
package helpers

import (
	"encoding/json"
	"fmt"
)

// MapInterfaceToStruct converts an interface{} to a struct of a specified type.
// It first marshals the interface{} to JSON bytes and
// then unmarshals those bytes into the struct.
func MapInterfaceToStruct[T any](data interface{}) (*T, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interface to JSON: %w", err)
	}

	// Initialize the generic type
	var v T

	// Unmarshal the JSON bytes into the struct
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to PriceSettings: %w", err)
	}
	return &v, nil
}
