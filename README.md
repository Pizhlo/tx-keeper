# 🚀 tx-keeper

A Go library that provides an abstract transaction wrapper for atomic execution of operations with built-in commit and rollback support.

## 📖 Overview

tx-keeper is a lightweight, storage-agnostic transaction management library that helps you implement atomic operations in your Go applications. It provides a clean abstraction over transaction patterns, allowing you to define commit and rollback operations with their respective arguments.

## ✨ Key Features

- 🔄 **Abstract Transaction Wrapper**: Provides a clean abstraction over transaction operations without being tied to any specific storage backend
- ⚡ **Atomic Execution**: Ensures that operations are executed atomically - either all succeed or all fail
- 🔒 **Commit and Rollback Management**: Helps store and manage commit and rollback operations with their arguments
- 🎯 **Storage Agnostic**: Independent of any specific storage implementation, making it flexible for various use cases
- 🛡️ **Error Handling**: Comprehensive error handling with detailed error messages for debugging

## 📦 Installation

### Latest version
```bash
go get github.com/Pizhlo/tx-keeper
```

### Specific version
```bash
go get github.com/Pizhlo/tx-keeper@v1.0.0
```

## 🚀 Usage

### 📝 Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Pizhlo/tx-keeper/transaction"
)

func main() {
    ctx := context.Background()
    
    // Create a new transaction
    tx := transaction.NewTransaction()
    
    // Define commit operations
    commit := transaction.NewCommit(
        func(ctx context.Context, args ...any) error {
            fmt.Println("Executing commit operation with args:", args)
            return nil
        },
        "arg1", "arg2",
    )
    
    // Define rollback operations
    rollback := transaction.NewRollback(
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
    }
}
```

### 🔧 Advanced Usage with Multiple Operations

```go
// Create a transaction with rollback requirement
tx := transaction.NewTransaction()

// Add multiple commit operations
commit := &transaction.Commit{
    Fns: []transaction.Function{
        {Fn: saveToDatabase, Args: []any{"user", userData}},
        {Fn: sendNotification, Args: []any{"email", emailData}},
        {Fn: updateCache, Args: []any{"user_cache", cacheData}},
    },
}

// Add multiple rollback operations
rollback := &transaction.Rollback{
    Fns: []transaction.Function{
        {Fn: deleteFromDatabase, Args: []any{"user", userID}},
        {Fn: cancelNotification, Args: []any{"email", emailID}},
        {Fn: invalidateCache, Args: []any{"user_cache", userID}},
    },
}

// Execute transaction
if err := tx.WithCommit(commit).WithRollback(rollback).DoCommit(ctx); err != nil {
    // Automatic rollback on failure
    tx.DoRollback(ctx)
}
```

### 🔍 Rollback Check Behavior

By default, tx-keeper requires a rollback function to be set before allowing commit operations. This ensures that you always have a way to undo changes if something goes wrong.

```go
// This will fail - no rollback function provided
tx := transaction.NewTransaction()
commit := transaction.NewCommit(
    func(ctx context.Context, args ...any) error {
        fmt.Println("Executing commit operation")
        return nil
    },
)

// This will return ErrCannotDoCommit
if err := tx.WithCommit(commit).DoCommit(ctx); err != nil {
    fmt.Printf("Commit failed: %v\n", err) // Will print: Commit failed: cannot do commit. Rollback function is not set
}
```

### ⚙️ Disabling Rollback Check

You can disable the rollback check requirement using the `WithNoCheckRollback` option:

```go
// Create transaction without rollback check
tx := transaction.NewTransaction(transaction.WithNoCheckRollback())

commit := transaction.NewCommit(
    func(ctx context.Context, args ...any) error {
        fmt.Println("Executing commit operation without rollback")
        return nil
    },
)

// This will succeed even without a rollback function
if err := tx.WithCommit(commit).DoCommit(ctx); err != nil {
    fmt.Printf("Commit failed: %v\n", err)
} else {
    fmt.Println("Commit successful!")
}
```

⚠️ **Note**: Use `WithNoCheckRollback` carefully, as it removes the safety mechanism that ensures you have a rollback strategy in place.

## 📚 API Reference

### 🔧 Transaction

The main transaction struct that manages commit and rollback operations.

```go
type Transaction struct {
    commit         *Commit
    rollback       *Rollback
    checkRollback  bool  // whether to check for rollback function presence during commit.
}
```

### ✅ Commit

Represents a collection of functions to be executed during commit.

```go
type Commit struct {
    Fns []Function
}
```

### 🔄 Rollback

Represents a collection of functions to be executed during rollback.

```go
type Rollback struct {
    Fns []Function
}
```

### ⚙️ Function

Represents a function with its arguments.

```go
type Function struct {
    Fn   Func
    Args []any
}
```

### 🔧 Func

Represents a function type that takes a context and variable arguments and returns an error.

```go
type Func func(ctx context.Context, args ...any) error
```

## 🚨 Error Handling

The library provides specific error types for different failure scenarios:

- ❌ `ErrCannotDoCommit`: Returned when trying to commit without a rollback function (when `needRollback` is true)
- 🔄 `ErrCannotDoRollback`: Returned when trying to rollback without a rollback function


## 🤝 Contributing

1. 🍴 Fork the repository
2. 🌿 Create your feature branch (`git checkout -b feature/amazing-feature`)
3. 💾 Commit your changes (`git commit -m 'Add some amazing feature'`)
4. 🚀 Push to the branch (`git push origin feature/amazing-feature`)
5. 📝 Open a Pull Request

## 🧪 Testing

Run the tests with:

```bash
go test ./...
```

Or use the provided Makefile:

```bash
make test
```

---

<div align="center">

**Made with ❤️ for the Go community**

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>
