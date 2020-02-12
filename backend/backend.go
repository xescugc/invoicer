package backend

type Backend int

//go:generate enumer -type=Backend -output=backend_string.go

const (
	FS Backend = iota
)
