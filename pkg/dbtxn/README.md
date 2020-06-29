# DBTXN

Help database transaction

In `Repository` layer
```go
func (r *RepoImpl) Delete(ctx context.Context) (int64, error) {
    // Use dbtxn for potential transaction usecase
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
    }
    
    // Get the transaction object 
    db := txn.DB

    // result, err := ...

	if err != nil {
        // Set the error when failed
		txn.SetError(err)
		return -1, err
	}

	return result.RowsAffected()
}
```

In `Service` layer
```go
func (s *SvcImpl) SomeOperation(ctx context.Context) error{
    // To begin the transaction and commit or rollback in end function
    defer dbtxn.Begin(&ctx)()

    // ...
}
```