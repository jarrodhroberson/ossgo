package firestore

import (
	"context"
	"time"

	fs "cloud.google.com/go/firestore"
)

// Queries interface defines the methods for creating query instances.
type Queries interface {
	ByKey(key string) *byKeyQuery
	BeforeTime(field string, before time.Time) *timeBeforeQuery
	AfterTime(field string, after time.Time) *timeAfterQuery
	TimeBetween(field string, start time.Time, end time.Time) *timeBetweenQuery
	TimeEquals(field string, value time.Time) *timeEqualsQuery
	StringEquals(field string, value string) *stringEqualsQuery
	IntEquals(field string, value int) *intEqualsQuery
	BoolEquals(field string, value bool) *boolEqualsQuery
}

// Query interface defines the Execute method for all query types.
type Query interface {
	Execute(ctx context.Context) (*fs.DocumentIterator, error)
}

// byKeyQuery represents a query to retrieve a document by key (document ID).
type byKeyQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	key        string            // Firestore document ID
}

// Execute executes the byKeyQuery against Firestore.
func (q *byKeyQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Limit(1).Documents(ctx), nil
}

// timeBeforeQuery represents a query to retrieve documents before a time.
type timeBeforeQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	before     time.Time         // Time before (exclusive)
}

// Execute executes the timeBeforeQuery against Firestore.
func (q *timeBeforeQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, "<", q.before).Documents(ctx), nil
}

// timeAfterQuery represents a query to retrieve documents after a time.
type timeAfterQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	after      time.Time         // Time after (exclusive)
}

// Execute executes the timeAfterQuery against Firestore.
func (q *timeAfterQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, ">", q.after).Documents(ctx), nil
}

// timeBetweenQuery represents a query to retrieve documents between two times.
type timeBetweenQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	start      time.Time         // Time start (inclusive)
	end        time.Time         // Time end (inclusive)
}

// Execute executes the timeBetweenQuery against Firestore.
func (q *timeBetweenQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, ">=", q.start).Where(q.field, "<=", q.end).Documents(ctx), nil
}

// timeEqualsQuery represents a query to retrieve documents based on a time field.
type timeEqualsQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	value      time.Time         // Value to match
}

// Execute executes the timeEqualsQuery against Firestore.
func (q *timeEqualsQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, "==", q.value).Documents(ctx), nil
}

// stringEqualsQuery represents a query to retrieve documents based on a string field.
type stringEqualsQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	value      string            // Value to match
}

// Execute executes the stringEqualsQuery against Firestore.
func (q *stringEqualsQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, "==", q.value).Documents(ctx), nil
}

// IntEqualsQuery represents a query to retrieve documents based on an int field.
type intEqualsQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	value      int               // Value to match
}

// Execute executes the IntEqualsQuery against Firestore.
func (q *intEqualsQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, "==", q.value).Documents(ctx), nil
}

// BoolEqualsQuery represents a query to retrieve documents based on a bool field.
type boolEqualsQuery struct {
	collection *fs.CollectionRef // Firestore collection name
	field      string            // Field name to query
	value      bool              // Value to match
}

// Execute executes the BoolEqualsQuery against Firestore.
func (q *boolEqualsQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	return q.collection.Where(q.field, "==", q.value).Documents(ctx), nil
}

// newQuery struct to hold the query creation functions.
type newQuery struct {
	collection *fs.CollectionRef // Firestore collection name
}

func (nq newQuery) Execute(ctx context.Context) (*fs.DocumentIterator, error) {
	//TODO implement me
	panic("implement me")
}

// ByKey creates a new ByKeyQuery instance.
func (nq newQuery) ByKey(key string) *byKeyQuery {
	return &byKeyQuery{collection: nq.collection, key: key}
}

// BeforeTime creates a new TimeBeforeQuery.
func (nq newQuery) BeforeTime(field string, before time.Time) *timeBeforeQuery {
	return &timeBeforeQuery{collection: nq.collection, field: field, before: before}
}

// AfterTime creates a new TimeAfterQuery.
func (nq newQuery) AfterTime(field string, after time.Time) *timeAfterQuery {
	return &timeAfterQuery{collection: nq.collection, field: field, after: after}
}

// TimeBetween creates a new TimeBetweenQuery.
func (nq newQuery) TimeBetween(field string, start time.Time, end time.Time) *timeBetweenQuery {
	return &timeBetweenQuery{collection: nq.collection, field: field, start: start, end: end}
}

// TimeEquals creates a new TimeEqualsQuery.
func (nq newQuery) TimeEquals(field string, value time.Time) *timeEqualsQuery {
	return &timeEqualsQuery{collection: nq.collection, field: field, value: value}
}

// StringEquals creates a new stringEqualsQuery for equality checks.
func (nq newQuery) StringEquals(field string, value string) *stringEqualsQuery {
	return &stringEqualsQuery{collection: nq.collection, field: field, value: value}
}

// IntEquals creates a new IntEqualsQuery for equality checks.
func (nq newQuery) IntEquals(field string, value int) *intEqualsQuery {
	return &intEqualsQuery{collection: nq.collection, field: field, value: value}
}

// BoolEquals creates a new BoolEqualsQuery for equality checks.
func (nq newQuery) BoolEquals(field string, value bool) *boolEqualsQuery {
	return &boolEqualsQuery{collection: nq.collection, field: field, value: value}
}
