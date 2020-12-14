package dbtxn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// ContextKey to get transaction
const ContextKey key = iota

type (
	key int
	// Context of transaction
	Context struct {
		Tx  Tx
		Err error
	}
	// CommitFn is commit function to close the transaction
	CommitFn func() error
	// UseHandler responsible to handle transaction
	UseHandler struct {
		*Context
		DB sq.BaseRunner
	}
	// Tx is interface for *db.Tx
	Tx interface {
		sq.BaseRunner
		Rollback() error
		Commit() error
	}
)

// Begin transaction
func Begin(parent *context.Context) *Context {
	c := &Context{}
	*parent = context.WithValue(*parent, ContextKey, c)
	return c
}

// Use transaction if possible
func Use(ctx context.Context, db *sql.DB) (*UseHandler, error) {
	c := Find(ctx)
	if ctx == nil {
		return nil, errors.New("dbtxn: missing context.Context")
	}
	if c == nil { // NOTE: not transactional
		return &UseHandler{DB: db}, nil
	}
	if c.Tx == nil {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.Err = fmt.Errorf("dbtxn: %w", err)
			return nil, c.Err
		}
		c.Tx = tx
	}
	return &UseHandler{DB: c.Tx, Context: c}, nil
}

// Find transaction context
func Find(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}
	c, _ := ctx.Value(ContextKey).(*Context)
	return c
}

// Error of transaction
func Error(ctx context.Context) error {
	if c := Find(ctx); c != nil {
		return c.Err
	}
	return nil
}

//
// Context
//

// Commit if no error
func (c *Context) Commit() error {
	if c.Tx == nil {
		return nil
	}
	if c.Err != nil {
		return c.Tx.Rollback()
	}
	return c.Tx.Commit()
}

// SetError to set error to txn context
func (c *Context) SetError(err error) bool {
	if c != nil && err != nil {
		c.Err = err
		return true
	}
	return false
}
