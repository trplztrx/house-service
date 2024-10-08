package middleware

import (
	"house-service/internal/transport/handlers"
	"house-service/pkg"
	"net/http"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var respBoby []byte
		token := r.Header.Get("authorization")

		if token == "" {
			respBoby = handlers.CreateErrorResponse(r.Context(), handlers.NotAuthorizedError, handlers.NotAuthorizedErrorMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBoby)
			return
		}

		_, err := pkg.ValidateJWTToken(token)
		if err != nil {
			respBoby = handlers.CreateErrorResponse(r.Context(), handlers.NotAuthorizedError, handlers.NotAuthorizedErrorMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBoby)
			return
		}

		handler.ServeHTTP(w, r)
	})
}