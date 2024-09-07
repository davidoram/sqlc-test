package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/davidoram/sqlc-test/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	pool, ctx := mustConnect(t)
	defer pool.Close()

	// Insert a new row
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	beforeCustCount, beforeRevisionCount := rowCounts(t, ctx, conn)

	// Create a customer and save it
	cust, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	afterCustCount, afterRevisionCount := rowCounts(t, ctx, conn)
	assert.Equal(t, beforeCustCount+1, afterCustCount)
	assert.Equal(t, beforeRevisionCount, afterRevisionCount)
}

func TestUpdate(t *testing.T) {
	pool, ctx := mustConnect(t)
	defer pool.Close()

	// Insert a new row
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	// Create a customer and save it
	cust, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	beforeCustCount, beforeRevisionCount := rowCounts(t, ctx, conn)

	// Update the customer
	cust, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{ID: cust.ID, Revision: cust.Revision, Name: "Fred v2"})
	assert.NoError(t, err)
	assert.Equal(t, "Fred v2", cust.Name)
	assert.Equal(t, int32(2), cust.Revision)

	afterCustCount, afterRevisionCount := rowCounts(t, ctx, conn)
	assert.Equal(t, beforeCustCount, afterCustCount)
	assert.Equal(t, beforeRevisionCount+1, afterRevisionCount)

	// Chect that the old revision is still there
	revisions, err := db.New(conn).GetCustomerRevisions(ctx, cust.ID)
	assert.NoError(t, err)
	assert.Len(t, revisions, 2)
	assert.Equal(t, "Fred v1", revisions[0].Name)
	assert.Equal(t, "Fred v2", revisions[1].Name)
	assert.Equal(t, int32(1), revisions[0].Revision)
	assert.Equal(t, int32(2), revisions[1].Revision)
}

func TestUpdateOptimisticLock(t *testing.T) {
	pool, ctx := mustConnect(t)
	defer pool.Close()

	// Insert a new row
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	// Create a customer and save it
	cust, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	beforeCustCount, beforeRevisionCount := rowCounts(t, ctx, conn)

	// Update the customer twice with the same revision, the second should fail
	_, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{ID: cust.ID, Revision: cust.Revision, Name: "Fred A"})
	assert.NoError(t, err)
	_, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{ID: cust.ID, Revision: cust.Revision, Name: "Fred B"})
	assert.ErrorContains(t, err, "no rows in result set")

	// Read the customer and check it got the first update
	cust, err = db.New(conn).GetCustomerByID(ctx, cust.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Fred A", cust.Name)
	assert.Equal(t, int32(2), cust.Revision)

	afterCustCount, afterRevisionCount := rowCounts(t, ctx, conn)
	assert.Equal(t, beforeCustCount, afterCustCount)
	assert.Equal(t, beforeRevisionCount+1, afterRevisionCount)
}

func TestGetCustomerByID(t *testing.T) {
	pool, ctx := mustConnect(t)
	defer pool.Close()

	// Get a connection
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	// Check that if we try to get a customer that doesn't exist we get an error
	_, err = db.New(conn).GetCustomerByID(ctx, uuid.New())
	assert.ErrorContains(t, err, "no rows in result set")

	// Create a customer and save it
	cust, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	cust, err = db.New(conn).GetCustomerByID(ctx, cust.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	// Update the customer
	cust, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{ID: cust.ID, Revision: cust.Revision, Name: "Fred v2"})
	assert.NoError(t, err)

	// Check that GetCustomerByID returns the updated customer
	cust, err = db.New(conn).GetCustomerByID(ctx, cust.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Fred v2", cust.Name)
	assert.Equal(t, int32(2), cust.Revision)
}

func TestManyRevisions(t *testing.T) {
	pool, ctx := mustConnect(t)
	defer pool.Close()

	// Get a connection
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	// Create a customer and save it
	cust, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	assert.NoError(t, err)
	assert.Equal(t, "Fred v1", cust.Name)
	assert.Equal(t, int32(1), cust.Revision)

	// Update the customer 100 times
	for i := 2; i <= 100; i++ {
		cust, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{ID: cust.ID, Revision: cust.Revision, Name: fmt.Sprintf("Fred v%d", i)})
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Fred v%d", i), cust.Name)
		assert.Equal(t, int32(i), cust.Revision)
	}

	// Check that GetCustomerRevisions returns all the revisions
	revisions, err := db.New(conn).GetCustomerRevisions(ctx, cust.ID)
	assert.NoError(t, err)
	assert.Len(t, revisions, 100)
	for i, rev := range revisions {
		assert.Equal(t, fmt.Sprintf("Fred v%d", i+1), rev.Name)
		assert.Equal(t, int32(i+1), rev.Revision)
	}
}

// mustConnect returns a connection to the test database, and a context
func mustConnect(t testing.TB) (*pgxpool.Pool, context.Context) {
	ctx := context.Background()
	log.Printf("connecting to postgres...")
	postgresUrl := "postgres://postgres:postgres@localhost:5432/sqlc_test?sslmode=disable"
	pool, err := pgxpool.New(ctx, postgresUrl)
	require.NoError(t, err)
	return pool, ctx
}

// rowCounts returns the number of rows in the customers and customer_revisions tables
func rowCounts(t *testing.T, ctx context.Context, conn *pgxpool.Conn) (int64, int64) {
	// Get before counts
	beforeCustCount, err := db.New(conn).CountCustomers(ctx)
	assert.NoError(t, err)
	beforeRevisionCount, err := db.New(conn).CountCustomerRevisions(ctx)
	assert.NoError(t, err)
	return beforeCustCount, beforeRevisionCount
}
