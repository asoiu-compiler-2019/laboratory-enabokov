package main

import "fmt"
import "testing"


def nameFunction(a, b) {
	# this is just a comment
	var a float
	var b uint64
	var c double

	testing.Compare(a, b)

	a = 5 * 4 + 3
	fmt.Println(a)

	c = 12 * 23 + 1234 - (12 * 23)
	fmt.Println(b)

	b = 41 + 1 / (12 * 23)

	if b > 0 {
		b += 2
	}

	c := 123.123
	testing.Compare(a, c, b)
}
