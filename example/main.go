package main

import (
	"context"
	"fmt"
	"log"

	txkeeper "github.com/Pizhlo/tx-keeper" //nolint:depguard // для примера импортируем пакет
)

//nolint:forbidigo
func main() {
	ctx := context.Background()

	// Create a new transaction
	tx := txkeeper.NewTransaction()

	// Define commit operations
	commit := txkeeper.NewCommit(
		func(_ context.Context, args ...any) error {
			fmt.Printf("Executing commit operation with args: %v\n", args)
			return nil
		},
		"arg1", "arg2",
	)

	// Define rollback operations
	rollback := txkeeper.NewRollback(
		func(_ context.Context, args ...any) error {
			fmt.Printf("Executing rollback operation with args: %v\n", args)
			return nil
		},
		"rollback_arg1", "rollback_arg2",
	)

	// Execute the transaction
	if err := tx.WithCommit(commit).WithRollback(rollback).DoCommit(ctx); err != nil {
		log.Printf("Commit failed: %v\n", err)
		// Execute rollback
		if rollbackErr := tx.DoRollback(ctx); rollbackErr != nil {
			log.Printf("Rollback failed: %v\n", rollbackErr)
		}
	} else {
		fmt.Println("Transaction completed successfully!")
	}
}
