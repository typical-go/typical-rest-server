package dbtype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON type
type JSON json.RawMessage

// Value is implement of Valuer
func (j JSON) Value() (driver.Value, error) {
	byteArr := []byte(j)
	return driver.Value(byteArr), nil
}

// Scan is implement of scanner
func (j *JSON) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}
	return nil
}

// MarshalJSON to marshal to json formatted
func (j *JSON) MarshalJSON() ([]byte, error) {
	return *j, nil
}

// UnmarshalJSON to unmarshal
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
