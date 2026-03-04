package mongodb_helpers_service

import (
	"context"
	"encoding/json"
	"fmt"
	errors_constants "go_boilerplate_project/constants/errors"
	mongodb_models "go_boilerplate_project/models/databases/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *service) resolveContext(ctx context.Context, cancelFunc context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx != nil && cancelFunc != nil {
		return ctx, cancelFunc
	}

	return s.Input.Services.Context.GetContext()
}

func (s *service) getDatabase(db string) *mongo.Database {
	if db != "" {
		return s.Input.Client.MongoDBClient.Client.Database(db)
	}
	return s.Input.Client.MongoDBClient.Client.Database(s.Input.Client.MongoDBClient.Database)
}

func (s *service) resolveDatabaseName(db string) string {
	if db != "" {
		return db
	}
	return s.Input.Client.MongoDBClient.Database
}

func toJSON(v any) string {
	if v == nil {
		return "{}"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(b)
}

// ──────────────────────────────────────────────
// Insert
// ──────────────────────────────────────────────

func (s *service) InsertOne(input *mongodb_models.InsertOneInput) (*mongodb_models.InsertOneOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrInsertOneInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.Document == nil {
		return nil, errors_constants.ErrDocumentRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.insertOne(%s)`,
		dbName, input.Collection, toJSON(input.Document))

	s.Input.Logger.Debugw("MongoDB InsertOne",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	result, err := db.Collection(input.Collection).InsertOne(ctx, input.Document)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB InsertOne failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB InsertOne success",
		"database", dbName,
		"collection", input.Collection,
		"insertedID", result.InsertedID,
	)

	return &mongodb_models.InsertOneOutput{
		ID: result.InsertedID,
	}, nil
}

func (s *service) InsertMany(input *mongodb_models.InsertManyInput) (*mongodb_models.InsertManyOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrInsertManyInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if len(input.Documents) == 0 {
		return nil, errors_constants.ErrDocumentsRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.insertMany(%s)`,
		dbName, input.Collection, toJSON(input.Documents))

	s.Input.Logger.Debugw("MongoDB InsertMany",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
		"documentCount", len(input.Documents),
	)

	result, err := db.Collection(input.Collection).InsertMany(ctx, input.Documents)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB InsertMany failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB InsertMany success",
		"database", dbName,
		"collection", input.Collection,
		"insertedCount", len(result.InsertedIDs),
	)

	return &mongodb_models.InsertManyOutput{
		IDs: result.InsertedIDs,
	}, nil
}

// ──────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────

func (s *service) UpdateOne(input *mongodb_models.UpdateOneInput) (*mongodb_models.UpdateOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrUpdateOneInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		return nil, errors_constants.ErrFilterRequired
	}
	if input.Update == nil {
		return nil, errors_constants.ErrUpdateRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.updateOne(%s, %s, {"upsert": %t})`,
		dbName, input.Collection, toJSON(input.Filter), toJSON(input.Update), input.Upsert)

	s.Input.Logger.Debugw("MongoDB UpdateOne",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
		"upsert", input.Upsert,
	)

	opts := options.Update().SetUpsert(input.Upsert)
	result, err := db.Collection(input.Collection).UpdateOne(ctx, input.Filter, input.Update, opts)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB UpdateOne failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB UpdateOne success",
		"database", dbName,
		"collection", input.Collection,
		"matchedCount", result.MatchedCount,
		"modifiedCount", result.ModifiedCount,
		"upsertedCount", result.UpsertedCount,
	)

	return &mongodb_models.UpdateOutput{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}, nil
}

func (s *service) UpdateMany(input *mongodb_models.UpdateManyInput) (*mongodb_models.UpdateOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrUpdateManyInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		return nil, errors_constants.ErrFilterRequired
	}
	if input.Update == nil {
		return nil, errors_constants.ErrUpdateRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.updateMany(%s, %s, {"upsert": %t})`,
		dbName, input.Collection, toJSON(input.Filter), toJSON(input.Update), input.Upsert)

	s.Input.Logger.Debugw("MongoDB UpdateMany",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
		"upsert", input.Upsert,
	)

	opts := options.Update().SetUpsert(input.Upsert)
	result, err := db.Collection(input.Collection).UpdateMany(ctx, input.Filter, input.Update, opts)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB UpdateMany failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB UpdateMany success",
		"database", dbName,
		"collection", input.Collection,
		"matchedCount", result.MatchedCount,
		"modifiedCount", result.ModifiedCount,
		"upsertedCount", result.UpsertedCount,
	)

	return &mongodb_models.UpdateOutput{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}, nil
}

func (s *service) UpdateOneWithID(input *mongodb_models.UpdateOneWithIDInput) (*mongodb_models.UpdateOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrUpdateOneWithIDInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.ID == nil {
		return nil, errors_constants.ErrIDRequired
	}
	if input.Update == nil {
		return nil, errors_constants.ErrUpdateRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.updateOne({"_id": %s}, %s, {"upsert": %t})`,
		dbName, input.Collection, toJSON(input.ID), toJSON(input.Update), input.Upsert)

	s.Input.Logger.Debugw("MongoDB UpdateOneWithID",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
		"id", input.ID,
		"upsert", input.Upsert,
	)

	opts := options.Update().SetUpsert(input.Upsert)
	result, err := db.Collection(input.Collection).UpdateByID(ctx, input.ID, input.Update, opts)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB UpdateOneWithID failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB UpdateOneWithID success",
		"database", dbName,
		"collection", input.Collection,
		"matchedCount", result.MatchedCount,
		"modifiedCount", result.ModifiedCount,
	)

	return &mongodb_models.UpdateOutput{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}, nil
}

// ──────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────

func (s *service) DeleteOne(input *mongodb_models.DeleteOneInput) (*mongodb_models.DeleteOneOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrDeleteOneInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		return nil, errors_constants.ErrFilterRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.deleteOne(%s)`,
		dbName, input.Collection, toJSON(input.Filter))

	s.Input.Logger.Debugw("MongoDB DeleteOne",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	result, err := db.Collection(input.Collection).DeleteOne(ctx, input.Filter)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB DeleteOne failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB DeleteOne success",
		"database", dbName,
		"collection", input.Collection,
		"deletedCount", result.DeletedCount,
	)

	return &mongodb_models.DeleteOneOutput{
		DeletedCount: result.DeletedCount,
	}, nil
}

func (s *service) DeleteMany(input *mongodb_models.DeleteManyInput) (*mongodb_models.DeleteManyOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrDeleteManyInputNil
	}
	if input.Collection == "" {
		return nil, errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		return nil, errors_constants.ErrFilterRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return nil, errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.deleteMany(%s)`,
		dbName, input.Collection, toJSON(input.Filter))

	s.Input.Logger.Debugw("MongoDB DeleteMany",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	result, err := db.Collection(input.Collection).DeleteMany(ctx, input.Filter)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB DeleteMany failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("MongoDB DeleteMany success",
		"database", dbName,
		"collection", input.Collection,
		"deletedCount", result.DeletedCount,
	)

	return &mongodb_models.DeleteManyOutput{
		DeletedCount: result.DeletedCount,
	}, nil
}

// ──────────────────────────────────────────────
// Find
// ──────────────────────────────────────────────

func (s *service) FindOne(input *mongodb_models.FindOneInput) error {
	if input == nil {
		return errors_constants.ErrFindOneInputNil
	}
	if input.Collection == "" {
		return errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		input.Filter = bson.M{}
	}
	if input.Result == nil {
		return errors_constants.ErrResultPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.findOne(%s`,
		dbName, input.Collection, toJSON(input.Filter))
	if input.Projection != nil {
		query += ", " + toJSON(input.Projection)
	}
	query += ")"
	if input.Sort != nil {
		query += ".sort(" + toJSON(input.Sort) + ")"
	}

	s.Input.Logger.Debugw("MongoDB FindOne",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	opts := options.FindOne()
	if input.Projection != nil {
		opts.SetProjection(input.Projection)
	}
	if input.Sort != nil {
		opts.SetSort(input.Sort)
	}

	err := db.Collection(input.Collection).FindOne(ctx, input.Filter, opts).Decode(input.Result)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB FindOne failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MongoDB FindOne success",
		"database", dbName,
		"collection", input.Collection,
	)
	return nil
}

func (s *service) FindMany(input *mongodb_models.FindManyInput) error {
	if input == nil {
		return errors_constants.ErrFindManyInputNil
	}
	if input.Collection == "" {
		return errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		input.Filter = bson.M{}
	}
	if input.Results == nil {
		return errors_constants.ErrResultsPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.find(%s`,
		dbName, input.Collection, toJSON(input.Filter))
	if input.Projection != nil {
		query += ", " + toJSON(input.Projection)
	}
	query += ")"
	if input.Sort != nil {
		query += ".sort(" + toJSON(input.Sort) + ")"
	}

	s.Input.Logger.Debugw("MongoDB FindMany",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	opts := options.Find()
	if input.Sort != nil {
		opts.SetSort(input.Sort)
	}
	if input.Projection != nil {
		opts.SetProjection(input.Projection)
	}

	cursor, err := db.Collection(input.Collection).Find(ctx, input.Filter, opts)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB FindMany failed",
			"query", query,
			"error", err,
		)
		return err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, input.Results)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB FindMany cursor decode failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MongoDB FindMany success",
		"database", dbName,
		"collection", input.Collection,
	)
	return nil
}

func (s *service) FindManyWithFilters(input *mongodb_models.FindManyWithFiltersInput) error {
	if input == nil {
		return errors_constants.ErrFindManyWithFiltersInputNil
	}
	if input.Collection == "" {
		return errors_constants.ErrCollectionRequired
	}
	if input.Filter == nil {
		input.Filter = bson.M{}
	}
	if input.Results == nil {
		return errors_constants.ErrResultsPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.find(%s`,
		dbName, input.Collection, toJSON(input.Filter))
	if input.Projection != nil {
		query += ", " + toJSON(input.Projection)
	}
	query += ")"
	if input.Sort != nil {
		query += ".sort(" + toJSON(input.Sort) + ")"
	}
	if input.Skip != nil {
		query += fmt.Sprintf(".skip(%d)", *input.Skip)
	}
	if input.Limit != nil {
		query += fmt.Sprintf(".limit(%d)", *input.Limit)
	}

	s.Input.Logger.Debugw("MongoDB FindManyWithFilters",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	opts := options.Find()
	if input.Sort != nil {
		opts.SetSort(input.Sort)
	}
	if input.Skip != nil {
		opts.SetSkip(*input.Skip)
	}
	if input.Limit != nil {
		opts.SetLimit(*input.Limit)
	}
	if input.Projection != nil {
		opts.SetProjection(input.Projection)
	}

	cursor, err := db.Collection(input.Collection).Find(ctx, input.Filter, opts)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB FindManyWithFilters failed",
			"query", query,
			"error", err,
		)
		return err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, input.Results)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB FindManyWithFilters cursor decode failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MongoDB FindManyWithFilters success",
		"database", dbName,
		"collection", input.Collection,
	)
	return nil
}

// ──────────────────────────────────────────────
// Aggregate
// ──────────────────────────────────────────────

func (s *service) Aggregate(input *mongodb_models.AggregateInput) error {
	if input == nil {
		return errors_constants.ErrAggregateInputNil
	}
	if input.Collection == "" {
		return errors_constants.ErrCollectionRequired
	}
	if input.Pipeline == nil {
		return errors_constants.ErrPipelineRequired
	}
	if input.Results == nil {
		return errors_constants.ErrResultsPointerRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()
	db := s.getDatabase(input.Database)
	if db == nil {
		return errors_constants.ErrDatabaseRequired
	}

	dbName := s.resolveDatabaseName(input.Database)
	query := fmt.Sprintf(`db.getSiblingDB("%s").%s.aggregate(%s)`,
		dbName, input.Collection, toJSON(input.Pipeline))

	s.Input.Logger.Debugw("MongoDB Aggregate",
		"query", query,
		"database", dbName,
		"collection", input.Collection,
	)

	cursor, err := db.Collection(input.Collection).Aggregate(ctx, input.Pipeline)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB Aggregate failed",
			"query", query,
			"error", err,
		)
		return err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, input.Results)
	if err != nil {
		s.Input.Logger.Errorw("MongoDB Aggregate cursor decode failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("MongoDB Aggregate success",
		"database", dbName,
		"collection", input.Collection,
	)
	return nil
}

// ──────────────────────────────────────────────
// Transaction
// ──────────────────────────────────────────────

// RunTransaction executes the provided callback inside a MongoDB multi-document
// transaction. Every helper method called within the callback should receive
// the context.Context supplied to the callback so that all operations enlist
// in the same session.
//
// If the callback returns an error the transaction is automatically aborted and
// all changes are rolled back.
//
// NOTE: MongoDB transactions require a replica set or sharded cluster; they are
// not supported on standalone servers.
func (s *service) RunTransaction(input *mongodb_models.TransactionInput) error {
	if input == nil {
		return errors_constants.ErrTransactionInputNil
	}
	if input.Callback == nil {
		return errors_constants.ErrTransactionCallbackRequired
	}

	s.Input.Logger.Debugw("MongoDB RunTransaction starting")

	session, err := s.Input.Client.MongoDBClient.Client.StartSession()
	if err != nil {
		s.Input.Logger.Errorw("MongoDB RunTransaction failed to start session", "error", err)
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (any, error) {
		return nil, input.Callback(sessCtx)
	})
	if err != nil {
		s.Input.Logger.Errorw("MongoDB RunTransaction failed (rolled back)", "error", err)
		return err
	}

	s.Input.Logger.Debugw("MongoDB RunTransaction committed successfully")
	return nil
}
