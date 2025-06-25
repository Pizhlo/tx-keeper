package transaction

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		opts []Option
		want *Transaction
	}

	tests := []testCase{
		{name: "default", opts: []Option{}, want: &Transaction{needRollback: false}},
		{name: "with need rollback", opts: []Option{WithNeedRollback()}, want: &Transaction{needRollback: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transaction := NewTransaction(tt.opts...)

			assert.NotNil(t, transaction)
			assert.Empty(t, transaction.commit.fns)
			assert.Empty(t, transaction.rollback.fns)
		})
	}
}

func TestWithCommit(t *testing.T) {
	t.Parallel()

	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := transactionFunc(func(_ context.Context, _ ...any) error {
		return nil
	})

	tx.withCommit(NewCommit(fn, args...))

	assert.Len(t, tx.commit.fns, 1)
	assert.Equal(t, args, tx.commit.fns[0].args)
}

func TestWithRollback(t *testing.T) {
	t.Parallel()

	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := transactionFunc(func(_ context.Context, _ ...any) error {
		return nil
	})

	tx.withRollback(NewRollback(fn, args...))

	assert.Len(t, tx.rollback.fns, 1)
	assert.Equal(t, args, tx.rollback.fns[0].args)
}

func TestDoCommit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		transaction *Transaction
		commit      Commit
		expectedSum int
		wantError   bool
		err         error
	}

	sum := 0
	mu := sync.Mutex{}

	count := func(n int) {
		mu.Lock()
		sum += n
		mu.Unlock()
	}

	tests := []testCase{
		{
			name:        "positive case",
			transaction: NewTransaction(),
			commit: Commit{
				fns: []function{
					{
						fn: transactionFunc(func(_ context.Context, args ...any) error {
							for _, arg := range args {
								count(arg.(int)) //nolint:forcetypeassert // we know that arg is int
							}

							return nil
						}),
						args: []any{1, 2, 3, 4, 5},
					},
				},
			},
			expectedSum: 15,
		},
		{
			name:        "negative case: need rollback, but rollback is not set",
			transaction: &Transaction{needRollback: true},
			commit: Commit{
				fns: []function{
					{
						fn: transactionFunc(func(_ context.Context, _ ...any) error {
							return nil
						}),
						args: []any{1, 2, 3, 4, 5},
					},
				},
			},
			wantError: true,
			err:       ErrCannotDoCommit,
		},
		{
			name:        "negative case: fn error",
			transaction: NewTransaction(),
			commit: Commit{
				fns: []function{
					{
						fn: transactionFunc(func(_ context.Context, _ ...any) error {
							return fmt.Errorf("some error")
						}),
					},
				},
			},
			wantError: true,
			err:       fmt.Errorf("tx-keeper: error commit on func 0: some error"),
		},
	}

	getCount := func() int {
		mu.Lock()
		defer mu.Unlock()

		return sum
	}

	setCount := func(n int) {
		mu.Lock()
		sum = n
		mu.Unlock()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tx := tt.transaction
			setCount(0)

			tx.commit = &tt.commit

			err := tx.doCommit(t.Context())
			if tt.wantError {
				require.Error(t, err)
				require.EqualError(t, tt.err, err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSum, getCount())
		})
	}
}

func TestDoRollback(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		rollback    *Rollback
		transaction *Transaction
		expectedSum int
		wantError   bool
		err         error
	}

	sum := 0
	mu := sync.Mutex{}

	count := func(n int) {
		mu.Lock()
		sum += n
		mu.Unlock()
	}

	tests := []testCase{
		{
			name:        "positive case",
			transaction: NewTransaction(),
			rollback: &Rollback{
				fns: []function{
					{
						fn: transactionFunc(func(_ context.Context, args ...any) error {
							for _, arg := range args {
								count(arg.(int)) //nolint:forcetypeassert // we know that arg is int
							}

							return nil
						}),
						args: []any{1, 2, 3, 4, 5},
					},
				},
			},
			expectedSum: 15,
		},
		{
			name:        "negative case: need rollback, but rollback is not set",
			transaction: &Transaction{},
			wantError:   true,
			err:         ErrCannotDoRollback,
		},
		{
			name:        "negative case: fn error",
			transaction: NewTransaction(),
			rollback: &Rollback{
				fns: []function{
					{
						fn: transactionFunc(func(_ context.Context, _ ...any) error {
							return fmt.Errorf("some error")
						}),
					},
				},
			},
			wantError: true,
			err:       fmt.Errorf("tx-keeper: error rollback on func 0: some error"),
		},
	}

	getCount := func() int {
		mu.Lock()
		defer mu.Unlock()

		return sum
	}

	setCount := func(n int) {
		mu.Lock()
		sum = n
		mu.Unlock()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tx := tt.transaction
			setCount(0)

			tx.rollback = tt.rollback

			err := tx.doRollback(t.Context())
			if tt.wantError {
				require.Error(t, err)
				require.EqualError(t, tt.err, err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSum, getCount())
		})
	}
}
