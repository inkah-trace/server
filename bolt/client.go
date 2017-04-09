package bolt

import (
	"time"

	"github.com/boltdb/bolt"
	"log"
	"github.com/inkah-trace/server"
)

// Client represents a client to the underlying BoltDB data store.
type Client struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	// Services
	traceService TraceService

	db *bolt.DB
}

func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.traceService = TraceService{}
	c.traceService.client = c
	return c
}

// Open opens and initializes the BoltDB database.
func (c *Client) Open() error {
	// Open database file.
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("Error opening DB: %s", err)
		return err
	}
	c.db = db

	// Initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("Traces")); err != nil {
		return err
	}

	if _, err := tx.CreateBucketIfNotExists([]byte("TraceEvents")); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes then underlying BoltDB database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ProfileService returns the profile service associated with the client.
func (c *Client) TraceService() server.TraceService {
	return &c.traceService
}