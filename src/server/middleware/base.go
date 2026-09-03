package middleware

import (
	"net/http"
	"strings"
)

func Base(basepath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if basepath == "" {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == basepath {
				http.Redirect(w, r, basepath+"/", http.StatusFound)
				return
			}
			if !strings.HasPrefix(r.URL.Path, basepath) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			r2 := r.Clone(r.Context())
			r2.URL.Path = strings.TrimPrefix(r.URL.Path, basepath)
			next.ServeHTTP(w, r2)
		})
	}
}
