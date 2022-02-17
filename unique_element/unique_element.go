// Package main Чудные вхождения в массив
package main

import (
	"fmt"
)

func main() {
	A := []int{4, 4, 3, 5, 1, 5, 1}

	fmt.Println(Solution(A))
}

func Solution(A []int) int {
	xorResult := 0

	for _, number := range A {
		xorResult = xorResult ^ number
	}

	return xorResult
}
