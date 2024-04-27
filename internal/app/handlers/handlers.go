package handlers

import (
	"io"
	"net/http"
)

func HandlerPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := string(body) + "/shortURL"
	//shortURL := createSgortURL(string(body))  тут будет логика

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func HandlerGet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Path[1:]
	originalURL := "https://" + id //тут будет логика

	if originalURL != "" {

		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
