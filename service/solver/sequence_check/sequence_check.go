// Package sequence_check Проверка последовательности
package sequence_check

import (
	"sort"
)

//Solution решение задачи
func Solution(set []int) int {
	sort.Ints(set)

	for i, number := range set {
		if i+1 != number {
			return 0
		}
	}

	return 1
}
