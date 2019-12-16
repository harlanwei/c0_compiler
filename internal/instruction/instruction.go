package instruction

const (
	Nop     = 0x0
	Bipush  = 0x1
	Ipush   = 0x2
	Pop     = 0x4
	Pop2    = 0x5
	Popn    = 0x6
	Dup     = 0x7
	Dup2    = 0x8
	Loadc   = 0x9
	Loada   = 0xa
	New     = 0xb
	Snew    = 0xc
	Iload   = 0x10
	Dload   = 0x11
	Aload   = 0x12
	Iaload  = 0x18
	Daload  = 0x19
	Aaload  = 0x1a
	Istore  = 0x20
	Dstore  = 0x21
	Astore  = 0x22
	Iastore = 0x28
	Dastore = 0x29
	Aastore = 0x2a
	Iadd    = 0x30
	Dadd    = 0x31
	Isub    = 0x34
	Dsub    = 0x35
	Imul    = 0x38
	Dmul    = 0x39
	Idiv    = 0x3c
	Ddiv    = 0x3d
	Ineg    = 0x40
	Dneg    = 0x41
	Icmp    = 0x44
	Dcmp    = 0x45
	I2d     = 0x60
	D2i     = 0x61
	I2c     = 0x62
	Jmp     = 0x70
	Je      = 0x71
	Jne     = 0x72
	Jl      = 0x73
	Jge     = 0x74
	Jg      = 0x75
	Jle     = 0x76
	Call    = 0x80
	Ret     = 0x88
	Iret    = 0x89
	Dret    = 0x8a
	Aret    = 0x8b
	Iprint  = 0xa0
	Dprint  = 0xa1
	Cprint  = 0xa2
	Sprint  = 0xa3
	Printl  = 0xaf
	Iscan   = 0xb0
	Dscan   = 0xb1
	Cscan   = 0xb2
)

type Instruction struct {
	Code                 int
	Representation       string
	nOperands            int
	offset               int
	Operands             []int
	changesToStackSize   int
	variableStackChanges bool
}

var Instructions = map[int]Instruction{
	Nop:     {Code: Nop, Representation: "nop", nOperands: 0, offset: 1},
	Bipush:  {Code: Bipush, Representation: "bipush", nOperands: 1, offset: 2, Operands: []int{1}, changesToStackSize: 1},
	Ipush:   {Code: Ipush, Representation: "ipush", nOperands: 1, offset: 5, Operands: []int{4}, changesToStackSize: 1},
	Pop:     {Code: Pop, Representation: "pop", nOperands: 0, offset: 1, changesToStackSize: 1},
	Pop2:    {Code: Pop2, Representation: "pop2", nOperands: 0, offset: 1, changesToStackSize: 2},
	Popn:    {Code: Popn, Representation: "popn", nOperands: 1, offset: 5, Operands: []int{4}, variableStackChanges: true},
	Dup:     {Code: Dup, Representation: "dup", nOperands: 0, offset: 1, changesToStackSize: 1},
	Dup2:    {Code: Dup2, Representation: "dup2", nOperands: 0, offset: 1, changesToStackSize: 2},
	Loadc:   {Code: Loadc, Representation: "loadc", nOperands: 1, offset: 3, Operands: []int{2}, variableStackChanges: true},
	Loada:   {Code: Loada, Representation: "loada", nOperands: 2, offset: 7, Operands: []int{2, 4}, changesToStackSize: 1},
	New:     {Code: New, Representation: "new", nOperands: 0, offset: 1},
	Snew:    {Code: Snew, Representation: "snew", nOperands: 1, offset: 5, Operands: []int{4}, variableStackChanges: true},
	Iload:   {Code: Iload, Representation: "iload", nOperands: 0, offset: 1, changesToStackSize: 1},
	Dload:   {Code: Dload, Representation: "dload", nOperands: 0, offset: 1, changesToStackSize: 2},
	Aload:   {Code: Aload, Representation: "aload", nOperands: 0, offset: 1, changesToStackSize: 1},
	Iaload:  {Code: Iaload, Representation: "iaload", nOperands: 0, offset: 1, changesToStackSize: -1},
	Daload:  {Code: Daload, Representation: "daload", nOperands: 0, offset: 1},
	Aaload:  {Code: Aaload, Representation: "aaload", nOperands: 0, offset: 1, changesToStackSize: -1},
	Istore:  {Code: Istore, Representation: "istore", nOperands: 0, offset: 1, changesToStackSize: -2},
	Dstore:  {Code: Dstore, Representation: "dstore", nOperands: 0, offset: 1, changesToStackSize: -3},
	Astore:  {Code: Astore, Representation: "astore", nOperands: 0, offset: 1, changesToStackSize: -2},
	Iastore: {Code: Iastore, Representation: "iastore", nOperands: 0, offset: 1, changesToStackSize: -3},
	Dastore: {Code: Dastore, Representation: "dastore", nOperands: 0, offset: 1, changesToStackSize: -4},
	Aastore: {Code: Aastore, Representation: "aastore", nOperands: 0, offset: 1, changesToStackSize: -3},
	Iadd:    {Code: Iadd, Representation: "iadd", nOperands: 0, offset: 1, changesToStackSize: -1},
	Dadd:    {Code: Dadd, Representation: "dadd", nOperands: 0, offset: 1, changesToStackSize: -2},
	Isub:    {Code: Isub, Representation: "isub", nOperands: 0, offset: 1, changesToStackSize: -1},
	Dsub:    {Code: Dsub, Representation: "dsub", nOperands: 0, offset: 1, changesToStackSize: -2},
	Imul:    {Code: Imul, Representation: "imul", nOperands: 0, offset: 1, changesToStackSize: -1},
	Dmul:    {Code: Dmul, Representation: "dmul", nOperands: 0, offset: 1, changesToStackSize: -2},
	Idiv:    {Code: Idiv, Representation: "idiv", nOperands: 0, offset: 1, changesToStackSize: -1},
	Ddiv:    {Code: Ddiv, Representation: "ddiv", nOperands: 0, offset: 1, changesToStackSize: -2},
	Ineg:    {Code: Ineg, Representation: "ineg", nOperands: 0, offset: 1},
	Dneg:    {Code: Dneg, Representation: "dneg", nOperands: 0, offset: 1},
	Icmp:    {Code: Icmp, Representation: "icmp", nOperands: 0, offset: 1, changesToStackSize: -1},
	Dcmp:    {Code: Dcmp, Representation: "dcmp", nOperands: 0, offset: 1, changesToStackSize: -3},
	I2d:     {Code: I2d, Representation: "i2d", nOperands: 0, offset: 1, changesToStackSize: 1},
	D2i:     {Code: D2i, Representation: "d2i", nOperands: 0, offset: 1, changesToStackSize: -1},
	I2c:     {Code: I2c, Representation: "i2c", nOperands: 0, offset: 1},
	Jmp:     {Code: Jmp, Representation: "jmp", nOperands: 1, offset: 3, Operands: []int{2}},
	Je:      {Code: Je, Representation: "je", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Jne:     {Code: Jne, Representation: "jne", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Jl:      {Code: Jl, Representation: "jl", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Jge:     {Code: Jge, Representation: "jge", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Jg:      {Code: Jg, Representation: "jg", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Jle:     {Code: Jle, Representation: "jle", nOperands: 1, offset: 3, Operands: []int{2}, changesToStackSize: -1},
	Call:    {Code: Call, Representation: "call", nOperands: 1, offset: 3, Operands: []int{2}},
	Ret:     {Code: Ret, Representation: "ret", nOperands: 0, offset: 1},
	Iret:    {Code: Iret, Representation: "iret", nOperands: 0, offset: 1},
	Dret:    {Code: Dret, Representation: "dret", nOperands: 0, offset: 1},
	Aret:    {Code: Aret, Representation: "aret", nOperands: 0, offset: 1},
	Iprint:  {Code: Iprint, Representation: "iprint", nOperands: 0, offset: 1, changesToStackSize: -1},
	Dprint:  {Code: Dprint, Representation: "dprint", nOperands: 0, offset: 1, changesToStackSize: -2},
	Cprint:  {Code: Cprint, Representation: "cprint", nOperands: 0, offset: 1, changesToStackSize: -1},
	Sprint:  {Code: Sprint, Representation: "sprint", nOperands: 0, offset: 1, changesToStackSize: -1},
	Printl:  {Code: Printl, Representation: "printl", nOperands: 0, offset: 1},
	Iscan:   {Code: Iscan, Representation: "iscan", nOperands: 0, offset: 1, changesToStackSize: 1},
	Dscan:   {Code: Dscan, Representation: "dscan", nOperands: 0, offset: 1, changesToStackSize: 2},
	Cscan:   {Code: Cscan, Representation: "cscan", nOperands: 0, offset: 1, changesToStackSize: 1},
}

func GetInstruction(code int) Instruction {
	return Instructions[code]
}

func GetCodeFrom(representation string) *Instruction {
	for _, i := range Instructions {
		if i.Representation == representation {
			return &i
		}
	}
	return nil
}

func (instruction Instruction) IsValidInstruction(operands ...int) bool {
	return len(operands) == instruction.nOperands
}
