package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendErrorMsg(t *testing.T) {
	testCases := []struct {
		name           string
		errorMsg       string
		body           io.Reader
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid request",
			body:           nil,
			errorMsg:       "MSG",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `<div class="error-msg" role="alert"><strong>Error: </strong> <span>MSG</span></div>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", tc.body)

			sendErrorMsg(w, r, tc.errorMsg)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode,
				"unexpected status code. Expected :%d, got: %d", tc.expectedStatus, w.Result().StatusCode)

			body, _ := io.ReadAll(w.Result().Body)
			assert.Equal(t, tc.expectedBody, string(body), "expected response body")
		})
	}
}

func Test_Page404(t *testing.T) {
	testCases := []struct {
		name           string
		body           io.Reader
		expectedStatus int
	}{
		{
			name:           "not found page",
			body:           nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", tc.body)

			handler := &Handlers{}
			handler.Page404(w, r)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode,
				"unexpected status code. Expected :%d, got: %d", tc.expectedStatus, w.Result().StatusCode)
		})
	}
}
