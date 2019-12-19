package instruction

type Line struct {
	I        Instruction
	Operands *[]int
}

func (l *Line) SetFirstOperandTo(operand int) {
	if len(*l.Operands) < 1 {
		copied := []int{operand}
		l.Operands = &copied
		return
	}
	(*l.Operands)[0] = operand
}

func (l *Line) ChangeInstructionTo(instruction int, operands ...int) {
	copied := make([]int, len(operands))
	copy(copied, operands)
	l.I = GetInstruction(instruction)
	l.Operands = &copied
}
