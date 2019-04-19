package config

import (
	"reflect"
	"strconv"
)

// InfoDetail contain information detail
type InfoDetail struct {
	Name     string
	Type     string
	Required string
	Default  string
}

// Informations return details of configuration
func Informations() []InfoDetail {
	var slice []InfoDetail
	informations(&slice, &Config{})
	return slice
}

func informations(slice *[]InfoDetail, obj interface{}) {
	elem := reflect.ValueOf(obj).Elem()

	for i := 0; i < elem.NumField(); i++ {
		fieldValue := elem.Field(i)
		fieldType := elem.Type().Field(i)

		if fieldValue.Kind() == reflect.Struct && fieldType.Anonymous {
			informations(slice, fieldValue.Addr().Interface())
		}

		tag := fieldType.Tag

		name := tag.Get("envconfig")
		ignored, _ := strconv.ParseBool(tag.Get("ignored"))
		if name == "" || ignored {
			continue
		}

		*slice = append(*slice, InfoDetail{
			Name:     name,
			Type:     fieldValue.Type().String(),
			Required: tag.Get("required"),
			Default:  tag.Get("default"),
		})
	}

}
