package mysql_helpers_service

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	errors_constants "go_boilerplate_project/constants/errors"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	"go_boilerplate_project/storage/mysql"

	"gorm.io/gorm"
)

func (s *service) resolveContext(ctx context.Context, cancelFunc context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx != nil && cancelFunc != nil {
		return ctx, cancelFunc
	}
	return s.Input.Services.Context.GetContext()
}

func (s *service) getDB(database string) *gorm.DB {
	if database != "" {
		client, err := mysql.GetMySQLClientByDatabase(s.Input.Env, s.Input.Logger, database)
		if err != nil {
			s.Input.Logger.Errorw("Failed to get MySQL client by database", "error", err)
			return nil
		}
		return client
	}
	return s.Input.Client.MySQLClient
}

func modelName(v any) string {
	if v == nil {
		return "<nil>"
	}
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t.Name()
}

func interpolateSQL(sql string, args []any) string {
	if len(args) == 0 {
		return sql
	}
	result := sql
	for _, arg := range args {
		var val string
		switch v := arg.(type) {
		case string:
			val = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "\\'"))
		case nil:
			val = "NULL"
		case bool:
			if v {
				val = "1"
			} else {
				val = "0"
			}
		default:
			val = fmt.Sprintf("%v", v)
		}
		result = strings.Replace(result, "?", val, 1)
	}
	return result
}

func formatConditions(conditions any, args []any) string {
	if conditions == nil {
		return ""
	}
	switch c := conditions.(type) {
	case string:
		return interpolateSQL(c, args)
	default:
		return fmt.Sprintf("%+v", c)
	}
}

func formatValues(values any) string {
	if values == nil {
		return ""
	}
	return fmt.Sprintf("%+v", values)
}

// ──────────────────────────────────────────────
// Create
// ──────────────────────────────────────────────

func (s *service) Create(input *mysql_models.CreateInput) error {
	if input == nil {
		return errors_constants.ErrMySQLCreateInputNil
	}
	if input.Value == nil {
		return errors_constants.ErrMySQLValueRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Value)
	inTx := input.Database != ""
	query := fmt.Sprintf("INSERT INTO %s SET %s", table, formatValues(input.Value))

	s.Input.Logger.Debugw("MySQL Create",
		"query", query,
		"table", table,
		"inTransaction", inTx,
	)

	err := s.getDB(input.Database).WithContext(ctx).Create(input.Value).Error
	if err != nil {
		s.Input.Logger.Errorw("MySQL Create failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL Create success", "table", table)
	return nil
}

func (s *service) CreateInBatches(input *mysql_models.CreateInBatchesInput) error {
	if input == nil {
		return errors_constants.ErrMySQLCreateInBatchesInputNil
	}
	if input.Value == nil {
		return errors_constants.ErrMySQLValueRequired
	}
	if input.BatchSize <= 0 {
		return errors_constants.ErrMySQLBatchSizeRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Value)
	inTx := input.Database != ""
	query := fmt.Sprintf("INSERT INTO %s (batch, size=%d) VALUES %s",
		table, input.BatchSize, formatValues(input.Value))

	s.Input.Logger.Debugw("MySQL CreateInBatches",
		"query", query,
		"table", table,
		"batchSize", input.BatchSize,
		"inTransaction", inTx,
	)

	err := s.getDB(input.Database).WithContext(ctx).CreateInBatches(input.Value, input.BatchSize).Error
	if err != nil {
		s.Input.Logger.Errorw("MySQL CreateInBatches failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL CreateInBatches success", "table", table)
	return nil
}

// ──────────────────────────────────────────────
// Find
// ──────────────────────────────────────────────

func (s *service) FindOne(input *mysql_models.FindOneInput) error {
	if input == nil {
		return errors_constants.ErrMySQLFindOneInputNil
	}
	if input.Result == nil {
		return errors_constants.ErrMySQLResultPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Result)
	inTx := input.Database != ""

	query := fmt.Sprintf("SELECT * FROM %s", table)
	if where := formatConditions(input.Conditions, input.Args); where != "" {
		query += " WHERE " + where
	}
	if input.Order != "" {
		query += " ORDER BY " + input.Order
	}
	query += " LIMIT 1"

	s.Input.Logger.Debugw("MySQL FindOne",
		"query", query,
		"table", table,
		"inTransaction", inTx,
	)

	q := s.getDB(input.Database).WithContext(ctx)
	if input.Conditions != nil {
		q = q.Where(input.Conditions, input.Args...)
	}
	if input.Order != "" {
		q = q.Order(input.Order)
	}

	err := q.First(input.Result).Error
	if err != nil {
		s.Input.Logger.Errorw("MySQL FindOne failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL FindOne success", "table", table)
	return nil
}

func (s *service) FindMany(input *mysql_models.FindManyInput) error {
	if input == nil {
		return errors_constants.ErrMySQLFindManyInputNil
	}
	if input.Results == nil {
		return errors_constants.ErrMySQLResultsPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Results)
	inTx := input.Database != ""

	query := fmt.Sprintf("SELECT * FROM %s", table)
	if where := formatConditions(input.Conditions, input.Args); where != "" {
		query += " WHERE " + where
	}
	if input.Order != "" {
		query += " ORDER BY " + input.Order
	}
	if input.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *input.Limit)
	}
	if input.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *input.Offset)
	}

	s.Input.Logger.Debugw("MySQL FindMany",
		"query", query,
		"table", table,
		"inTransaction", inTx,
	)

	q := s.getDB(input.Database).WithContext(ctx)
	if input.Conditions != nil {
		q = q.Where(input.Conditions, input.Args...)
	}
	if input.Order != "" {
		q = q.Order(input.Order)
	}
	if input.Limit != nil {
		q = q.Limit(*input.Limit)
	}
	if input.Offset != nil {
		q = q.Offset(*input.Offset)
	}

	err := q.Find(input.Results).Error
	if err != nil {
		s.Input.Logger.Errorw("MySQL FindMany failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL FindMany success", "table", table)
	return nil
}

// ──────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────

func (s *service) Update(input *mysql_models.UpdateInput) (*mysql_models.UpdateOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrMySQLUpdateInputNil
	}
	if input.Model == nil {
		return nil, errors_constants.ErrMySQLModelRequired
	}
	if input.Conditions == nil {
		return nil, errors_constants.ErrMySQLConditionsRequired
	}
	if input.Values == nil {
		return nil, errors_constants.ErrMySQLValuesRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Model)
	inTx := input.Database != ""

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table, formatValues(input.Values), formatConditions(input.Conditions, input.Args))

	s.Input.Logger.Debugw("MySQL Update",
		"query", query,
		"table", table,
		"inTransaction", inTx,
	)

	result := s.getDB(input.Database).WithContext(ctx).
		Model(input.Model).
		Where(input.Conditions, input.Args...).
		Updates(input.Values)

	if result.Error != nil {
		s.Input.Logger.Errorw("MySQL Update failed",
			"query", query,
			"error", result.Error,
		)
		return nil, result.Error
	}

	s.Input.Logger.Debugw("MySQL Update success",
		"table", table,
		"rowsAffected", result.RowsAffected,
	)

	return &mysql_models.UpdateOutput{
		RowsAffected: result.RowsAffected,
	}, nil
}

// ──────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────

func (s *service) Delete(input *mysql_models.DeleteInput) (*mysql_models.DeleteOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrMySQLDeleteInputNil
	}
	if input.Model == nil {
		return nil, errors_constants.ErrMySQLModelRequired
	}
	if input.Conditions == nil {
		return nil, errors_constants.ErrMySQLConditionsRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	table := modelName(input.Model)
	inTx := input.Database != ""

	query := fmt.Sprintf("DELETE FROM %s WHERE %s",
		table, formatConditions(input.Conditions, input.Args))

	s.Input.Logger.Debugw("MySQL Delete",
		"query", query,
		"table", table,
		"inTransaction", inTx,
	)

	result := s.getDB(input.Database).WithContext(ctx).
		Where(input.Conditions, input.Args...).
		Delete(input.Model)

	if result.Error != nil {
		s.Input.Logger.Errorw("MySQL Delete failed",
			"query", query,
			"error", result.Error,
		)
		return nil, result.Error
	}

	s.Input.Logger.Debugw("MySQL Delete success",
		"table", table,
		"rowsAffected", result.RowsAffected,
	)

	return &mysql_models.DeleteOutput{
		RowsAffected: result.RowsAffected,
	}, nil
}

// ──────────────────────────────────────────────
// Raw Query (SELECT)
// ──────────────────────────────────────────────

func (s *service) RawQuery(input *mysql_models.RawQueryInput) error {
	if input == nil {
		return errors_constants.ErrMySQLRawQueryInputNil
	}
	if input.SQL == "" {
		return errors_constants.ErrMySQLSQLRequired
	}
	if input.Result == nil {
		return errors_constants.ErrMySQLResultPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	inTx := input.Database != ""
	query := interpolateSQL(input.SQL, input.Args)

	s.Input.Logger.Debugw("MySQL RawQuery",
		"query", query,
		"inTransaction", inTx,
	)

	err := s.getDB(input.Database).WithContext(ctx).Raw(input.SQL, input.Args...).Scan(input.Result).Error
	if err != nil {
		s.Input.Logger.Errorw("MySQL RawQuery failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL RawQuery success", "query", query)
	return nil
}

// ──────────────────────────────────────────────
// Exec (non-SELECT)
// ──────────────────────────────────────────────

func (s *service) Exec(input *mysql_models.ExecInput) (*mysql_models.ExecOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrMySQLExecInputNil
	}
	if input.SQL == "" {
		return nil, errors_constants.ErrMySQLSQLRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	inTx := input.Database != ""
	query := interpolateSQL(input.SQL, input.Args)

	s.Input.Logger.Debugw("MySQL Exec",
		"query", query,
		"inTransaction", inTx,
	)

	result := s.getDB(input.Database).WithContext(ctx).Exec(input.SQL, input.Args...)
	if result.Error != nil {
		s.Input.Logger.Errorw("MySQL Exec failed",
			"query", query,
			"error", result.Error,
		)
		return nil, result.Error
	}

	s.Input.Logger.Debugw("MySQL Exec success",
		"query", query,
		"rowsAffected", result.RowsAffected,
	)

	return &mysql_models.ExecOutput{
		RowsAffected: result.RowsAffected,
	}, nil
}

// ──────────────────────────────────────────────
// Transaction
// ──────────────────────────────────────────────

// RunTransaction wraps GORM's db.Transaction. The callback receives a *gorm.DB
// transaction handle — pass it via the Tx field of other helper input structs
// so every operation joins the same transaction. If the callback returns an
// error the transaction is rolled back automatically.
func (s *service) RunTransaction(input *mysql_models.TransactionInput) error {
	if input == nil {
		return errors_constants.ErrMySQLTransactionInputNil
	}
	if input.Callback == nil {
		return errors_constants.ErrMySQLCallbackRequired
	}

	s.Input.Logger.Debugw("MySQL RunTransaction starting")

	err := s.Input.Client.MySQLClient.Transaction(input.Callback)
	if err != nil {
		s.Input.Logger.Errorw("MySQL RunTransaction failed (rolled back)", "error", err)
		return err
	}

	s.Input.Logger.Debugw("MySQL RunTransaction committed successfully")
	return nil
}

func (s *service) AutoMigrate(input *mysql_models.AutoMigrateInput) error {
	if input == nil {
		return errors_constants.ErrMySQLAutoMigrateInputNil
	}
	if input.Model == nil {
		return errors_constants.ErrMySQLModelsRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	err := s.getDB(input.Database).WithContext(ctx).AutoMigrate(input.Model)
	if err != nil {
		s.Input.Logger.Errorw("MySQL AutoMigrate failed",
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MySQL AutoMigrate success")
	return nil
}
