package raft

import "github.com/rotisserie/eris"

var (
	ErrCmdNotFound      = eris.New("Command not found")
	ErrSyntaxParseCmdOp = eris.New("Parse cmd op syntax error")
)
