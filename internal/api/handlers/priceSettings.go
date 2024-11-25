package handlers

import (
	"fmt"
	"net/http"
)

// GetPriceSettings retrieves the price settings for specified user
func GetPriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// CreatePriceSettings creates a new price settings for new user
func CreatePriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// PatchPriceSettings updates the price settings for specified user
func PatchPriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// DeletePriceSettings deletes the price settings when user was deleted or removed
func DeletePriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
