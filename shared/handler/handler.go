package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/czxrny/veh-sense-backend/shared/apierrors"
)

// Checks if there is a request body, invokes the inner handler and writes the response.
func GetAllHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, url.Values) ([]T, error)) {
	if !requestBodyIsEmpty(r) {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	items, err := innerHandler(r.Context(), r.URL.Query())
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Checks if there is a request body, reads ID from path, invokes the inner handler and writes the response.
func GetByIdHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) (*T, error)) {
	if !requestBodyIsEmpty(r) {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	id, err := getIdFromPath(r)
	if err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	item, err := innerHandler(r.Context(), id)
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Same as GetByIdHandler - just skipping the ID part (used for the /me/* endpoints)
func GetSimpleHandler[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context) (*T, error)) {
	if !requestBodyIsEmpty(r) {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	item, err := innerHandler(r.Context())
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Decodes and validates the request body, invokes the inner handler and writes the response.
func PostHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T) (*R, error)) {
	var newItem T
	if err := decodeAndValidateRequestBody(r, &newItem); err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	item, err := innerHandler(r.Context(), &newItem)
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Decodes and validates the request body, invokes the inner handler but does not return the response.
func PostHandlerSilent[T any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T) error) {
	var newItem T
	if err := decodeAndValidateRequestBody(r, &newItem); err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	if err := innerHandler(r.Context(), &newItem); err != nil {
		handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Decodes and validates the request body, reads the id from path, invokes the inner handler and writes the response.
func PatchHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T, int) (*R, error)) {
	id, err := getIdFromPath(r)
	if err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	var updatedItem T
	if err := decodeAndValidateRequestBody(r, &updatedItem); err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	item, err := innerHandler(r.Context(), &updatedItem, id)
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Same as PatchByIdHandler - just skipping the ID part (used for the /me/* endpoints)
func PatchSimpleHandler[T, R any](w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, *T) (*R, error)) {
	var updatedItem T
	if err := decodeAndValidateRequestBody(r, &updatedItem); err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	item, err := innerHandler(r.Context(), &updatedItem)
	if err != nil {
		handleErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Checks if there is a request body, reads ID from path, invokes the inner handler and writes the StatusNoContent.
func DeleteHandler(w http.ResponseWriter, r *http.Request, innerHandler func(context.Context, int) error) {
	if !requestBodyIsEmpty(r) {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	id, err := getIdFromPath(r)
	if err != nil {
		handleErrors(w, apierrors.ErrBadRequest)
		return
	}

	if err := innerHandler(r.Context(), id); err != nil {
		handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
