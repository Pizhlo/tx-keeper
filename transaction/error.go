package transaction

// TxError represents a transaction-specific error with a custom message.
type TxError struct {
	msg string
}

// Error returns the formatted error message with the package prefix.
func (e *TxError) Error() string {
	return "tx-keeper: " + e.msg
}

var (
	// ErrCannotDoCommit is returned when attempting to commit a transaction  but the rollback function has not been set.
	ErrCannotDoCommit = &TxError{msg: "cannot do commit. Rollback function is not set"}

	// ErrCannotDoRollback is returned when attempting to rollback a transaction but the rollback function has not been set.
	ErrCannotDoRollback = &TxError{msg: "cannot do rollback. Rollback function is not set"}
)
