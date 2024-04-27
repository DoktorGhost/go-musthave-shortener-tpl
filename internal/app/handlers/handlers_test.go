package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerPost(t *testing.T) {
	type want struct {
		code        int
		contentType string
		url         string
	}
	type test struct {
		name   string
		method string
		body   *strings.Reader
		want   want
	}
	tests := []test{
		{
			name:   "not POST method",
			method: http.MethodGet,
			body:   strings.NewReader("https://yandex.ru"),
			want: want{
				code: 405,
			},
		},
		{
			name:   "body is nil",
			method: http.MethodPost,
			body:   strings.NewReader(""),
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name:   "normal",
			method: http.MethodPost,
			body:   strings.NewReader("ya.ru"),
			want: want{
				code:        201,
				contentType: "text/plain",
				url:         "ya.ru/shortURL",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", tt.body)
			w := httptest.NewRecorder()
			HandlerPost(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.url, w.Body.String())
		})
	}
}

func TestHandlerGet(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	type test struct {
		name   string
		method string
		target string
		want   want
	}
	tests := []test{
		{
			name:   "not Get method",
			method: http.MethodPost,
			target: "yandex.ru",
			want: want{
				code: 405,
			},
		},
		{
			name:   "normal",
			method: http.MethodGet,
			target: "yandex.ru",
			want: want{
				code:     307,
				location: "https://yandex.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/"+tt.target, nil)
			w := httptest.NewRecorder()
			HandlerGet(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.location, w.Header().Get("Location"))
		})
	}
}
