package transaction

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type testCase struct {
		name string
		opts []transactionOption
		want *Transaction
	}

	tests := []testCase{
		{name: "default", opts: []transactionOption{}, want: &Transaction{needRollback: false}},
		{name: "with need rollback", opts: []transactionOption{WithNeedRollback()}, want: &Transaction{needRollback: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := NewTransaction(tt.opts...)

			assert.NotNil(t, transaction)
			assert.Empty(t, transaction.commit.fns)
			assert.Empty(t, transaction.rollback.fns)

			assert.Equal(t, &commit{fns: make([]function, 0)}, transaction.commit)
			assert.Equal(t, &rollback{fns: make([]function, 0)}, transaction.rollback)
		})
	}

}

func TestWithCommit(t *testing.T) {
	//nolint:varnamelen // стандартное название для такой переменной
	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := transactionFunc(func(ctx context.Context, args ...any) error {
		return nil
	})

	tx.withCommit(NewCommit(fn, args...))

	assert.Len(t, tx.commit.fns, 1)
	assert.Equal(t, args, tx.commit.fns[0].args)
}

func TestWithRollback(t *testing.T) {
	//nolint:varnamelen // стандартное название для такой переменной
	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := transactionFunc(func(ctx context.Context, args ...any) error {
		return nil
	})

	tx.withRollback(NewRollback(fn, args...))

	assert.Len(t, tx.rollback.fns, 1)
	assert.Equal(t, args, tx.rollback.fns[0].args)
}

func TestDoCommit(t *testing.T) {
	type testCase struct {
		name        string
		transaction *Transaction
		commit      commit
		expectedSum int
		wantError   bool
		err         error
	}

	sum := 0

	tests := []testCase{
		{
			name:        "positive case",
			transaction: NewTransaction(),
			commit: commit{
				fns: []function{
					{
						fn: transactionFunc(func(ctx context.Context, args ...any) error {
							for _, arg := range args {
								sum += arg.(int) //nolint:gosec // мы уверены, что arg - int, т.к. мы сами добавляем его в args
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
			commit: commit{
				fns: []function{
					{
						fn: transactionFunc(func(ctx context.Context, args ...any) error {
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
			commit: commit{
				fns: []function{
					{
						fn: transactionFunc(func(ctx context.Context, args ...any) error {
							return fmt.Errorf("some error")
						}),
					},
				},
			},
			wantError: true,
			err:       fmt.Errorf("tx-keeper: error commit on func 0: some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.transaction
			sum = 0

			tx.commit = &tt.commit

			err := tx.doCommit(t.Context())
			if tt.wantError {
				require.Error(t, err)
				require.EqualError(t, tt.err, err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSum, sum)
		})
	}
}

func TestDoRollback(t *testing.T) {
	type testCase struct {
		name        string
		rollback    *rollback
		transaction *Transaction
		expectedSum int
		wantError   bool
		err         error
	}

	sum := 0

	tests := []testCase{
		{
			name:        "positive case",
			transaction: NewTransaction(),
			rollback: &rollback{
				fns: []function{
					{
						fn: transactionFunc(func(ctx context.Context, args ...any) error {
							for _, arg := range args {
								sum += arg.(int) //nolint:gosec // мы уверены, что arg - int, т.к. мы сами добавляем его в args
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
			rollback: &rollback{
				fns: []function{
					{
						fn: transactionFunc(func(ctx context.Context, args ...any) error {
							return fmt.Errorf("some error")
						}),
					},
				},
			},
			wantError: true,
			err:       fmt.Errorf("tx-keeper: error rollback on func 0: some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.transaction
			sum = 0

			tx.rollback = tt.rollback

			err := tx.doRollback(t.Context())
			if tt.wantError {
				require.Error(t, err)
				require.EqualError(t, tt.err, err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSum, sum)
		})
	}
}
