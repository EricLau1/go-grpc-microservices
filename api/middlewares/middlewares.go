package middlewares

import (
	"log"
	"microservices/api/restutil"
	"microservices/security"
	"net/http"
	"time"
)

func LogRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next(w, r)
		log.Printf(`{"proto": "%s", "method": "%s", "route": "%s%s", "request_time": "%v"}`,
			r.Proto, r.Method, r.Host, r.URL.Path, time.Since(t))
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := security.ExtractToken(r)
		if err != nil {
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		token, err := security.ParseToken(tokenString)
		if err != nil {
			log.Println("error on parse token:", err.Error())
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("invalid token", tokenString)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}

		next(w, r)
	}
}
