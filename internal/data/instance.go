package data

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/encoding"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/pointer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Instance manages MongoDB collection operations
type Instance[T any] struct {
	ColName        string
	templateObject T
	db             *mongo.Database
	col            *mongo.Collection
}

func NewDBInstance[T any](colName string) *Instance[T] {
	var templateObject T
	typeof := reflect.TypeOf(templateObject)
	if typeof == nil || typeof.Kind() == reflect.Interface {
		return &Instance[T]{ColName: colName}
	}
	switch typeof.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice:
		return &Instance[T]{ColName: colName, templateObject: templateObject}
	default:
		panic("template object must be Pointer, Map, or Slice, got " + typeof.Kind().String())
	}
}

func (m *Instance[T]) ApplyDatabase(database *mongo.Database) *Instance[T] {
	m.db = database
	m.col = database.Collection(m.ColName)
	return m
}

func (m *Instance[T]) CreateIndex(keys bson.D, options *options.IndexOptions) error {
	_, err := m.col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options,
	})
	return err
}

// GetChangeStreamWithOpt Get oplog with change stream option
func (m *Instance[T]) GetChangeStreamWithOpt(dbName string, collectionName string, opts *options.ChangeStreamOptions, cb func(bson.M)) (err error) {

	ctx := context.Background()
	cur := &mongo.ChangeStream{}
	pipelineData := []bson.D{}
	client := m.db.Client()

	defer cur.Close(ctx)

	if dbName != "" && collectionName != "" { // Watching a collection
		fmt.Println("Watching", dbName+"."+collectionName)

		coll := client.Database(dbName).Collection(collectionName)
		cur, err = coll.Watch(ctx, pipelineData, opts)

	} else if dbName != "" { // Watching a database

		fmt.Println("Watching", dbName)
		db := client.Database(dbName)
		cur, err = db.Watch(ctx, pipelineData, opts)

	} else { // Watching all

		fmt.Println("Watching all")
		cur, err = client.Watch(ctx, pipelineData, opts)
	}

	if err != nil {
		return
	}

	// loop forever look change data
	for cur.Next(ctx) {
		data := bson.M{}
		cur.Decode(&data)
		cb(data)
	}

	return
}

func (m *Instance[T]) GetChangeStream(dbName string, collectionName string, cb func(bson.M)) (err error) {

	opts := options.ChangeStream()
	opts.SetFullDocument(options.UpdateLookup)

	return m.GetChangeStreamWithOpt(dbName, collectionName, opts, cb)
}

// ApplyTransaction
//
// @handler: the transaction will be committed when give a non-error
// @isolation: will be default value when given nil attributes
func (m *Instance[T]) ApplyTransaction(handler func(ctx SessionContext) ([]T, error), isolation *Isolation) ([]T, error) {
	// setup Isolation & txn option
	if isolation == nil {
		isolation = &defaultIsolation
	} else {
		if isolation.Read == nil {
			isolation.Read = defaultIsolation.Read
		}
		if isolation.Write == nil {
			isolation.Write = defaultIsolation.Write
		}
	}
	txnOpts := options.Transaction().SetWriteConcern(isolation.Write).SetReadConcern(isolation.Read)

	// start session
	session, err := m.db.Client().StartSession()
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeDBTransactionFailed, "Failed to start session on "+m.db.Name()+" with error: "+err.Error())
	}
	defer session.EndSession(context.TODO())

	var wrapHandler = func(ctx SessionContext) (any, error) {
		val, err := handler(ctx)
		return val, err
	}

	// apply transaction
	result, txnErr := session.WithTransaction(context.TODO(), wrapHandler, txnOpts)
	if txnErr != nil {
		return nil, errors.InternalServer(common.ErrorCodeDBTransactionFailed, "Failed to commit transaction with error: "+txnErr.Error())
	}
	if results, ok := result.([]T); ok {
		return results, nil
	}

	return nil, errors.InternalServer(common.ErrorCodeInternalError, "internal Error")
}

// convertToBson Go object to map (to get / query)
func (m *Instance[T]) convertToBson(ent any) (bson.M, error) {
	if ent == nil {
		return bson.M{}, nil
	}

	sel, err := bson.Marshal(ent)
	if err != nil {
		return nil, err
	}

	obj := bson.M{}
	err = bson.Unmarshal(sel, &obj)

	return obj, err
}

// bsonDToObject convert bson.D to object
func (m *Instance[T]) bsonDToObject(b bson.D) (T, error) {
	var obj T

	if b == nil {
		return obj, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return obj, err
	}

	err = bson.Unmarshal(bytes, &obj)
	return obj, err
}

// convertToObject convert bson to object
func (m *Instance[T]) convertToObject(b bson.M) (T, error) {
	var obj T
	if b == nil {
		return obj, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return obj, err
	}

	err = bson.Unmarshal(bytes, &obj)
	return obj, err
}

// newObject return new object with same type of TemplateObject
func (m *Instance[T]) newObject() *T {
	var result *T
	return result
}

// newList return new object with same type of TemplateObject
func (m *Instance[T]) newList(limit int) []T {
	return make([]T, 0, limit)
}

func (m *Instance[T]) interfaceSlice(slice any) ([]any, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, common.ErrInterfaceNonSlice
	}

	ret := make([]any, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func (m *Instance[T]) parseSingleResult(result *mongo.SingleResult, action string) ([]T, error) {
	var obj T
	err := result.Decode(&obj)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, "DB Error: "+err.Error())
	}

	// put to slice
	var sliceData = make([]T, 0, 1)
	sliceData = append(sliceData, obj)
	return sliceData, nil
}

// ========================================= METHOD =====================================

// Create inserts a single document
func (m *Instance[T]) Create(ctx context.Context, entity any, opts ...*options.InsertOneOptions) ([]T, error) {
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	bsonM, err := encoding.ToBSON(entity)
	if err != nil {
		return nil, errors.BadRequest(common.ErrorCodeInvalidBson, "bson conversion error: "+err.Error())
	}

	if bsonM["created_time"] == nil {
		bsonM["created_time"] = time.Now()
	}
	bsonM["last_updated_time"] = bsonM["created_time"]

	result, err := m.col.InsertOne(ctx, bsonM, opts...)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCreate, "insert error: "+err.Error())
	}

	bsonM["_id"] = result.InsertedID
	data, err := m.convertToObject(bsonM)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "convert error: "+err.Error())
	}

	return []T{data}, nil
}

func (m *Instance[T]) CreateMany(ctx context.Context, entityList any) ([]interface{}, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, "DB error: Collection "+m.ColName+" is not initialized yet")
	}

	list, err := m.interfaceSlice(entityList)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeNonSlice, "DB error: "+err.Error())
	}

	var bsonList []any
	now := time.Now()
	for _, item := range list {
		// convert to bson
		bsonM, err := encoding.ToBSON(item)
		if err != nil {
			return nil, errors.InternalServer(common.ErrorCodeInvalidBson, "DB error: Invalid bson object -  "+err.Error())
		}

		if bsonM["created_time"] == nil {
			bsonM["created_time"] = now
		}
		bsonM["last_updated_time"] = bsonM["created_time"]

		bsonList = append(bsonList, bsonM)
	}

	opt := options.InsertMany()
	opt.Ordered = pointer.GetPointer(false)

	result, err := m.col.InsertMany(ctx, bsonList, opt)
	if err != nil && len(result.InsertedIDs) == 0 {
		return nil, errors.InternalServer(common.ErrorCreate, "DB error: "+err.Error())
	}

	return result.InsertedIDs, nil
}

// Query retrieves documents with pagination and sorting
func (m *Instance[T]) Query(ctx context.Context, query any, offset, limit int64, sortFields *primitive.D) ([]T, error) {
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	opt := &options.FindOptions{Limit: &limit, Skip: &offset}
	if sortFields != nil {
		opt.Sort = sortFields
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
		opt.Limit = &limit
	}

	bsonM, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "bson conversion error: "+err.Error())
	}

	result, err := m.col.Find(ctx, bsonM, opt)
	if err != nil || result.Err() != nil {
		return nil, errors.NotFound(common.ErrorNotFound, fmt.Sprintf("no %s found", m.ColName))
	}
	defer result.Close(ctx)

	list := make([]T, 0, limit)
	if err = result.All(ctx,

		&list); err != nil || len(list) == 0 {
		return nil, errors.NotFound(common.ErrorNotFound, fmt.Sprintf("no %s found", m.ColName))
	}

	return list, nil
}

func (m *Instance[T]) Update(ctx context.Context, query any, updater any,
	opts ...*options.FindOneAndUpdateOptions) ([]T, error) {
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	bUpdater, err := encoding.ToBSON(updater)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, fmt.Sprintf("DB error: %v", err.Error()))
	}

	bQuery, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, fmt.Sprintf("DB error: %v", err.Error()))
	}

	bUpdater["last_updated_time"] = time.Now()
	bUpdater = bson.M{
		"$set": bUpdater,
	}

	// do update
	result := m.col.FindOneAndUpdate(ctx, bQuery, bUpdater, opts...)
	if result.Err() != nil {
		detail := ""
		if result != nil {
			detail = result.Err().Error()
		}
		return nil, errors.NotFound(common.ErrorNotFound, "Not found any matched "+m.ColName+". Error detail: "+detail)
	}

	return m.parseSingleResult(result, "UpdateOne")
}

func (m *Instance[T]) UpdateOne(ctx context.Context, query any, updater any) ([]T, error) {
	return m.Update(ctx, query, updater, options.FindOneAndUpdate().SetReturnDocument(options.After))
}

func (m *Instance[T]) UpdateMany(ctx context.Context, query any, updater any,
	opts ...*options.UpdateOptions) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	bUpdater, err := encoding.ToBSON(updater)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, fmt.Sprintf("DB Error: %s", err.Error()))
	}

	bUpdater["last_updated_time"] = time.Now()

	// convert to bUpdater
	bQuery, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	// do update
	result, err := m.col.UpdateMany(ctx, bQuery, bson.M{
		"$set": bUpdater,
	}, opts...)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorUpdate, "Update error: UpdateManyWithCtx - "+err.Error())
	}

	if result.MatchedCount == 0 {
		return nil, errors.NotFound(common.ErrorNotFound, "Not found any "+m.ColName+".")
	}

	return nil, nil
}

func (m *Instance[T]) Upsert(ctx context.Context, query any, updater any) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	// convert to bson
	bUpdater, err := encoding.ToBSON(updater)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, fmt.Sprintf("DB Error: %s", err.Error()))
	}

	if bUpdater["_id"] != nil {
		delete(bUpdater, "_id")
	}

	bUpdater["last_updated_time"] = time.Now()
	createdTime, ok := bUpdater["created_time"]
	if !ok || createdTime == nil {
		createdTime = bUpdater["last_updated_time"]
	} else {
		delete(bUpdater, "created_time")
	}

	// convert to bson
	bQuery, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	upsertOpt := &options.FindOneAndUpdateOptions{
		ReturnDocument: pointer.GetPointer(options.After),
		Upsert:         pointer.GetPointer(true),
	}

	result := m.col.FindOneAndUpdate(ctx, bQuery, bson.M{
		"$set": bUpdater,
		"$setOnInsert": bson.M{
			"created_time": createdTime,
		},
	}, upsertOpt)
	if result.Err() != nil {
		detail := ""
		if result != nil {
			detail = result.Err().Error()
		}
		return nil, errors.NotFound(common.ErrorNotFound, "Not found any matched "+m.ColName+". Error detail: "+detail)
	}

	return m.parseSingleResult(result, "UpdateOne")
}

func (m *Instance[T]) ReleaseOne(ctx context.Context, query any, replacement any, opts ...*options.FindOneAndReplaceOptions) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	// convert
	bReplacement, err := m.convertToBson(replacement)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorCodeInvalidBson, fmt.Sprintf("DB Error: %s", err.Error()))
	}

	if bReplacement["created_time"] == nil {
		bReplacement["created_time"] = time.Now()
	}
	bReplacement["last_updated_time"] = time.Now()

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: ReplaceOne - Cannot convert object - "+err.Error())
	}

	// do replace
	result := m.col.FindOneAndReplace(context.TODO(), converted, bReplacement, opts...)
	if result.Err() != nil {
		detail := ""
		if result != nil {
			detail = result.Err().Error()
		}
		return nil, errors.NotFound(common.ErrorNotFound, "Not found any matched "+m.ColName+". Error detail: "+detail)
	}

	return m.parseSingleResult(result, "ReplaceOne")
}

func (m *Instance[T]) Delete(ctx context.Context, query any,
	opts ...*options.DeleteOptions) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	// convert query
	converted, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	_, err = m.col.DeleteMany(ctx, converted, opts...)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorDelete, "DB error: "+err.Error())
	}

	return nil, nil
}

func (m *Instance[T]) Count(ctx context.Context, query any, opts ...*options.CountOptions) (int64, error) {
	// check col
	if m.col == nil {
		return 0, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	// convert query
	converted, err := encoding.ToBSON(query)
	if err != nil {
		return 0, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	// if query is empty -> count by EstimatedDocumentCount else count by CountDocuments
	count := int64(0)
	if len(converted) == 0 {
		count, err = m.col.EstimatedDocumentCount(ctx, nil)
	} else {
		count, err = m.col.CountDocuments(ctx, converted, opts...)
	}
	if err != nil {
		return 0, errors.InternalServer(common.ErrorCount, "DB error: "+err.Error())
	}

	return count, nil
}

func (m *Instance[T]) IncreaseOne(ctx context.Context, query any, fieldName string, value int) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	updater := bson.M{
		"$inc": bson.D{
			{Key: fieldName, Value: value},
		},
		"$setOnInsert": bson.M{
			"created_time": time.Now(),
		},
		"$currentDate": bson.M{
			"last_updated_time": bson.M{"$type": "date"},
		},
	}

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: pointer.GetPointer(options.After),
		Upsert:         pointer.GetPointer(true),
	}

	// convert query
	bsonM, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	result := m.col.FindOneAndUpdate(ctx, bsonM, updater, &opt)

	return m.parseSingleResult(result, "Increase "+fieldName+" of")
}

func (m *Instance[T]) Aggregate(ctx context.Context, pipeline any, result any,
	opts ...*options.AggregateOptions) ([]T, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	q, err := m.col.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorAggregate, "DB error: AggregateWithCtx - "+err.Error())
	}
	err = q.All(context.TODO(), result)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorAggregate, "DB error: AggregateWithCtx - "+err.Error())
	}

	return nil, nil
}

// DistinctWithCtx ...
func (m *Instance[T]) Distinct(ctx context.Context, query any, field string,
	opts ...*options.DistinctOptions) ([]interface{}, error) {
	// check col
	if m.col == nil {
		return nil, errors.InternalServer(common.ErrorColNotInit, fmt.Sprintf("collection %s not initialized", m.ColName))
	}

	// convert query
	converted, err := encoding.ToBSON(query)
	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: Cannot convert object - "+err.Error())
	}

	result, err := m.col.Distinct(ctx, field, converted, opts...)

	if err != nil {
		return nil, errors.InternalServer(common.ErrorConversion, "DB error: DistinctWithCtx "+err.Error())
	}

	return result, nil
}
