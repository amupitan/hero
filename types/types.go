package types

type Type interface {
	IsType(value string) bool
	String() string
}
