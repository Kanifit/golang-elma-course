// Package search_element Поиск отсутствующего элемента
package search_element

//Solution решение задачи
func Solution(set []int) int {
	progressionSum := (len(set) + 1) * (len(set) + 2) / 2

	elementsSum := 0
	for _, number := range set {
		elementsSum += number
	}

	return progressionSum - elementsSum
}
