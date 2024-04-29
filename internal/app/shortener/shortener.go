package shortener

import (
	"math/rand"
	"time"
)

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	scr := rand.NewSource(time.Now().UnixNano())
	r := rand.New(scr)
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[r.Intn(len(charset))])
	}
	return string(result)
}
