package vm

type rtype int64

const (
	untyped = rtype(iota)
	rInt
	rFloat
	rBool
	rAtom
	rString
)
