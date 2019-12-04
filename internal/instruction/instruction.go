package instruction

const (
	Nop          = 0x0
	Bipush       = 0x1
	Ipush        = 0x2
	Pop          = 0x4
	Pop2         = 0x5
	Popn         = 0x6
	Dup          = 0x7
	Dup2         = 0x8
	Loadc        = 0x9
	Loada        = 0xa
	New          = 0xb
	Snew         = 0xc
	Iload        = 0x10
	Dload        = 0x11
	Aload        = 0x12
	Iaload       = 0x18
	Daload       = 0x19
	Aaload       = 0x1a
	Istore       = 0x20
	Dstore       = 0x21
	Astore       = 0x22
	Iastore      = 0x28
	Dastore      = 0x29
	Aastore      = 0x2a
	Iadd         = 0x30
	Dadd         = 0x31
	Isub         = 0x34
	Dsub         = 0x35
	Imul         = 0x38
	Dmul         = 0x39
	Idiv         = 0x3c
	Ddiv         = 0x3d
	Ineg         = 0x40
	Dneg         = 0x41
	Icmp         = 0x44
	Dcmp         = 0x45
	I2d          = 0x60
	D2i          = 0x61
	I2c          = 0x62
	Jmp          = 0x70
	Je           = 0x71
	Jne          = 0x72
	Jl           = 0x73
	Jge          = 0x74
	Jg           = 0x75
	Jle          = 0x76
	Call         = 0x80
	Ret          = 0x88
	Iret         = 0x89
	Dret         = 0x8a
	Aret         = 0x8b
	Iprint       = 0xa0
	Dprint       = 0xa1
	Cprint       = 0xa2
	Sprint       = 0xa3
	Printl       = 0xaf
	Iscan        = 0xb0
	Dscan        = 0xb1
	Cscan        = 0xb2
	AnalyzerCall = 0xffff // for call instructions before assembling
)

type Instruction struct {
	code           int
	representation string
	nOperands      int
	offset         int
}

var instructions = map[int]Instruction{
	Nop:          {code: Nop, representation: "nop", nOperands: 0, offset: 1},
	Bipush:       {code: Bipush, representation: "bipush", nOperands: 1, offset: 2},
	Ipush:        {code: Ipush, representation: "ipush", nOperands: 1, offset: 5},
	Pop:          {code: Pop, representation: "pop", nOperands: 0, offset: 1},
	Pop2:         {code: Pop2, representation: "pop2", nOperands: 0, offset: 1},
	Popn:         {code: Popn, representation: "popn", nOperands: 1, offset: 5},
	Dup:          {code: Dup, representation: "dup", nOperands: 0, offset: 1},
	Dup2:         {code: Dup2, representation: "dup2", nOperands: 0, offset: 1},
	Loadc:        {code: Loadc, representation: "loadc", nOperands: 1, offset: 3},
	Loada:        {code: Loada, representation: "loada", nOperands: 2, offset: 7},
	New:          {code: New, representation: "new", nOperands: 0, offset: 1},
	Snew:         {code: Snew, representation: "snew", nOperands: 1, offset: 5},
	Iload:        {code: Iload, representation: "iload", nOperands: 0, offset: 1},
	Dload:        {code: Dload, representation: "dload", nOperands: 0, offset: 1},
	Aload:        {code: Aload, representation: "aload", nOperands: 0, offset: 1},
	Iaload:       {code: Iaload, representation: "iaload", nOperands: 0, offset: 1},
	Daload:       {code: Daload, representation: "daload", nOperands: 0, offset: 1},
	Aaload:       {code: Aaload, representation: "aaload", nOperands: 0, offset: 1},
	Istore:       {code: Istore, representation: "istore", nOperands: 0, offset: 1},
	Dstore:       {code: Dstore, representation: "dstore", nOperands: 0, offset: 1},
	Astore:       {code: Astore, representation: "astore", nOperands: 0, offset: 1},
	Iastore:      {code: Iastore, representation: "iastore", nOperands: 0, offset: 1},
	Dastore:      {code: Dastore, representation: "dastore", nOperands: 0, offset: 1},
	Aastore:      {code: Aastore, representation: "aastore", nOperands: 0, offset: 1},
	Iadd:         {code: Iadd, representation: "iadd", nOperands: 0, offset: 1},
	Dadd:         {code: Dadd, representation: "dadd", nOperands: 0, offset: 1},
	Isub:         {code: Isub, representation: "isub", nOperands: 0, offset: 1},
	Dsub:         {code: Dsub, representation: "dsub", nOperands: 0, offset: 1},
	Imul:         {code: Imul, representation: "imul", nOperands: 0, offset: 1},
	Dmul:         {code: Dmul, representation: "dmul", nOperands: 0, offset: 1},
	Idiv:         {code: Idiv, representation: "idiv", nOperands: 0, offset: 1},
	Ddiv:         {code: Ddiv, representation: "ddiv", nOperands: 0, offset: 1},
	Ineg:         {code: Ineg, representation: "ineg", nOperands: 0, offset: 1},
	Dneg:         {code: Dneg, representation: "dneg", nOperands: 0, offset: 1},
	Icmp:         {code: Icmp, representation: "icmp", nOperands: 0, offset: 1},
	Dcmp:         {code: Dcmp, representation: "dcmp", nOperands: 0, offset: 1},
	I2d:          {code: I2d, representation: "i2d", nOperands: 0, offset: 1},
	D2i:          {code: D2i, representation: "d2i", nOperands: 0, offset: 1},
	I2c:          {code: I2c, representation: "i2c", nOperands: 0, offset: 1},
	Jmp:          {code: Jmp, representation: "jmp", nOperands: 1, offset: 3},
	Je:           {code: Je, representation: "je", nOperands: 1, offset: 3},
	Jne:          {code: Jne, representation: "jne", nOperands: 1, offset: 3},
	Jl:           {code: Jl, representation: "jl", nOperands: 1, offset: 3},
	Jge:          {code: Jge, representation: "jge", nOperands: 1, offset: 3},
	Jg:           {code: Jg, representation: "jg", nOperands: 1, offset: 3},
	Jle:          {code: Jle, representation: "jle", nOperands: 1, offset: 3},
	Call:         {code: Call, representation: "call", nOperands: 1, offset: 3},
	Ret:          {code: Ret, representation: "ret", nOperands: 0, offset: 1},
	Iret:         {code: Iret, representation: "iret", nOperands: 0, offset: 1},
	Dret:         {code: Dret, representation: "dret", nOperands: 0, offset: 1},
	Aret:         {code: Aret, representation: "aret", nOperands: 0, offset: 1},
	Iprint:       {code: Iprint, representation: "iprint", nOperands: 0, offset: 1},
	Dprint:       {code: Dprint, representation: "dprint", nOperands: 0, offset: 1},
	Cprint:       {code: Cprint, representation: "cprint", nOperands: 0, offset: 1},
	Sprint:       {code: Sprint, representation: "sprint", nOperands: 0, offset: 1},
	Printl:       {code: Printl, representation: "printl", nOperands: 0, offset: 1},
	Iscan:        {code: Iscan, representation: "iscan", nOperands: 0, offset: 1},
	Dscan:        {code: Dscan, representation: "dscan", nOperands: 0, offset: 1},
	Cscan:        {code: Cscan, representation: "cscan", nOperands: 0, offset: 1},
	AnalyzerCall: {code: AnalyzerCall, representation: "#", nOperands: 1, offset: 3},
}

func GetInstruction(code int) Instruction {
	return instructions[code]
}

func (instruction Instruction) IsValidInstruction(operands ...int) bool {
	return len(operands) == instruction.nOperands
}
