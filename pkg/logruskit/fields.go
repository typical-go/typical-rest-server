package logruskit

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKey int

const fieldsKey ctxKey = iota

// PutField key value to context
func PutField(c *context.Context, key string, value interface{}) {
	fields := GetFields(*c)
	if fields != nil {
		fields[key] = value
	} else {
		*c = context.WithValue(*c, fieldsKey, logrus.Fields{
			key: value,
		})
	}
}

// GetFields get fields from context
func GetFields(ctx context.Context) logrus.Fields {
	v := ctx.Value(fieldsKey)
	if v == nil {
		return nil
	}
	fields, ok := v.(logrus.Fields)
	if !ok {
		return nil
	}

	return fields
}
