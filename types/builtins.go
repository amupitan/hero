package types

type verify func(value string) bool

type builtin struct {
	name string
	verify
}

func (b builtin) String() string {
	return b.name
}

func (b builtin) IsType(value string) bool {
	return b.verify(value)
}

var Int = &builtin{
	name:   `int`,
	verify: func(value string) bool { return true },
}

var Float = &builtin{
	name:   `float`,
	verify: func(value string) bool { return true },
}

var String = &builtin{
	name:   `string`,
	verify: func(value string) bool { return true },
}

var Rune = &builtin{
	name:   `rune`,
	verify: func(value string) bool { return true },
}

var Bool = &builtin{
	name:   `bool`,
	verify: func(value string) bool { return true },
}

var Generic = &builtin{
	name:   `generic`,
	verify: func(value string) bool { return true },
}
