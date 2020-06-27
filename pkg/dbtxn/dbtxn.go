package dbtxn

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// ContextKey to get transaction
const ContextKey key = iota

type (
	key int
	// Context of transaction
	Context struct {
		tx  *sql.Tx
		err error
	}
	// CommitFn is commit function to close the transaction
	CommitFn func() error
	// Handler responsible to handle transaction
	Handler struct {
		db  sq.BaseRunner
		txn bool
		*Context
	}
)

// Begin transaction
func Begin(parent *context.Context) CommitFn {
	c := &Context{}
	*parent = context.WithValue(*parent, ContextKey, c)

	return func() error {
		if c.tx == nil {
			return nil
		}
		if c.err != nil {
			return c.tx.Rollback()
		}
		return c.tx.Commit()
	}
}

// Use transaction if possible
func Use(ctx context.Context, db *sql.DB) (*Handler, error) {
	c := Retrieve(ctx)

	// NOTE: not transactional
	if c == nil {
		return &Handler{db: db, txn: false}, nil
	}

	if c.tx == nil {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.err = fmt.Errorf("dbtxn: %w", err)
			return nil, c.err
		}
		c.tx = tx
	}

	return &Handler{db: c.tx, txn: true, Context: c}, nil
}

// Retrieve transaction context
func Retrieve(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}

	c, _ := ctx.Value(ContextKey).(*Context)
	return c
}

// Error of transaction
func Error(ctx context.Context) error {
	if c := Retrieve(ctx); c != nil {
		return c.err
	}
	return nil
}

//
// Handler
//

// SetError to set error to txn context
func (t *Handler) SetError(err error) bool {
	if t.Context != nil {
		t.Context.err = err
		return true
	}
	return false
}

// DB return base runner to run the query
func (t *Handler) DB() sq.BaseRunner {
	return t.db
}

// Txn return true when using transaction
func (t *Handler) Txn() bool {
	return t.txn
}
