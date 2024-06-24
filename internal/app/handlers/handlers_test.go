package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// общая функция для отправки запросов
func testRequest(t *testing.T, ts *httptest.Server, method,
	path, body string) (*http.Response, string) {

	// Создаем кастомный клиент
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Отключаем автоматический редирект
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)

}

func TestRoute(t *testing.T) {
	db := maps.NewMapStorage()
	storage := usecase.NewShortURLUseCase(db)

	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started")
	conf := config.ParseConfig()

	ts := httptest.NewServer(InitRoutes(*storage, conf))
	defer ts.Close()
	//добавим в бд тестовую запись
	id := "admin1"
	_ = db.Create("SHORTurl", (ts.URL + "/SHORTurl"), "https://vk.com", id)
	_ = db.Create("SHORTurl_2", (ts.URL + "/SHORTurl_2"), ".ru", id)

	oneTest, _ := db.ReadOriginal("https://vk.com")
	twoTest, _ := db.ReadOriginal(".ru")

	type values struct {
		url    string
		method string
		body   string
	}

	type want struct {
		status int
		body   string
	}

	var tests = []struct {
		name   string
		values values
		want   want
	}{
		{
			name: "Test #1 Запись в бд",
			values: values{
				url:    "/",
				method: "POST",
				body:   "https://ya.ru",
			},
			want: want{
				status: http.StatusCreated,
			},
		},
		{
			name: "Test #2 Метод GET вместо POST",
			values: values{
				url:    "/",
				method: "GET",
				body:   "https://ya.ru",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #3 без тела запроса",
			values: values{
				url:    "/",
				method: "POST",
				body:   "",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #4 не валидный URL",
			values: values{
				url:    "/",
				method: "POST",
				body:   "ggf.fdfhk/fsdf",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #5 добавление дубликата",
			values: values{
				url:    "/",
				method: "POST",
				body:   "https://vk.com",
			},
			want: want{
				status: http.StatusConflict,
				body:   ts.URL + "/SHORTurl",
			},
		},
		{
			name: "Test #6 проверка извлечения URL по сокращенной ссылке",
			values: values{
				url:    "/SHORTurl",
				method: "GET",
				body:   "",
			},
			want: want{
				status: http.StatusTemporaryRedirect,
				body:   "https://vk.com",
			},
		},
		{
			name: "Test #7 метод POST вмсето GET",
			values: values{
				url:    "/" + oneTest,
				method: "POST",
				body:   "",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #8 несуществующая сокращенная ссылка",
			values: values{
				url:    "/sdfjvu88934nkdkl",
				method: "GET",
				body:   "",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #9 оригинальная ссылка ' ' ",
			values: values{
				url:    "/" + twoTest,
				method: "GET",
				body:   "",
			},
			want: want{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test #10 не тот метод HandlerApiPost",
			values: values{
				url:    "/api/shorten",
				method: "GET",
				body:   `{"url":"https://ya.ru/"}`,
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #11 не корректный JSON HandlerApiPost",
			values: values{
				url:    "/api/shorten",
				method: "POST",
				body:   `{"https://ya.ru/"}`,
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "Test #12 не валидный url HandlerApiPost",
			values: values{
				url:    "/api/shorten",
				method: "POST",
				body:   `{"url":"sobaka"}`,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #13 валидный url HandlerApiPost",
			values: values{
				url:    "/api/shorten",
				method: "POST",
				body:   `{"url":"https://vks.com"}`,
			},
			want: want{
				status: http.StatusCreated,
			},
		},
	}
	for num, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, url := testRequest(t, ts, test.values.method, test.values.url, test.values.body)
			defer resp.Body.Close()

			assert.Equal(t, test.want.status, resp.StatusCode)
			if num == 4 {
				assert.Equal(t, test.want.body, url)
			} else {
				assert.Equal(t, test.want.body, resp.Header.Get("Location"))
			}

		})

	}

	requestBody := `{
        "url": "https://vk.com"
    }`

	// ожидаемое содержимое тела ответа при успешном запросе
	successBody := fmt.Sprintf(`{"result": "%s/SHORTurl"}`, ts.URL)

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)

		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", ts.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)

		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", ts.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// Если ответ сжат, декомпрессируем его
		if resp.Header.Get("Content-Encoding") == "gzip" {
			zr, err := gzip.NewReader(bytes.NewReader(b))
			require.NoError(t, err)
			defer zr.Close()
			b, err = io.ReadAll(zr)
			require.NoError(t, err)
		}
		require.JSONEq(t, successBody, string(b))
	})

	conf = &config.Config{
		BaseURL: "http://localhost",
	}

	t.Run("batch request", func(t *testing.T) {
		reqBody := []models.RequestBatch{
			{ID: "1", OriginalURL: "https://example.com"},
			{ID: "2", OriginalURL: "https://example.org"},
		}
		reqBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBytes))
		req = req.WithContext(context.WithValue(req.Context(), "userID", "user123"))
		w := httptest.NewRecorder()

		HandlerBatch(w, req, *storage, conf)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var resBody []models.ResponseBatch
		err := json.NewDecoder(resp.Body).Decode(&resBody)
		require.NoError(t, err)
		require.Len(t, resBody, 2)
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/shorten/batch", nil)
		w := httptest.NewRecorder()

		HandlerBatch(w, req, *storage, conf)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("missing userID", func(t *testing.T) {
		reqBody := []models.RequestBatch{
			{ID: "1", OriginalURL: "https://example.com"},
		}
		reqBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBytes))
		w := httptest.NewRecorder()

		HandlerBatch(w, req, *storage, conf)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer([]byte("invalid json")))
		req = req.WithContext(context.WithValue(req.Context(), "userID", "user123"))
		w := httptest.NewRecorder()

		HandlerBatch(w, req, *storage, conf)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
