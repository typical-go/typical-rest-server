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
	// Handler responsible to handle transaction
	Handler struct {
		DB      sq.BaseRunner
		Context *Context
	}
	// Tx is interface for database transaction
	Tx interface {
		sq.BaseRunner
		Rollback() error
		Commit() error
	}
)

// Begin transaction
func Begin(parent *context.Context) CommitFn {
	c := &Context{}
	*parent = context.WithValue(*parent, ContextKey, c)
	return c.Commit
}

// Use transaction if possible
func Use(ctx context.Context, db *sql.DB) (*Handler, error) {
	c := Find(ctx)
	if ctx == nil {
		return nil, errors.New("dbtxn: missing context.Context")
	}
	// NOTE: not transactional
	if c == nil {
		return &Handler{DB: db}, nil
	}
	if c.Tx == nil {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.Err = fmt.Errorf("dbtxn: %w", err)
			return nil, c.Err
		}
		c.Tx = tx
	}
	return &Handler{DB: c.Tx, Context: c}, nil
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

//
// Handler
//

// SetError to set error to txn context
func (t *Handler) SetError(err error) bool {
	if t.Context != nil && err != nil {
		t.Context.Err = err
		return true
	}
	return false
}
