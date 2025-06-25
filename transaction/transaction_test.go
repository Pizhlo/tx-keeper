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
		{name: "default", opts: []Option{}, want: &Transaction{checkRollback: true}},
		{name: "with need rollback", opts: []Option{WithNoCheckRollback()}, want: &Transaction{checkRollback: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transaction := NewTransaction(tt.opts...)

			assert.NotNil(t, transaction)
			assert.Equal(t, tt.want.checkRollback, transaction.checkRollback)
			assert.Empty(t, transaction.commit.Fns)
			assert.Empty(t, transaction.rollback.Fns)
		})
	}
}

func TestWithCommit(t *testing.T) {
	t.Parallel()

	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := Func(func(_ context.Context, _ ...any) error {
		return nil
	})

	tx.WithCommit(NewCommit(fn, args...))

	assert.Len(t, tx.commit.Fns, 1)
	assert.Equal(t, args, tx.commit.Fns[0].Args)
}

func TestWithRollback(t *testing.T) {
	t.Parallel()

	tx := NewTransaction()

	args := []any{"1", "2", "3"}

	fn := Func(func(_ context.Context, _ ...any) error {
		return nil
	})

	tx.WithRollback(NewRollback(fn, args...))

	assert.Len(t, tx.rollback.Fns, 1)
	assert.Equal(t, args, tx.rollback.Fns[0].Args)
}

func TestDoCommit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		commit      Commit
		opts        []Option
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
			name: "positive case",
			opts: []Option{WithNoCheckRollback()},
			commit: Commit{
				Fns: []Function{
					{
						Fn: Func(func(_ context.Context, args ...any) error {
							for _, arg := range args {
								count(arg.(int)) //nolint:forcetypeassert // we know that arg is int
							}

							return nil
						}),
						Args: []any{1, 2, 3, 4, 5},
					},
				},
			},
			expectedSum: 15,
		},
		{
			name: "negative case: need rollback, but rollback is not set",
			commit: Commit{
				Fns: []Function{
					{
						Fn: Func(func(_ context.Context, _ ...any) error {
							return nil
						}),
					},
				},
			},
			wantError: true,
			err:       ErrCannotDoCommit,
		},
		{
			name: "negative case: fn error",
			opts: []Option{WithNoCheckRollback()},
			commit: Commit{
				Fns: []Function{
					{
						Fn: Func(func(_ context.Context, _ ...any) error {
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
			tx := NewTransaction(tt.opts...)
			setCount(0)

			tx.commit = &tt.commit

			err := tx.DoCommit(t.Context())
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
		opts        []Option
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
			name: "positive case",
			opts: []Option{WithNoCheckRollback()},
			rollback: &Rollback{
				Fns: []Function{
					{
						Fn: Func(func(_ context.Context, args ...any) error {
							for _, arg := range args {
								count(arg.(int)) //nolint:forcetypeassert // we know that arg is int
							}

							return nil
						}),
						Args: []any{1, 2, 3, 4, 5},
					},
				},
			},
			expectedSum: 15,
		},
		{
			name:      "negative case: need rollback, but rollback is not set",
			wantError: true,
			rollback:  &Rollback{},
			err:       ErrCannotDoRollback,
		},
		{
			name: "negative case: fn error",
			rollback: &Rollback{
				Fns: []Function{
					{
						Fn: Func(func(_ context.Context, _ ...any) error {
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
			tx := NewTransaction(tt.opts...)
			setCount(0)

			tx.rollback = tt.rollback

			err := tx.DoRollback(t.Context())
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
