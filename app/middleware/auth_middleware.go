package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rizface/golang-api-template/app/errorgroup"
	"github.com/rizface/golang-api-template/system/security"
)

type UserContext string

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenSegment := strings.Split(authHeader, " ")
		if len(tokenSegment) != 2 || tokenSegment[0] != "bearer" {
			panic(errorgroup.TOKEN_INVALID)
		}

		claims := security.DecodeToken(tokenSegment[1], "bearer")
		r = r.WithContext(context.WithValue(r.Context(), "user", claims))
		r.Header.Add("Test", "Asd")
		next.ServeHTTP(w, r)
	})
}
