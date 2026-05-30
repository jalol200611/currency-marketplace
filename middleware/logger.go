package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		fmt.Printf(
			"[REQUEST] %s %s\n",
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)

		fmt.Printf(
			"[DONE] %s | %v\n",
			r.URL.Path,
			time.Since(start),
		)
	})
}
