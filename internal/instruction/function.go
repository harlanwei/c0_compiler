package instruction

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/common"
	"c0_compiler/internal/token"
	"container/heap"
)

type Error = cc0_error.Error
type PriorityQueue = common.PriorityQueue

type FnInstructions struct {
	lines  *[]Line
	offset int
}

type Fn struct {
	instructions          *FnInstructions
	ReturnType            int
	Parameters            *[]string
	emptyMemorySlots      *PriorityQueue
	currentConstantOffset int
	stackSize             int
	RelatedSymbolTable    *SymbolTable
}

func InitFn(returnType int) (res *Fn) {
	res = &Fn{
		instructions:          &FnInstructions{lines: &[]Line{}, offset: 0},
		Parameters:            &[]string{},
		ReturnType:            returnType,
		currentConstantOffset: 0,
		emptyMemorySlots:      &PriorityQueue{0},
	}
	heap.Init(res.emptyMemorySlots)
	return
}

func (f *Fn) GetLines() *[]Line {
	return f.instructions.lines
}

func (f *Fn) GetCurrentOffset() int {
	return len(*f.instructions.lines)
}

func (f *Fn) ReplaceNopAt(offset int, instruction int, operands ...int) {
	f.ChangeInstructionTo(offset, instruction, operands...)
}

func (f *Fn) GetCurrentConstantOffset() int {
	return f.currentConstantOffset
}

func (f *Fn) GetAnEmptyConstantSlot() (result int) {
	result = f.currentConstantOffset
	f.currentConstantOffset++
	return
}

func (f *Fn) GetPreviousLine() *Line {
	lines := *f.instructions.lines
	return &(lines[len(lines)-2])
}

func (f *Fn) GetCurrentLine() *Line {
	lines := *f.instructions.lines
	return &(lines[len(lines)-1])
}

func (f *Fn) generateLine(instruction int, operands ...int) Line {
	i := GetInstruction(instruction)
	f.stackSize += i.changesToStackSize
	if !i.IsValidInstruction(operands...) {
		cc0_error.PrintfToStdErr("Incorrect usage of instruction 0x%x!\n", instruction)
		cc0_error.ThrowAndExit(cc0_error.Analyzer)
	}
	f.instructions.offset += i.offset
	return Line{I: i, Operands: &operands}
}

func (f *Fn) Append(instruction int, operands ...int) {
	*f.instructions.lines = append(*f.instructions.lines, f.generateLine(instruction, operands...))
}

func (f *Fn) NextMemorySlot(kind int) (slot int) {
	queue := f.emptyMemorySlots
	slot = queue.Pop().(int)
	if queue.Len() == 0 {
		if kind == token.Double {
			queue.Push(slot + 2)
		} else {
			queue.Push(slot + 1)
		}
	} else {
		if kind != token.Double {
			return
		}
		slotBackup := []int{}
		for {
			nextSlot := queue.Pop().(int)
			if nextSlot-slot >= 2 {
				for _, backedSlot := range slotBackup {
					queue.Push(backedSlot)
				}
				return
			}
			slotBackup = append(slotBackup, slot)
			slot = nextSlot
		}
	}
	return
}

func (f *Fn) PopStack(reservedSize int) {
	f.Append(Popn, f.stackSize-reservedSize)
}

func (f *Fn) ChangeInstructionTo(offset, instruction int, operands ...int) {
	l := &((*f.instructions.lines)[offset])
	copied := make([]int, len(operands))
	copy(copied, operands)
	l.I = GetInstruction(instruction)
	l.Operands = &copied
}
