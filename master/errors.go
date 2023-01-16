package master

import "errors"

var (
	ErrWrongEndNode = errors.New("wrong end node")
	ErrWrongDelay   = errors.New("wrong delay")
)
