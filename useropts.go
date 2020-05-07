package svd

import "io"

type UserOptions struct {
	Out           io.Writer
	Dump          bool
	Pkg           string
	InputFilename string
	Tags          string
	Import 	string
}
