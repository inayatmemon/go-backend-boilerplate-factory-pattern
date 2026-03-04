package errors_constants

import "errors"

var (
	ErrInvalidEnvironmentMode = errors.New("invalid environment mode")

	// ── MongoDB Input Nil Errors ─────────────────────────────────────

	ErrInsertOneInputNil           = errors.New("insert one input is nil")
	ErrInsertManyInputNil          = errors.New("insert many input is nil")
	ErrUpdateOneInputNil           = errors.New("update one input is nil")
	ErrUpdateManyInputNil          = errors.New("update many input is nil")
	ErrUpdateOneWithIDInputNil     = errors.New("update one with id input is nil")
	ErrUpdateManyWithIDInputNil    = errors.New("update many with id input is nil")
	ErrDeleteOneInputNil           = errors.New("delete one input is nil")
	ErrDeleteManyInputNil          = errors.New("delete many input is nil")
	ErrFindOneInputNil             = errors.New("find one input is nil")
	ErrFindManyInputNil            = errors.New("find many input is nil")
	ErrFindManyWithFiltersInputNil = errors.New("find many with filters input is nil")
	ErrAggregateInputNil           = errors.New("aggregate input is nil")
	ErrTransactionInputNil         = errors.New("transaction input is nil")

	// ── MongoDB Field Validation Errors ──────────────────────────────

	ErrCollectionRequired          = errors.New("collection is required for mongodb operation")
	ErrDocumentRequired            = errors.New("document is required for mongodb operation and it cannot be nil")
	ErrDocumentsRequired           = errors.New("documents slice is required and cannot be empty")
	ErrDatabaseRequired            = errors.New("database is required for mongodb operation")
	ErrFilterRequired              = errors.New("filter is required for mongodb operation")
	ErrUpdateRequired              = errors.New("update is required for mongodb operation")
	ErrIDRequired                  = errors.New("id is required for mongodb operation")
	ErrIDsRequired                 = errors.New("ids slice is required and cannot be empty")
	ErrPipelineRequired            = errors.New("pipeline is required for aggregate operation")
	ErrResultPointerRequired       = errors.New("result pointer is required to decode the document into")
	ErrResultsPointerRequired      = errors.New("results pointer is required to decode the documents into")
	ErrTransactionCallbackRequired = errors.New("transaction callback function is required")

	// ── Redis Input Nil Errors ───────────────────────────────────────

	ErrRedisSetInputNil    = errors.New("redis set input is nil")
	ErrRedisSetNXInputNil  = errors.New("redis set nx input is nil")
	ErrRedisGetInputNil    = errors.New("redis get input is nil")
	ErrRedisUpdateInputNil = errors.New("redis update input is nil")
	ErrRedisDeleteInputNil = errors.New("redis delete input is nil")
	ErrRedisExistsInputNil = errors.New("redis exists input is nil")
	ErrRedisExpireInputNil = errors.New("redis expire input is nil")

	// ── Redis Field Validation Errors ────────────────────────────────

	ErrRedisKeyRequired   = errors.New("key is required for redis operation")
	ErrRedisKeysRequired  = errors.New("keys slice is required and cannot be empty for redis operation")
	ErrRedisValueRequired = errors.New("value is required for redis operation")

	// ── MySQL Input Nil Errors ───────────────────────────────────────

	ErrMySQLCreateInputNil          = errors.New("mysql create input is nil")
	ErrMySQLCreateInBatchesInputNil = errors.New("mysql create in batches input is nil")
	ErrMySQLFindOneInputNil         = errors.New("mysql find one input is nil")
	ErrMySQLFindManyInputNil        = errors.New("mysql find many input is nil")
	ErrMySQLUpdateInputNil          = errors.New("mysql update input is nil")
	ErrMySQLDeleteInputNil          = errors.New("mysql delete input is nil")
	ErrMySQLRawQueryInputNil        = errors.New("mysql raw query input is nil")
	ErrMySQLExecInputNil            = errors.New("mysql exec input is nil")
	ErrMySQLTransactionInputNil     = errors.New("mysql transaction input is nil")
	ErrMySQLAutoMigrateInputNil     = errors.New("mysql auto migrate input is nil")

	// ── MySQL Field Validation Errors ────────────────────────────────

	ErrMySQLValueRequired          = errors.New("value is required for mysql create operation")
	ErrMySQLBatchSizeRequired      = errors.New("batch size must be greater than zero")
	ErrMySQLResultPointerRequired  = errors.New("result pointer is required for mysql find operation")
	ErrMySQLResultsPointerRequired = errors.New("results pointer is required for mysql find operation")
	ErrMySQLModelRequired          = errors.New("model is required for mysql operation (used for table inference)")
	ErrMySQLConditionsRequired     = errors.New("conditions are required to prevent accidental global updates/deletes")
	ErrMySQLValuesRequired         = errors.New("values are required for mysql update operation")
	ErrMySQLSQLRequired            = errors.New("sql statement is required")
	ErrMySQLCallbackRequired       = errors.New("transaction callback function is required")
	ErrMySQLModelsRequired         = errors.New("models are required for mysql auto migrate operation")
)
