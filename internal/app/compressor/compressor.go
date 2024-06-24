package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipResponseWriter обертывает http.ResponseWriter для поддержки сжатия данных.
type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write перенаправляет запись в сжатый Writer.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipAndDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, содержит ли заголовок Content-Encoding информацию о сжатии
		if r.Header.Get("Content-Encoding") == "gzip" {
			// Декомпрессия тела запроса
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer zr.Close()
			r.Body = http.MaxBytesReader(w, zr, 10<<20) // Устанавливаем максимальный размер тела запроса
		}

		// Проверяем поддержку клиентом сжатия ответов
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Создаем сжатый Writer для записи ответа
			gz := gzip.NewWriter(w)
			defer gz.Close()

			// Создаем ResponseWriter с оберткой для сжатия данных
			gzWriter := &gzipWriter{Writer: gz, ResponseWriter: w}

			// Устанавливаем заголовок Content-Encoding для указания, что ответ будет сжат
			gzWriter.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gzWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
