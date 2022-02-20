// Package cyclic_rotation Циклическая ротация
package cyclic_rotation

func Solution(set []int, shift int) []int {

	if shift > len(set) {
		shift = shift % len(set)
	}

	return append(set[len(set)-shift:], set[0:len(set)-shift]...)
}
