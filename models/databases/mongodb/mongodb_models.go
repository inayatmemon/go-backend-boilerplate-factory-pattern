package mongodb_models

import "context"

// --- Insert ---

type InsertOneInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Document   any
}

type InsertOneOutput struct {
	ID any
}

type InsertManyInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Documents  []any
}

type InsertManyOutput struct {
	IDs []any
}

// --- Update ---

type UpdateOneInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
	Update     any
	Upsert     bool
}

type UpdateOutput struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    any
}

type UpdateManyInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
	Update     any
	Upsert     bool
}

type UpdateOneWithIDInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	ID         any
	Update     any
	Upsert     bool
}

type UpdateManyWithIDInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	IDs        []any
	Update     any
	Upsert     bool
}

// --- Delete ---

type DeleteOneInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
}

type DeleteOneOutput struct {
	DeletedCount int64
}

type DeleteManyInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
}

type DeleteManyOutput struct {
	DeletedCount int64
}

// --- Find ---

type FindOneInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
	Result     any
	Sort       any
	Projection any
}

type FindManyInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
	Results    any
	Sort       any
	Projection any
}

type FindManyWithFiltersInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Filter     any
	Results    any
	Sort       any
	Skip       *int64
	Limit      *int64
	Projection any
}

// --- Aggregate ---

type AggregateInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Database   string
	Collection string
	Pipeline   any
	Results    any
}

// --- Transaction ---

// TransactionInput holds the callback to execute within a MongoDB transaction.
// The callback receives a context.Context that carries the session — pass it
// through to every helper method's input.Context so all operations join the
// same transaction. If any operation returns an error the entire transaction
// is rolled back automatically.
type TransactionInput struct {
	Callback func(ctx context.Context) error
}
