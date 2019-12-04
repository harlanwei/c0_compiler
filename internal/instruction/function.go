package instruction

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/common"
	"container/heap"
	"os"
)

type Error = cc0_error.Error
type PriorityQueue = common.PriorityQueue

type FnInstructions struct {
	lines  *[]Line
	offset int
}

type Fn struct {
	instructions     *FnInstructions
	Parameters       *[]string
	emptyMemorySlots *PriorityQueue
}

func InitFn() (res *Fn) {
	res = &Fn{
		instructions:     &FnInstructions{lines: &[]Line{}, offset: 0},
		Parameters:       &[]string{},
		emptyMemorySlots: &PriorityQueue{0},
	}
	heap.Init(res.emptyMemorySlots)
	return
}

func (f *Fn) GetLines() *[]Line {
	return f.instructions.lines
}

func (f *Fn) GetCurrentOffset() int {
	return f.instructions.offset
}

func (f *Fn) GetPreviousLine() *Line {
	lines := *f.instructions.lines
	return &(lines[len(lines)-2])
}

func (f *Fn) GetCurrentLine() *Line {
	lines := *f.instructions.lines
	return &(lines[len(lines)-1])
}

func (f *Fn) Append(instruction int, operands ...int) {
	i := GetInstruction(instruction)
	if !i.IsValidInstruction(operands...) {
		os.Exit(-1)
	}
	f.instructions.offset += i.offset
	*f.instructions.lines = append(*f.instructions.lines, Line{I: i, Operands: &operands})
}

func (f *Fn) NextMemorySlot() (slot int) {
	queue := f.emptyMemorySlots
	slot = queue.Pop().(int)
	if queue.Len() == 0 {
		queue.Push(slot + 1)
	}
	return
}
