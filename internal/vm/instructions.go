package vm

var instructions []instruction

type code = int

type any = interface{}

type instruction struct {
	code     code
	operands []any
}

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

var code2string = map[int]string{
	Nop:     "nop",
	Bipush:  "bipush",
	Ipush:   "ipush",
	Pop:     "pop",
	Pop2:    "pop2",
	Popn:    "popn",
	Dup:     "dup",
	Dup2:    "dup2",
	Loadc:   "loadc",
	Loada:   "loada",
	New:     "new",
	Snew:    "snew",
	Iload:   "iload",
	Dload:   "dload",
	Aload:   "aload",
	Iaload:  "iaload",
	Daload:  "daload",
	Aaload:  "aaload",
	Istore:  "istore",
	Dstore:  "dstore",
	Astore:  "astore",
	Iastore: "iastore",
	Dastore: "dastore",
	Aastore: "aastore",
	Iadd:    "iadd",
	Dadd:    "dadd",
	Isub:    "isub",
	Dsub:    "dsub",
	Imul:    "imul",
	Dmul:    "dmul",
	Idiv:    "idiv",
	Ddiv:    "ddiv",
	Ineg:    "ineg",
	Dneg:    "dneg",
	Icmp:    "icmp",
	Dcmp:    "dcmp",
	I2d:     "i2d",
	D2i:     "d2i",
	I2c:     "i2c",
	Jmp:     "jmp",
	Je:      "je",
	Jne:     "jne",
	Jl:      "jl",
	Jge:     "jge",
	Jg:      "jg",
	Jle:     "jle",
	Call:    "call",
	Ret:     "ret",
	Iret:    "iret",
	Dret:    "dret",
	Aret:    "aret",
}

func AddInstruction(code code, operands ...any) {
	instructions = append(instructions, instruction{
		code:     code,
		operands: operands,
	})
}
