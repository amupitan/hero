package ast

type String struct {
	value string
	isRaw bool
}

func (s *String) Value() string {
	// TODO(DEV) use isRaw to determine output
	return s.value
}

func (s *String) String() string {
	return s.Value()
}
