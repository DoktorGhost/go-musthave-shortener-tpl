package auth

import (
	"context"
	"fmt"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/shortener"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type contextKey string

const UserIDKey contextKey = "userID"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// /ЗАМЕНИТЬ ПЕРЕМЕННЫМИ ОКРУЖЕНИЯ
const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: shortener.RandomString(5),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	// возвращаем строку токена
	return tokenString, nil
}

// извлекаем UserID и проверяем токен на валидность
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

func SetUserCookie(w http.ResponseWriter, token string) {

	cookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Cookie only accessible via HTTP(S), not JavaScript
		Secure:   true, // Cookie sent only over HTTPS
	}

	http.SetCookie(w, cookie)
}

func GetUserCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		return "", err
	}
	userID, err := GetUserID(cookie.Value)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserCookie(r)
		if err != nil {
			token, err := BuildJWTString()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			SetUserCookie(w, token)
			userID, _ := GetUserID(token)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
