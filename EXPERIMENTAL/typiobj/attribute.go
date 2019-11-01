package typiobj

import "reflect"

// Name of object. Return value from name field if available or return its type.
func Name(obj interface{}) (name string) {
	val := reflect.Indirect(reflect.ValueOf(obj)).FieldByName("Name")
	name = val.String()
	if name == "" || name == "<invalid Value>" {
		typ := reflect.TypeOf(obj)
		name = typ.Name()
		if name == "" {
			name = typ.String()
		}
	}
	return
}

// Description of Object. Return value from description field if available or return its type
func Description(obj interface{}) (description string) {
	val := reflect.Indirect(reflect.ValueOf(obj)).FieldByName("Description")
	description = val.String()
	if description == "<invalid Value>" {
		description = ""
	}
	return
}
