package vm

var constantIndex = 0

var constants = map[any]int{}

func AddConstant(value any) int {
	if index, ok := constants[value]; ok {
		return index
	}
	index := constantIndex
	constantIndex++
	constants[value] = index
	return index
}
