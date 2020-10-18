package reflectkit

import (
	"reflect"
)

// IsZero report if v is zero value
func IsZero(v interface{}) bool {
	val := reflect.ValueOf(v)
	return val.IsValid() && val.IsZero()
}
