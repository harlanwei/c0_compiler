package analyzer

import "container/heap"

type symbol struct {
	kind       c0Type
	isConstant bool
	isCallable bool
	address    int
	appendix   interface{}
}

type functionInfo struct {
	symbols      *symbolTable
	instructions []string
	parameters   []string
}

type symbolTable struct {
	symbols map[string]symbol
	parent  *symbolTable
}

type function struct {
	symbols symbolTable
}

var emptySlots = &PriorityQueue{0}

func nextSlot() (slot int) {
	slot = emptySlots.Pop().(int)
	if emptySlots.Len() == 0 {
		emptySlots.Push(slot + 1)
	}
	return
}

func (st symbolTable) hasDeclared(name string) bool {
	_, ok := st.symbols[name]
	return ok
}

func (st symbolTable) addAVariable(name string, kind c0Type) *Error {
	if st.hasDeclared(name) {
		return errorOf(RedeclaredAnIdentifier)
	}
	// TODO: assign proper address for variables
	st.symbols[name] = symbol{
		kind:       kind,
		isConstant: false,
		address:    -1,
	}
	return nil
}

func (st symbolTable) addAFunction(name string, returnType c0Type) *Error {
	return st.addAConstant(name, returnType, true)
}

func (sb symbol) addAnInstruction(instruction string) {
	sb.appendix.(functionInfo).addAnInstruction(instruction)
}

func (fi functionInfo) addAnInstruction(instruction string) {
	fi.instructions = append(fi.instructions, instruction)
}

func (st symbolTable) getSymbol(symbol string) *symbol {
	currentTable := &st
	for {
		if sb, ok := currentTable.symbols[symbol]; ok {
			return &sb
		}
		if currentTable.parent != nil {
			currentTable = currentTable.parent
		} else {
			return nil
		}
	}
}

func (st symbolTable) getAddressOf(symbol string) int {
	sb := st.getSymbol(symbol)
	if sb == nil {
		return -1
	}
	return sb.address
}

func (st symbolTable) removeConstant(name string) *Error {
	if !st.hasDeclared(name) || !st.symbols[name].isConstant {
		return errorOf(Bug)
	}
	emptySlots.Push(st.symbols[name].address)
	delete(st.symbols, name)
	return nil
}

func (st symbolTable) addAConstant(name string, kind c0Type, isAFunction bool) *Error {
	if st.hasDeclared(name) {
		return errorOf(RedeclaredAnIdentifier)
	}
	st.symbols[name] = symbol{
		kind:       kind,
		isConstant: true,
		isCallable: isAFunction,
		address:    nextSlot(),
	}
	return nil
}

func initGlobalSymbolTable() *symbolTable {
	heap.Init(emptySlots)
	return &symbolTable{map[string]symbol{}, nil}
}

func createChildTableFor(parent *symbolTable) *symbolTable {
	return &symbolTable{map[string]symbol{}, parent}
}
