package main

import (
	"context"
	"fmt"
	"log"

	txkeeper "github.com/Pizhlo/tx-keeper"
)

func main() {
	ctx := context.Background()

	// Create a new transaction
	tx := txkeeper.NewTransaction()

	// Define commit operations
	commit := txkeeper.NewCommit(
		func(ctx context.Context, args ...any) error {
			fmt.Println("Executing commit operation with args:", args)
			return nil
		},
		"arg1", "arg2",
	)

	// Define rollback operations
	rollback := txkeeper.NewRollback(
		func(ctx context.Context, args ...any) error {
			fmt.Println("Executing rollback operation with args:", args)
			return nil
		},
		"rollback_arg1", "rollback_arg2",
	)

	// Execute the transaction
	if err := tx.WithCommit(commit).WithRollback(rollback).DoCommit(ctx); err != nil {
		log.Printf("Commit failed: %v", err)
		// Execute rollback
		if rollbackErr := tx.DoRollback(ctx); rollbackErr != nil {
			log.Printf("Rollback failed: %v", rollbackErr)
		}
	} else {
		fmt.Println("Transaction completed successfully!")
	}
}
