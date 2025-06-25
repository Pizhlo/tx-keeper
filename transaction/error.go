package transaction

import "fmt"

type TxError struct {
	msg string
}

func (e *TxError) Error() string {
	return fmt.Sprintf("tx-keeper: %s", e.msg)
}

var (
	ErrCannotDoCommit   = &TxError{msg: "cannot do commit. Rollback function is not set"}
	ErrCannotDoRollback = &TxError{msg: "cannot do rollback. Rollback function is not set"}
)
