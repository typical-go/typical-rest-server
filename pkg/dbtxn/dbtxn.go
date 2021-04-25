package dbtxn

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-go/pkg/errkit"
)

// ContextKey to get transaction
const ContextKey key = iota

type (
	key int
	// Context of transaction
	Context struct {
		TxMap map[*sql.DB]Tx
		Errs  errkit.Errors
	}
	// CommitFn is commit function to close the transaction
	CommitFn func() error
	// UseHandler responsible to handle transaction
	UseHandler struct {
		*Context
		sq.StdSqlCtx
	}
	// Tx is interface for *db.Tx
	Tx interface {
		sq.StdSqlCtx
		Rollback() error
		Commit() error
	}
)

// NewContext return new instance of Context
func NewContext() *Context {
	return &Context{TxMap: make(map[*sql.DB]Tx)}
}

// Begin transaction
func Begin(parent *context.Context) *Context {
	c := NewContext()
	*parent = context.WithValue(*parent, ContextKey, c)
	return c
}

// Use transaction if possible
func Use(ctx context.Context, db *sql.DB) (*UseHandler, error) {
	if ctx == nil {
		return nil, errors.New("dbtxn: missing context.Context")
	}

	c := Find(ctx)
	if c == nil { // NOTE: not transactional
		return &UseHandler{StdSqlCtx: db}, nil
	}

	tx, err := c.Begin(ctx, db)
	if err != nil {
		return nil, err
	}

	return &UseHandler{StdSqlCtx: tx, Context: c}, nil
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
		return c.Errs.Unwrap()
	}
	return nil
}

//
// Context
//

// Begin transaction
func (c *Context) Begin(ctx context.Context, db *sql.DB) (sq.StdSqlCtx, error) {
	tx, ok := c.TxMap[db]
	if ok {
		return tx, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		c.AppendError(err)
		return nil, err
	}
	c.TxMap[db] = tx
	return tx, nil
}

// Commit if no error
func (c *Context) Commit() error {
	var errs errkit.Errors
	if len(c.Errs) > 0 {
		for _, tx := range c.TxMap {
			errs = append(errs, tx.Rollback())
		}
	} else {
		for _, tx := range c.TxMap {
			errs = append(errs, tx.Commit())
		}
	}

	return errs.Unwrap()
}

// AppendError to append error to txn context
func (c *Context) AppendError(err error) bool {
	if c != nil && err != nil {
		c.Errs = append(c.Errs, err)
		return true
	}
	return false
}
