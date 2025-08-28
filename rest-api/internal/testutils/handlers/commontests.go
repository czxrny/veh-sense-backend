package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func RunBasicGetAllTests(url string, handlerFunc func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	tests := []HttpTestStruct{
		{"GET ALL no token", http.MethodGet, "", url, nil, nil, http.StatusUnauthorized},
		{"GET ALL ok", http.MethodGet, "", url, GetAdminAuth(), nil, http.StatusOK},
		{"GET ALL with body", http.MethodGet, "", url, GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
	}
	RunTestCases(tests, handlerFunc, t)
}

func RunBasicAddTests(url, okRequestBody, badRequestBody string, handlerFunc func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	tests := []HttpTestStruct{
		{"POST no token", http.MethodPost, "", url, nil, strings.NewReader(okRequestBody), http.StatusUnauthorized},
		{"POST ok", http.MethodPost, "", url, GetAdminAuth(), strings.NewReader(okRequestBody), http.StatusOK},
		{"POST no request body", http.MethodPost, "", url, GetAdminAuth(), nil, http.StatusBadRequest},
		{"POST bad request body", http.MethodPost, "", url, GetAdminAuth(), strings.NewReader(badRequestBody), http.StatusBadRequest},
	}
	RunTestCases(tests, handlerFunc, t)
}

func RunBasicUpdateTests(url, okRequestBody, badRequestBody string, handlerFunc func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	tests := []HttpTestStruct{
		{"PATCH no token", http.MethodPatch, "1", url, nil, strings.NewReader(okRequestBody), http.StatusUnauthorized},
		{"PATCH ok", http.MethodPatch, "1", url, GetAdminAuth(), strings.NewReader(okRequestBody), http.StatusOK},
		{"PATCH no request body", http.MethodPatch, "1", url, GetAdminAuth(), nil, http.StatusBadRequest},
		{"PATCH bad request body", http.MethodPatch, "1", url, GetAdminAuth(), strings.NewReader(badRequestBody), http.StatusBadRequest},
		{"PATCH bad id in path", http.MethodPatch, "invalid", url, GetAdminAuth(), strings.NewReader(okRequestBody), http.StatusBadRequest},
	}
	RunTestCases(tests, handlerFunc, t)
}

func RunBasicGetByIdTests(url string, handler func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	tests := []HttpTestStruct{
		{"GET no token", http.MethodGet, "1", url, nil, nil, http.StatusUnauthorized},
		{"GET ok", http.MethodGet, "1", url, GetAdminAuth(), nil, http.StatusOK},
		{"GET with body", http.MethodGet, "1", url, GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
		{"GET bad id in path", http.MethodGet, "invalid", url, GetAdminAuth(), nil, http.StatusBadRequest},
	}
	RunTestCases(tests, handler, t)
}

func RunBasicDeleteByIdTests(url string, handler func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	tests := []HttpTestStruct{
		{"DELETE no token", http.MethodDelete, "1", url, nil, nil, http.StatusUnauthorized},
		{"DELETE ok", http.MethodDelete, "1", url, GetAdminAuth(), nil, http.StatusNoContent},
		{"DELETE with body", http.MethodDelete, "1", url, GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
		{"DELETE bad id in path", http.MethodDelete, "invalid", url, GetAdminAuth(), nil, http.StatusBadRequest},
	}
	RunTestCases(tests, handler, t)
}

func RunTestCases(tc []HttpTestStruct, handlerFunc func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	for _, testCase := range tc {
		t.Run(testCase.Name, func(t *testing.T) {
			req := CreateNewRequest(testCase.Method, testCase.Url, testCase.AuthInfo, testCase.Body)
			w := httptest.NewRecorder()

			if testCase.ID != "" {
				req = AddChiIdToContext(req, testCase.ID)
			}

			handlerFunc(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.WantStatus {
				t.Errorf("expected %d, got %d", testCase.WantStatus, res.StatusCode)
			}
		})
	}
}
