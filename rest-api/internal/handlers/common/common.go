package common

import (
	"context"
	"encoding/json"
	"net/http"
)

// Checks if there is a request body, invokes the inner handler and writes the response.
func GetAllHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context) ([]T, error)) {
	if !requestBodyIsEmpty(r) {
		http.Error(w, "Request body should be empty", http.StatusBadRequest)
	}

	items, err := innerHandler(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Checks if there is a request body, reads ID from path, invokes the inner handler and writes the response.
func GetByIdHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) (*T, error)) {
	if !requestBodyIsEmpty(r) {
		http.Error(w, "Request body should be empty", http.StatusBadRequest)
	}

	id, err := getIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	item, err := innerHandler(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Same as GetByIdHandler - just skipping the ID part (used for the /me/* endpoints)
func GetSimpleHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context) (*T, error)) {
	if !requestBodyIsEmpty(r) {
		http.Error(w, "Request body should be empty", http.StatusBadRequest)
	}

	item, err := innerHandler(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Decodes and validates the request body, invokes the inner handler and writes the response.
func PostHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T) (*R, error)) {
	var newItem T
	if err := decodeAndValidateRequestBody(r, &newItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := innerHandler(r.Context(), &newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Decodes and validates the request body, reads the id from path, invokes the inner handler and writes the response.
func PatchHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T, int) (*R, error)) {
	id, err := getIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	var updatedItem T
	if err := decodeAndValidateRequestBody(r, &updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := innerHandler(r.Context(), &updatedItem, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Same as PatchByIdHandler - just skipping the ID part (used for the /me/* endpoints)
func PatchSimpleHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T) (*R, error)) {
	var updatedItem T
	if err := decodeAndValidateRequestBody(r, &updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := innerHandler(r.Context(), &updatedItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Checks if there is a request body, reads ID from path, invokes the inner handler and writes the StatusNoContent.
func DeleteHandler(w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) error) {
	if !requestBodyIsEmpty(r) {
		http.Error(w, "Request body should be empty", http.StatusBadRequest)
	}

	id, err := getIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err := innerHandler(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
