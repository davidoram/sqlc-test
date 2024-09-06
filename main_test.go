package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/davidoram/sqlc-test/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDB(t *testing.T) {
	ctx := context.Background()
	log.Printf("connecting to postgres...")
	postgresUrl := "postgres://postgres:postgres@localhost:5432/sqlc_test?sslmode=disable"
	pool, err := pgxpool.New(ctx, postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Printf("connected to postgres")

	// Insert a new row
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Release()

	// Create a customer and save it
	fred, err := db.New(conn).CreateCustomer(ctx, "Fred v1")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created customer: (%s/%d) %s", fred.ID, fred.Revision, fred.Name)

	// Update the customer a few times
	for i := 0; i < 3; i++ {
		fred, err = db.New(conn).UpdateCustomer(ctx, db.UpdateCustomerParams{
			ID:   fred.ID,
			Name: fmt.Sprintf("Fred v%d", i+2),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Updated customer: (%s/%d) %s", fred.ID, fred.Revision, fred.Name)
	}

	// Show the current revision, should be 4
	customer, err := db.New(conn).GetCustomerByID(ctx, fred.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Re-read customer: (%s/%d) %s", customer.ID, customer.Revision, customer.Name)

	// Get all revisions of the customer and print them
	customers, err := db.New(conn).GetCustomerRevisions(ctx, fred.ID)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range customers {
		log.Printf("Revision: (%s/%d) %s", c.ID, c.Revision, c.Name)
	}
}
