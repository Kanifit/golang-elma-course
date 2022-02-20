// Package unique_element Чудные вхождения в массив
package unique_element

func Solution(set []int) int {
	xorResult := 0

	for _, number := range set {
		xorResult = xorResult ^ number
	}

	return xorResult
}
