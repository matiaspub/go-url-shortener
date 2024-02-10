package redirect_test

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/server/handlers/url/redirect"
	"url-shortener/internal/server/handlers/url/redirect/mocks"
	"url-shortener/internal/storage"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name     string
		alias    string
		err      error
		data     map[string]string
		expected string
		status   int
	}{
		{"alias found", "test", nil, map[string]string{"test": "https://example.com"}, "https://example.com", http.StatusFound},
		{"alias not found", "fake_test", storage.ErrURLNotFound, map[string]string{}, "not found", http.StatusOK},
		{"internal server error", "fake_test", errors.New("simulate internal error"), map[string]string{}, "Internal Server Error", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockURLGetter := mocks.NewURLGetter(t)
			mockURLGetter.On("GetURL", tc.alias).Return(tc.data[tc.alias], tc.err)

			mux := chi.NewRouter()
			mux.HandleFunc("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), mockURLGetter))
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/"+tc.alias, nil)
			mux.ServeHTTP(recorder, request)

			assert.Equal(t, tc.status, recorder.Code)
			if tc.status == http.StatusFound {
				assert.Equal(t, tc.expected, recorder.Header().Get("location"))
			} else {
				assert.Contains(t, recorder.Body.String(), tc.expected)
			}
		})
	}
}
