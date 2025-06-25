package transaction

import (
	"context"
	"fmt"
)

// Transaction represents a transaction that can be committed or rolled back.
// It contains commit and rollback functions, and a flag to check if rollback is required.
type Transaction struct {
	commit   *Commit   // function to call on commit.
	rollback *Rollback // function to call on rollback.

	needRollback bool // whether to check for rollback function presence during commit.
}

// transactionFunc is a function type that takes a context and variable arguments and returns an error.
type transactionFunc func(ctx context.Context, args ...any) error

// function represents a function with its arguments that will be executed during commit or rollback.
type function struct {
	fn   transactionFunc
	args []any
}

// Commit contains a slice of functions to be executed during Commit.
type Commit struct {
	fns []function
}

// Rollback contains a slice of functions to be executed during Rollback.
type Rollback struct {
	fns []function
}

// Option is a function type for configuring Transaction options.
type Option func(*Transaction)

// WithNeedRollback returns an option that sets the needRollback flag to true.
func WithNeedRollback() Option {
	return func(t *Transaction) {
		t.needRollback = true
	}
}

// NewTransaction creates a new Transaction with the given options.
// By default, needRollback is set to true.
func NewTransaction(opts ...Option) *Transaction {
	t := &Transaction{}

	for _, opt := range opts {
		opt(t)
	}

	t.commit = &Commit{fns: make([]function, 0)}
	t.rollback = &Rollback{fns: make([]function, 0)}

	return t
}

// NewCommit creates a new commit with a single function and its arguments.
func NewCommit(fn transactionFunc, args ...any) *Commit {
	return &Commit{
		fns: []function{{fn, args}},
	}
}

// NewRollback creates a new rollback with a single function and its arguments.
func NewRollback(fn transactionFunc, args ...any) *Rollback {
	return &Rollback{
		fns: []function{{fn, args}},
	}
}

// withCommit adds a commit function and its arguments to the transaction.
func (t *Transaction) withCommit(commit *Commit) *Transaction {
	t.commit = commit

	return t
}

// withRollback adds a rollback function and its arguments to the transaction.
func (t *Transaction) withRollback(rollback *Rollback) *Transaction {
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
func doCommit(ctx context.Context, commit *Commit) error {
	for i, fn := range commit.fns {
		if err := fn.fn(ctx, fn.args...); err != nil {
			return fmt.Errorf("tx-keeper: error commit on func %d: %w", i, err)
		}
	}

	return nil
}

// doRollback executes all functions in the rollback slice and returns the first error encountered.
func doRollback(ctx context.Context, rollback *Rollback) error {
	for i, fn := range rollback.fns {
		if err := fn.fn(ctx, fn.args...); err != nil {
			return fmt.Errorf("tx-keeper: error rollback on func %d: %w", i, err)
		}
	}

	return nil
}
