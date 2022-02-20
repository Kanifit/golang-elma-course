// Package search_element Поиск отсутствующего элемента
package search_element

//func main() {
//	A := []int{4, 3, 2, 5}
//
//	fmt.Println(Solution(A))
//}

func Solution(set []int) int {
	//Формула суммы первых n-членов арифметической прогрессии, исходя из чисел в диапазоне [1..(N + 1)]
	progressionSum := (2*1 + 1*(len(set)+1-1)) / 2 * (len(set) + 1)

	elementsSum := 0
	for _, number := range set {
		elementsSum += number
	}

	return progressionSum - elementsSum
}
