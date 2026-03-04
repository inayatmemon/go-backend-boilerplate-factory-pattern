package mysql_models

import (
	"context"

	"gorm.io/gorm"
)

// --- Create ---

type CreateInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Value      any    // pointer to model struct
}

// --- CreateInBatches ---

type CreateInBatchesInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Value      any    // pointer to slice of model structs
	BatchSize  int
}

// --- FindOne ---

type FindOneInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Result     any    // pointer to model struct to decode into
	Conditions any    // WHERE clause — string, struct, or map
	Args       []any  // placeholder values when Conditions is a string
	Order      string // e.g. "created_at DESC"
}

// --- FindMany ---

type FindManyInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Results    any    // pointer to slice of model structs
	Conditions any    // WHERE clause
	Args       []any  // placeholder values
	Order      string // e.g. "created_at DESC"
	Limit      *int
	Offset     *int
}

// --- Update ---

type UpdateInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Model      any    // pointer to model struct (used for table inference)
	Conditions any    // WHERE clause
	Args       []any  // placeholder values
	Values     any    // map or struct with the updated fields
}

type UpdateOutput struct {
	RowsAffected int64
}

// --- Delete ---

type DeleteInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Model      any    // pointer to model struct (used for table inference)
	Conditions any    // WHERE clause — required to prevent accidental global deletes
	Args       []any  // placeholder values
}

type DeleteOutput struct {
	RowsAffected int64
}

// --- Raw Query (SELECT) ---

type RawQueryInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	SQL        string
	Args       []any
	Result     any // pointer to struct / slice to scan into
}

// --- Exec (non-SELECT) ---

type ExecInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	SQL        string
	Args       []any
}

type ExecOutput struct {
	RowsAffected int64
}

// --- Transaction ---

// TransactionInput wraps a callback that receives a GORM transaction handle.
// All operations in the callback should either use this tx directly or pass
// it via the Tx field of other helper input structs. If the callback returns
// an error the transaction is rolled back automatically.
type TransactionInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Callback   func(tx *gorm.DB) error
}

// --- AutoMigrate ---

type AutoMigrateInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string // optional transaction handle; nil = use default client
	Model      any    // pointer to model struct (used for table inference)
}

type AutoMigrateOutput struct {
	RowsAffected int64
}
