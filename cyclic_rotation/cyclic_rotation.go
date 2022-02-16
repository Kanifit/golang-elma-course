package main

import (
	"fmt"
)

func main() {
	A := []int{3, 8, 9, 7, 6}
	K := 3

	if K > len(A) {
		K = K - len(A)*(K%len(A))
	}

	fmt.Println(Solution(A, K))
}

func Solution(A []int, K int) []int {
	return append(A[len(A)-K:], A[0:len(A)-K]...)
}
