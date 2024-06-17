package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

	shortURL, err := useCase.CreateShortURL(string(body), conf)
	if err != nil {
		if errors.Is(err, usecase.ErrShortURLAlreadyExists) {
			log.Println(usecase.ErrShortURLAlreadyExists)
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			return
		} else {
			log.Println("Ошибка при создании шорта")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	fullURL := ""

	if conf.BaseURL == "" {
		var scheme string
		if r.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		fullURL = scheme + r.Host + "/" + shortURL
	} else {
		fullURL = conf.BaseURL + "/" + shortURL
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullURL))
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

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	shortURL, err := useCase.CreateShortURL(req.URL, conf)
	if err != nil {
		if errors.Is(err, usecase.ErrShortURLAlreadyExists) {
			log.Println(usecase.ErrShortURLAlreadyExists)
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			return
		} else {
			log.Println("Ошибка при создании шорта")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	fullURL := ""

	if conf.BaseURL == "" {
		var scheme string
		if r.TLS != nil {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		fullURL = scheme + r.Host + "/" + shortURL
	} else {
		fullURL = conf.BaseURL + "/" + shortURL
	}

	resp := models.Response{
		Result: fullURL,
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

		shortURL, err := useCase.CreateShortURL(batch.OriginalURL, conf)
		if err != nil {
			if errors.Is(err, usecase.ErrShortURLAlreadyExists) {
				log.Println(usecase.ErrShortURLAlreadyExists)
				w.WriteHeader(http.StatusConflict) // 409 Conflict
				return
			} else {
				log.Println("Ошибка при создании шорта")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		fullURL := ""

		if conf.BaseURL == "" {
			var scheme string
			if r.TLS != nil {
				scheme = "https://"
			} else {
				scheme = "http://"
			}
			fullURL = scheme + r.Host + "/" + shortURL
		} else {
			fullURL = conf.BaseURL + "/" + shortURL
		}

		resp := models.ResponseBatch{
			ID:       batch.ID,
			ShortURL: fullURL,
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
