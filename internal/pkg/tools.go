package pkg

import (
	"regexp"
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

//func ValidPassword(password string) bool {
//	containNums, _ := regexp.Match(`[0-9]`, []byte(password))
//	containUpper, _ := regexp.Match(`[A-Z][a-z]`, []byte(password))
//	containSymbols, _ := regexp.Match(`[!@#$%^&*_-]`, []byte(password))
//
//	if (len(password) > 8 && len(password) < 20) && containNums && containUpper && containSymbols {
//		return true
//	}
//	return false
//}

func ValidLogin(login string) bool {
	containLatin, _ := regexp.Match(`[a-z0-9]`, []byte(login))
	containRussian, _ := regexp.Match(`[а-яА-Я]`, []byte(login))

	if len(login) > 10 || !containLatin || containRussian {
		return false
	}

	return true
}
