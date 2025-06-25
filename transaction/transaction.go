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

	checkRollback bool // whether to check for rollback function presence during commit. By default, it is true.
}

// Func is a function type that takes a context and variable arguments and returns an error.
type Func func(ctx context.Context, args ...any) error

// Function represents a Function with its arguments that will be executed during commit or rollback.
type Function struct {
	Fn   Func
	Args []any
}

// Commit contains a slice of functions to be executed during Commit.
type Commit struct {
	Fns []Function
}

// Rollback contains a slice of functions to be executed during Rollback.
type Rollback struct {
	Fns []Function
}

// Option is a function type for configuring Transaction options.
type Option func(*Transaction)

// WithNoCheckRollback returns an option that sets the checkRollback flag to false.
func WithNoCheckRollback() Option {
	return func(t *Transaction) {
		t.checkRollback = false
	}
}

// NewTransaction creates a new Transaction with the given options.
// By default, checkRollback is set to true.
func NewTransaction(opts ...Option) *Transaction {
	t := &Transaction{
		checkRollback: true,
	}

	for _, opt := range opts {
		opt(t)
	}

	t.commit = &Commit{Fns: make([]Function, 0)}
	t.rollback = &Rollback{Fns: make([]Function, 0)}

	return t
}

// NewCommit creates a new commit with a single function and its arguments.
func NewCommit(fn Func, args ...any) *Commit {
	return &Commit{
		Fns: []Function{{fn, args}},
	}
}

// NewRollback creates a new rollback with a single function and its arguments.
func NewRollback(fn Func, args ...any) *Rollback {
	return &Rollback{
		Fns: []Function{{fn, args}},
	}
}

// WithCommit adds a commit function and its arguments to the transaction.
func (t *Transaction) WithCommit(commit *Commit) *Transaction {
	t.commit = commit

	return t
}

// WithRollback adds a rollback function and its arguments to the transaction.
func (t *Transaction) WithRollback(rollback *Rollback) *Transaction {
	t.rollback = rollback

	return t
}

// DoCommit executes all commit functions. If needRollback is true and no rollback function is set,
// it returns an error.
func (t *Transaction) DoCommit(ctx context.Context) error {
	if t.checkRollback && len(t.rollback.Fns) == 0 {
		return ErrCannotDoCommit
	}

	return doCommit(ctx, t.commit)
}

// DoRollback executes all rollback functions. If no rollback function is set, it returns an error.
func (t *Transaction) DoRollback(ctx context.Context) error {
	if len(t.rollback.Fns) == 0 {
		return ErrCannotDoRollback
	}

	return doRollback(ctx, t.rollback)
}

// doCommit executes all functions in the commit slice and returns the first error encountered.
func doCommit(ctx context.Context, commit *Commit) error {
	for i, fn := range commit.Fns {
		if err := fn.Fn(ctx, fn.Args...); err != nil {
			return fmt.Errorf("tx-keeper: error commit on func %d: %w", i, err)
		}
	}

	return nil
}

// doRollback executes all functions in the rollback slice and returns the first error encountered.
func doRollback(ctx context.Context, rollback *Rollback) error {
	for i, fn := range rollback.Fns {
		if err := fn.Fn(ctx, fn.Args...); err != nil {
			return fmt.Errorf("tx-keeper: error rollback on func %d: %w", i, err)
		}
	}

	return nil
}
