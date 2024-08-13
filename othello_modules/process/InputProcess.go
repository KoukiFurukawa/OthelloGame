package process

import (
	"fmt"
)

func ScanUserInput() (int, int, error) {
	var userInput string
	fmt.Printf("置く場所を入力 : ")
	fmt.Scan(&userInput)
	i, j, err := isValid(userInput)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid input")
	}
	return i, j, err
}

func isValid(input string) (int, int, error) {
	if len(input) != 2 {
		return -1, -1, fmt.Errorf("allowed input lengths is 2")
	}
	allowedFirstString := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	allowedSecondString := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	if Contains(allowedFirstString, string(input[0])) && Contains(allowedSecondString, string(input[1])) {
		return Find(allowedFirstString, string(input[0])), Find(allowedSecondString, string(input[1])), nil
	}
	return -1, -1, fmt.Errorf("invalid input")
}

func Contains(slice []string, key string) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}
	return false
}

func Find(slice []string, key string) int {
	for i, s := range slice {
		if s == key {
			return i
		}
	}
	return -1
}
