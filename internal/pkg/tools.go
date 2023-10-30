package pkg

import (
	"strconv"
)

func ConvertStringToArray(data string) []int {
	var result []int
	for n := range data {
		// Проверяем переводится ли string в int
		if num, err := strconv.Atoi(string(data[n])); err == nil {
			result = append(result, num)
		}
	}
	return result
}

func ConvertArrayToString(data []int) string {
	var result string
	result += "{"
	for n := range data {
		intEl := strconv.Itoa(data[n])
		result += intEl
		if n != len(data)-1 {
			result += ", "
		}
	}
	result += "}"
	return result
}
