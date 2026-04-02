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
	}{
		{
			name:           "invalid request",
			body:           nil,
			errorMsg:       "MSG",
			expectedStatus: http.StatusBadRequest,
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
			bodyStr := string(body)
			assert.Contains(t, bodyStr, tc.errorMsg, "error message should be in response")
			assert.Contains(t, bodyStr, "Error:", "should contain Error label")
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
