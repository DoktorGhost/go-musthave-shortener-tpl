package server

import (
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartServer(t *testing.T) {
	// Создание фейкового HTTP-сервера
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Получение порта фейкового сервера
	_, port, err := net.SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to parse server address: %v", err)
	}

	// Запуск StartServer с адресом фейкового сервера
	err = StartServer(port)
	require.NoError(t, err)

}
