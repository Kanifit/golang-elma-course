// Package main Проверка последовательности
package main

import (
	"fmt"
	"sort"
)

func main() {
	A := []int{4, 1, 3, 2}

	fmt.Println(Solution(A))
}

func Solution(A []int) int {
	sort.Ints(A)

	for i, number := range A {
		if i+1 != number {
			return 0
		}
	}

	return 1
}
