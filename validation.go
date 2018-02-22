package main

import (
	"fmt"
	"strconv"
)

var errType = fmt.Errorf("Unsupported data type. Integer only")
var errMissX = fmt.Errorf("x is not provided")
var errMissY = fmt.Errorf("y is not provided")
var errDivideByZero = fmt.Errorf("Divide by zero")

// validation for add operation
func addValidation(x, y string) (int, int, error) {
	return baseValidation(x, y)
}

// validation for subtract operation
func subValidation(x, y string) (int, int, error) {
	return baseValidation(x, y)
}

// validation for multiply operation
func mulValidation(x, y string) (int, int, error) {
	return baseValidation(x, y)
}

// validation for divide operation
func divValidation(x, y string) (int, int, error) {
	intX, intY, err := baseValidation(x, y)
	if err != nil {
		return 0, 0, err
	}
	if intY == 0 {
		return intX, intY, errDivideByZero
	}
	return intX, intY, nil
}

// baseValidation will validate if both x and y exist
// also they both have int type
func baseValidation(x, y string) (int, int, error) {
	err := qsValidation(x, y)
	if err != nil {
		return 0, 0, err
	}
	intX, err := stringToInt(x)
	if err != nil {
		return 0, 0, errType
	}
	intY, err := stringToInt(y)
	if err != nil {
		return 0, 0, errType
	}
	return intX, intY, nil
}

// qsValidation checks if both x and y are provided
func qsValidation(x, y string) error {
	if x == "" {
		return errMissX
	}
	if y == "" {
		return errMissY
	}
	return nil
}

func stringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
