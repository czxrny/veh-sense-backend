package common

import (
	"context"
	"encoding/json"
	"net/http"
)

func GetAllHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context) ([]T, error)) {
	items, err := innerHandler(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetByIdHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) (*T, error)) {
	id, err := getIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	item, err := innerHandler(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

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

func DeleteHandler(w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) error) {
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
