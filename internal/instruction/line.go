package instruction

type Line struct {
	I        Instruction
	Operands *[]int
}

func (l *Line) SetFirstOperandTo(operand int) {
	if len(*l.Operands) < 1 {
		panic("trying to set an operand that does not exist")
	}
	(*l.Operands)[0] = operand
}
