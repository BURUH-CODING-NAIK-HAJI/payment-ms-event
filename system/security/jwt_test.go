package security_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rizface/golang-api-template/app/entity/securityentity"
	"github.com/rizface/golang-api-template/system/security"
	"github.com/stretchr/testify/assert"
)

var bearer, refresh string

func TestEncodeDataToJwt(t *testing.T) {
	userData := &securityentity.UserData{
		Id:   "id",
		Name: "Fariz",
	}
	generatedResponseJwt := security.GenerateToken(userData)
	bearer = generatedResponseJwt.TokenSchema.Bearer
	refresh = generatedResponseJwt.TokenSchema.Refresh
	assert.Equal(t, "securityentity.GeneratedResponseJwt", reflect.TypeOf(generatedResponseJwt).String())
	assert.Equal(t, "string", reflect.TypeOf(generatedResponseJwt.TokenSchema.Bearer).String())
	assert.Equal(t, "string", reflect.TypeOf(generatedResponseJwt.TokenSchema.Refresh).String())
}

func TestDecodeBearerJwt(t *testing.T) {
	defer func() {
		err := recover()
		fmt.Println(err)
	}()
	claim := security.DecodeToken(bearer, "bearer")
	assert.Equal(t, "security.JwtClaim", reflect.TypeOf(claim).String())
}

func TestDecodeRefreshToken(t *testing.T) {
	defer func() {
		err := recover()
		fmt.Println(err)
	}()
	claim := security.DecodeToken(refresh, "refresh")
	assert.Equal(t, "security.JwtClaim", reflect.TypeOf(claim).String())
}
