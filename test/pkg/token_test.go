package test

import (
	"reflect"
	"testing"

	"github.com/Troom-Corp/troom/internal/pkg"
)

// Тест который проверяет что функция Signjwt работает корректно
func TestSignjwt(t *testing.T) {
	var testuserid int
	testuserid = 100
	result := pkg.SignJWT(testuserid)
	if reflect.TypeOf(result) != reflect.TypeOf("") {
		t.Fatalf("Want type 'string', but got '%T'", reflect.TypeOf(result))
	}
}

// Тест который проверяет что функция GetIdentity работает корректно
func TestGetIdentity(t *testing.T) {
	var testuserid int
	testuserid = 100
	testjwttoken := pkg.SignJWT(testuserid)
	resultuserid, timejwt := pkg.GetIdentity(testjwttoken)
	if resultuserid != 100 {
		t.Fatalf("Want (100), but got (%d)", resultuserid)
	}
	if reflect.TypeOf(timejwt) != reflect.TypeOf(int64(0)) {
		t.Fatalf("Want type 'time.Time', but got '%T'", reflect.TypeOf(timejwt))
	}
}
