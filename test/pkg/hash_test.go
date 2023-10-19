package test

import (
	"reflect"
	"testing"

	"github.com/Troom-Corp/troom/internal/pkg"
	"golang.org/x/crypto/bcrypt"
)

// Тест который проверяет что функция Decode работает корректно
func TestDecode(t *testing.T) {
	testHashString := "Password"
	testHash, _ := bcrypt.GenerateFromPassword([]byte(testHashString), bcrypt.DefaultCost)
	gotError := pkg.Decode(testHash, []byte("Password"))
	if gotError != nil {
		t.Fatalf("Want (nil), but got (%e)", gotError)
	}
}

// Тест который проверяет что функция Encode работает корректно
func TestEncode(t *testing.T) {
	testPassword := []byte("Password")
	gotResult, gotError := pkg.Encode(testPassword)
	if len(gotResult) != 60 {
		t.Fatalf("Want lenght of slice (60), but got (%d)", len(gotResult))
	}
	if reflect.TypeOf(gotResult) != reflect.TypeOf([]byte{}) {
		t.Fatalf("Want type '[]byte', but got '%T'", reflect.TypeOf(gotResult))
	}
	if gotError != nil {
		t.Fatalf("Want error (nil), but got (%e)", gotError)
	}
}
