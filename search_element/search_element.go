// Package main Поиск отсутствующего элемента
package main

import (
	"fmt"
)

func main() {
	A := []int{4, 3, 2, 5}

	fmt.Println(Solution(A))
}

func Solution(A []int) int {
	//Формула суммы первых n-членов арифметической прогрессии, исходя из чисел в диапазоне [1..(N + 1)]
	progressionSum := (2*1 + 1*(len(A)+1-1)) / 2 * (len(A) + 1)

	elementsSum := 0
	for _, number := range A {
		elementsSum += number
	}

	return progressionSum - elementsSum
}
