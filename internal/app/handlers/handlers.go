package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/auth"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"net/http"
	"time"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, useCase usecase.ShortURLUseCase, conf *config.Config) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем userID из контекста
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	short, err := useCase.GetShort(originalURL)

	if err == nil {
		log.Println(usecase.ErrShortURLAlreadyExists)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		shortURL, _ := useCase.GetShortURL(short)
		w.Write([]byte(shortURL))
		return
	}

	short, err = useCase.GenerateShort(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := ""

	if conf.BaseURL == "" {
		var scheme string
		if r.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		shortURL = scheme + r.Host + "/" + short
	} else {
		shortURL = conf.BaseURL + "/" + short
	}

	err = useCase.CreateShortURL(short, shortURL, originalURL, userID, conf)

	if err != nil {
		log.Println("Ошибка при создании шорта", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func HandlerGet(res http.ResponseWriter, req *http.Request, useCase usecase.ShortURLUseCase) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(req, "shortURL")
	originalURL, err := useCase.GetOriginalURL(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(originalURL) > 4 {
		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func HandlerAPIPost(w http.ResponseWriter, r *http.Request, useCase usecase.ShortURLUseCase, conf *config.Config) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	originalURL := req.URL
	short, err := useCase.GetShort(originalURL)
	if err == nil {
		log.Println(usecase.ErrShortURLAlreadyExists)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		shortURL, _ := useCase.GetShortURL(short)
		resp := models.Response{Result: shortURL}
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	short, err = useCase.GenerateShort(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := ""

	if conf.BaseURL == "" {
		var scheme string
		if r.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		shortURL = scheme + r.Host + "/" + short
	} else {
		shortURL = conf.BaseURL + "/" + short
	}

	resp := models.Response{
		Result: shortURL,
	}

	err = useCase.CreateShortURL(short, shortURL, req.URL, userID, conf)
	if err != nil {
		log.Println("Ошибка при создании шорта", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func HandlerPing(res http.ResponseWriter, req *http.Request, conf *config.Config) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ps := conf.DatabaseDSN

	db, err := sql.Open("pgx", ps)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func HandlerBatch(w http.ResponseWriter, r *http.Request, useCase usecase.ShortURLUseCase, conf *config.Config) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req []models.RequestBatch
	var res []models.ResponseBatch
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for _, batch := range req {
		if batch.ID == "" || batch.OriginalURL == "" {
			continue
		}

		short, err := useCase.GenerateShort(batch.OriginalURL)
		if err != nil {
			log.Println("Ошибка при создании шорта", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := ""

		if conf.BaseURL == "" {
			var scheme string
			if r.TLS != nil {
				scheme = "https://"
			} else {
				scheme = "http://"
			}
			shortURL = scheme + r.Host + "/" + short
		} else {
			shortURL = conf.BaseURL + "/" + short
		}

		err = useCase.CreateShortURL(short, shortURL, batch.OriginalURL, userID, conf)

		if err != nil {
			continue
		}

		resp := models.ResponseBatch{
			ID:       batch.ID,
			ShortURL: shortURL,
		}

		res = append(res, resp)

	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func HandlerGetUserURL(w http.ResponseWriter, r *http.Request, useCase usecase.ShortURLUseCase) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем userID из контекста
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	urls, err := useCase.GetUserURL(userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
