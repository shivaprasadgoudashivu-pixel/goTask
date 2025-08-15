package utils

import "errors"

var (
	ErrInvalidSym = errors.New("invalid  SYM")
	ErrInvalidQty = errors.New("invalid  Qty")
	ErrInvalidAct = errors.New("invalid  Act")
	ErrInvalidPrc = errors.New("invalid  Prc")
)
