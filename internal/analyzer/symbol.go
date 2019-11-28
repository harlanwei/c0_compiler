package analyzer

import "container/heap"

type symbol struct {
	kind        c0Type
	isConstant  bool
	isAFunction bool
	address     int
	appendix    interface{}
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

func (st symbolTable) addVariable(name string, kind c0Type) *Error {
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

func (st symbolTable) addFunction(name string) *Error {
	// TODO
	return nil
}

func (st symbolTable) removeConstant(name string) *Error {
	if !st.hasDeclared(name) || !st.symbols[name].isConstant {
		return errorOf(Bug)
	}
	emptySlots.Push(st.symbols[name].address)
	delete(st.symbols, name)
	return nil
}

func (st symbolTable) addConstant(name string, kind c0Type) *Error {
	if st.hasDeclared(name) {
		return errorOf(RedeclaredAnIdentifier)
	}
	st.symbols[name] = symbol{
		kind:       kind,
		isConstant: true,
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
