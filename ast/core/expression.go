package core

type Expression interface {
	Statement
	// string value for debugging purposes
	String() string
}
