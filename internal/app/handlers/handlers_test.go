package handlers

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
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

	//добавим в бд тестовую запись
	oneTest := db.Create("SHORTurl", "https://vk.com")
	twoTest := db.Create("SHORTurl_2", ".ru")
	ts := httptest.NewServer(InitRoutes(*storage))
	defer ts.Close()

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
				status: http.StatusCreated,
				body:   ts.URL + "/SHORTurl",
			},
		},
		{
			name: "Test #6 проверка извлечения URL по сокращенной ссылке",
			values: values{
				url:    "/" + oneTest,
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
				body:   `{"url":"https://vk.com"}`,
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
}
