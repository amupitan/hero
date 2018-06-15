package core

type packageInfo struct {
	name string
}

type imports struct{}

type Runtime struct {
	Imports  imports
	Packages []packageInfo
	Body     Statement
}
