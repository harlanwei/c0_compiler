package instruction

const (
	ConstantKindInt = iota
	ConstantKindDouble
	ConstantKindString
)

type Constant struct {
	Kind    int
	Value   interface{}
	Address int
}
