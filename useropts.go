package svd

import "io"

type UserOptions struct {
	Out io.Writer
	Dump bool
	Pkg string
}
