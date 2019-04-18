package config

import (
	"reflect"
	"strconv"
)

type ConfigDetail struct {
	Name     string
	Type     string
	Required string
	Default  string
}

func details(slice *[]ConfigDetail, obj interface{}) {
	elem := reflect.ValueOf(obj).Elem()

	for i := 0; i < elem.NumField(); i++ {
		fieldValue := elem.Field(i)
		fieldType := elem.Type().Field(i)

		if fieldValue.Kind() == reflect.Struct && fieldType.Anonymous {
			details(slice, fieldValue.Addr().Interface())
		}

		tag := fieldType.Tag

		name := tag.Get("envconfig")
		ignored, _ := strconv.ParseBool(tag.Get("ignored"))
		if name == "" || ignored {
			continue
		}

		*slice = append(*slice, ConfigDetail{
			Name:     name,
			Type:     fieldValue.Type().String(),
			Required: tag.Get("required"),
			Default:  tag.Get("default"),
		})

	}

}
