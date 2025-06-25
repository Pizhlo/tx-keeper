// Package txkeeper provides an abstract transaction wrapper for atomic execution of operations
// with built-in commit and rollback support.
package txkeeper

import "github.com/Pizhlo/tx-keeper/transaction"

// Transaction represents a transaction that can be committed or rolled back.
// It contains commit and rollback functions, and a flag to check if rollback is required.
type Transaction = transaction.Transaction

// Func is a function type that takes a context and variable arguments and returns an error.
type Func = transaction.Func

// Function represents a Function with its arguments that will be executed during commit or rollback.
type Function = transaction.Function

// Commit contains a slice of functions to be executed during Commit.
type Commit = transaction.Commit

// Rollback contains a slice of functions to be executed during Rollback.
type Rollback = transaction.Rollback

// Option is a function type for configuring Transaction options.
type Option = transaction.Option

// WithNoCheckRollback returns an option that sets the checkRollback flag to false.
var WithNoCheckRollback = transaction.WithNoCheckRollback

// NewTransaction creates a new Transaction with the given options.
// By default, checkRollback is set to true.
var NewTransaction = transaction.NewTransaction

// NewCommit creates a new commit with a single function and its arguments.
var NewCommit = transaction.NewCommit

// NewRollback creates a new rollback with a single function and its arguments.
var NewRollback = transaction.NewRollback

// Error types
var (
	// ErrCannotDoCommit is returned when attempting to commit a transaction but the rollback function has not been set.
	ErrCannotDoCommit = transaction.ErrCannotDoCommit

	// ErrCannotDoRollback is returned when attempting to rollback a transaction but the rollback function has not been set.
	ErrCannotDoRollback = transaction.ErrCannotDoRollback
)
