package test

import (
	"reflect"
	"testing"

	"github.com/Troom-Corp/troom/internal/pkg"
)

// Тест который проверяет что функция CreateRefreshToken работает корректно
func TestCreateRefreshToken(t *testing.T) {
	var testuserid int
	testuserid = 100
	gotResult, gotError := pkg.CreateRefreshToken(testuserid)
	if gotError != nil {
		t.Fatalf("Want error (nil), but got (%e)", gotError)
	}
	if reflect.TypeOf(gotResult) != reflect.TypeOf("") {
		t.Fatalf("Want type 'string', but got '%T'", reflect.TypeOf(gotResult))
	}
}

// Тест который проверяет что функция CreateRefreshToken работает корректно
func TestCreateAccessToken(t *testing.T) {
	var testuserid int
	testuserid = 100
	gotResult, gotError := pkg.CreateAccessToken(testuserid)
	if gotError != nil {
		t.Fatalf("Want error (nil), but got (%e)", gotError)
	}
	if reflect.TypeOf(gotResult) != reflect.TypeOf("") {
		t.Fatalf("Want type 'string', but got '%T'", reflect.TypeOf(gotResult))
	}
}

// Тест который проверяет что функция GetIdentity работает корректно для refresh token
func TestGetIdentityForRerfreshToken(t *testing.T) {
	var testuserid int
	testuserid = 100
	testjwttoken, _ := pkg.CreateRefreshToken(testuserid)
	resultuserid, timejwt := pkg.GetIdentity(testjwttoken)
	if resultuserid != 100 {
		t.Fatalf("Want (100), but got (%d)", resultuserid)
	}
	if reflect.TypeOf(timejwt) != reflect.TypeOf(int64(0)) {
		t.Fatalf("Want type 'time.Time', but got '%T'", reflect.TypeOf(timejwt))
	}
}

// Тест который проверяет что функция GetIdentity работает корректно для access token
func TestGetIdentityForAccessToken(t *testing.T) {
	var testuserid int
	testuserid = 100
	testjwttoken, _ := pkg.CreateAccessToken(testuserid)
	resultuserid, timejwt := pkg.GetIdentity(testjwttoken)
	if resultuserid != 100 {
		t.Fatalf("Want (100), but got (%d)", resultuserid)
	}
	if reflect.TypeOf(timejwt) != reflect.TypeOf(int64(0)) {
		t.Fatalf("Want type 'time.Time', but got '%T'", reflect.TypeOf(timejwt))
	}
}
