package types

type builtin struct {
	name string
}

func (b builtin) String() string {
	return b.name
}

type Int struct{ builtin }

func (i Int) IsType(value string) bool {
	return true
}

type Float struct{ builtin }

func (f Float) IsType(value string) bool {
	return true
}

type String struct{ builtin }

func (s String) IsType(value string) bool {
	return true
}

type Rune struct{ builtin }

func (r Rune) IsType(value string) bool {
	return true
}

type Bool struct{ builtin }

func (b Bool) IsType(value string) bool {
	return true
}

type Generic struct{ builtin }

func (g Generic) IsType(value string) bool {
	return true
}
