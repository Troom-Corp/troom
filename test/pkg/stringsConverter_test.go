package test

import (
	"testing"

	"github.com/Troom-Corp/troom/internal/pkg"
)

// Тест который проверяет коректность работы функции ConvertStringToArray
func TestConvertStringToArray(t *testing.T) {
	testData := "{1, 0, 1, 0, 1, 0, 1, 0}"
	expectResult := []int{1, 0, 1, 0, 1, 0, 1, 0}
	gotResult := pkg.ConvertStringToArray(testData)
	for i := range gotResult {
		if gotResult[i] != expectResult[i] {
			t.Fatalf("Want (%v), but got (%v)", expectResult, gotResult)
		}
	}
}

// Тест который проверяет коректность работы функции ConvertArrayToString
func TestConvertArrayToString(t *testing.T) {
	testData := []int{1, 0, 1, 0, 1, 0, 1, 0}
	expectResult := "{1, 0, 1, 0, 1, 0, 1, 0}"
	gotResult := pkg.ConvertArrayToString(testData)
	if gotResult != expectResult {
		t.Fatalf("Want '%s', but got '%s'", expectResult, gotResult)
	}
}
