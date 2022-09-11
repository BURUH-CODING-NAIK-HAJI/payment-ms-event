package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rizface/golang-api-template/app/entity/securityentity"
	"github.com/rizface/golang-api-template/app/errorgroup"
	"github.com/rizface/golang-api-template/app/middleware"
	"github.com/rizface/golang-api-template/system/security"
	"github.com/stretchr/testify/assert"
)

func TestErrorMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := os.OpenFile("/notexistsfile.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	middleware.ErrorHandler(handler).ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAuthMiddleware(t *testing.T) {
	var token string
	var claim security.JwtClaim
	var ok bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, ok = r.Context().Value("user").(security.JwtClaim)
	})
	t.Run("Success decode token", func(t *testing.T) {
		userData := &securityentity.UserData{
			Id:   "id",
			Name: "Fariz",
		}
		generatedResponseJwt := security.GenerateToken(userData)
		token = generatedResponseJwt.TokenSchema.Bearer

		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Add("Authorization", "bearer "+token)
		w := httptest.NewRecorder()
		middleware.AuthHandler(handler).ServeHTTP(w, r)
		assert.True(t, ok)
		assert.Equal(t, userData.Id, claim.UserData.Id)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		defer func() {
			err, ok := recover().(errorgroup.Error)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, err.Code)
		}()
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Add("Authorization", "bearer "+"invalid token")
		w := httptest.NewRecorder()
		middleware.AuthHandler(handler).ServeHTTP(w, r)
	})

	t.Run("Non bearer token", func(t *testing.T) {
		defer func() {
			err, ok := recover().(errorgroup.Error)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, err.Code)
		}()
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Add("Authorization", "bearar "+"invalid token")
		w := httptest.NewRecorder()
		middleware.AuthHandler(handler).ServeHTTP(w, r)
	})

	t.Run("Invalid token segment length", func(t *testing.T) {
		defer func() {
			err, ok := recover().(errorgroup.Error)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, err.Code)
		}()
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Add("Authorization", "bearer"+token)
		w := httptest.NewRecorder()
		middleware.AuthHandler(handler).ServeHTTP(w, r)
	})
}
