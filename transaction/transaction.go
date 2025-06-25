package transaction

import (
	"context"
	"fmt"
)

// Transaction represents a transaction that can be committed or rolled back.
// It contains commit and rollback functions, and a flag to check if rollback is required.
type Transaction struct {
	commit   *commit   // function to call on commit.
	rollback *rollback // function to call on rollback.

	needRollback bool // whether to check for rollback function presence during commit.
}

// transactionFunc is a function type that takes a context and variable arguments and returns an error.
type transactionFunc func(ctx context.Context, args ...any) error

// function represents a function with its arguments that will be executed during commit or rollback.
type function struct {
	fn   transactionFunc
	args []any
}

// commit contains a slice of functions to be executed during commit.
type commit struct {
	fns []function
}

// rollback contains a slice of functions to be executed during rollback.
type rollback struct {
	fns []function
}

// transactionOption is a function type for configuring Transaction options.
type transactionOption func(*Transaction)

// WithNeedRollback returns an option that sets the needRollback flag to true.
func WithNeedRollback() transactionOption {
	return func(t *Transaction) {
		t.needRollback = true
	}
}

// NewTransaction creates a new Transaction with the given options.
// By default, needRollback is set to true.
func NewTransaction(opts ...transactionOption) *Transaction {
	t := &Transaction{
		needRollback: true,
	}

	for _, opt := range opts {
		opt(t)
	}

	t.commit = &commit{fns: make([]function, 0)}
	t.rollback = &rollback{fns: make([]function, 0)}

	return t
}

// NewCommit creates a new commit with a single function and its arguments.
func NewCommit(fn transactionFunc, args ...any) *commit {
	return &commit{
		fns: []function{{fn, args}},
	}
}

// NewRollback creates a new rollback with a single function and its arguments.
func NewRollback(fn transactionFunc, args ...any) *rollback {
	return &rollback{
		fns: []function{{fn, args}},
	}
}

// withCommit adds a commit function and its arguments to the transaction.
func (t *Transaction) withCommit(commit *commit) *Transaction {
	t.commit = commit

	return t
}

// withRollback adds a rollback function and its arguments to the transaction.
func (t *Transaction) withRollback(rollback *rollback) *Transaction {
	t.rollback = rollback

	return t
}

// doCommit executes all commit functions. If needRollback is true and no rollback function is set,
// it returns an error.
func (t *Transaction) doCommit(ctx context.Context) error {
	if t.needRollback && t.rollback == nil {
		return ErrCannotDoCommit
	}

	return doCommit(ctx, t.commit)
}

// doRollback executes all rollback functions. If no rollback function is set, it returns an error.
func (t *Transaction) doRollback(ctx context.Context) error {
	if t.rollback == nil {
		return ErrCannotDoRollback
	}

	return doRollback(ctx, t.rollback)
}

// doCommit executes all functions in the commit slice and returns the first error encountered.
func doCommit(ctx context.Context, commit *commit) error {
	for i, fn := range commit.fns {
		if err := fn.fn(ctx, fn.args...); err != nil {
			return fmt.Errorf("tx-keeper: error commit on func %d: %+v", i, err)
		}
	}

	return nil
}

// doRollback executes all functions in the rollback slice and returns the first error encountered.
func doRollback(ctx context.Context, rollback *rollback) error {
	for i, fn := range rollback.fns {
		if err := fn.fn(ctx, fn.args...); err != nil {
			return fmt.Errorf("tx-keeper: error rollback on func %d: %+v", i, err)
		}
	}

	return nil
}
