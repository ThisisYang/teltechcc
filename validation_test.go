package main

import (
	"testing"
)

func TestQSValidation(t *testing.T) {
	cases := []struct {
		name, x, y string
		expErr     error
	}{
		{name: "case 1", x: "", y: "", expErr: errMissX},
		{name: "case 2", x: "3", y: "", expErr: errMissY},
		{name: "case 3", x: "3", y: "3", expErr: nil},
	}
	for _, c := range cases {
		gotErr := qsValidation(c.x, c.y)
		if gotErr != c.expErr {
			t.Errorf("error on: %v\ngot err:\n %v \nexp err\n %v \n", c.name, gotErr, c.expErr)
		}
	}
}

func TestBaseValidation(t *testing.T) {
	cases := []struct {
		name, x, y string
		expX, expY int
		expErr     error
	}{
		{name: "case 1", x: "3", y: "4", expX: 3, expY: 4, expErr: nil},
		{name: "case 2", x: "a", y: "4", expX: 0, expY: 0, expErr: errType},
		{name: "case 3", x: "3", y: "a", expX: 0, expY: 0, expErr: errType},
	}
	for _, c := range cases {
		gotX, gotY, gotErr := baseValidation(c.x, c.y)
		if gotX != c.expX {
			t.Errorf("error on: %v\ngot x:\n %v \nexp x\n %v \n", c.name, gotX, c.expX)
		}
		if gotX != c.expX {
			t.Errorf("error on: %v\ngot y:\n %v \nexp y\n %v \n", c.name, gotY, c.expY)
		}

		if gotErr != c.expErr {
			t.Errorf("error on: %v\ngot err:\n %v \nexp err\n %v \n", c.name, gotErr, c.expErr)
		}
	}
}

func TestDivValidation(t *testing.T) {
	cases := []struct {
		name, x, y string
		expX, expY int
		expErr     error
	}{
		{name: "case 1", x: "3", y: "0", expX: 3, expY: 0, expErr: errDivideByZero},
		{name: "case 2", x: "0", y: "4", expX: 0, expY: 4, expErr: nil},
	}
	for _, c := range cases {
		gotX, gotY, gotErr := divValidation(c.x, c.y)
		if gotX != c.expX {
			t.Errorf("error on: %v\ngot x:\n %v \nexp x\n %v \n", c.name, gotX, c.expX)
		}
		if gotX != c.expX {
			t.Errorf("error on: %v\ngot y:\n %v \nexp y\n %v \n", c.name, gotY, c.expY)
		}

		if gotErr != c.expErr {
			t.Errorf("error on: %v\ngot err:\n %v \nexp err\n %v \n", c.name, gotErr, c.expErr)
		}
	}
}
