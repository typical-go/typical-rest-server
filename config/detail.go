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

func details(conf interface{}) (details []ConfigDetail) {
	val := reflect.ValueOf(conf).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		tag := typeField.Tag

		name := tag.Get("envconfig")
		ignored, _ := strconv.ParseBool(tag.Get("ignored"))
		if name == "" || ignored {
			continue
		}

		details = append(details, ConfigDetail{
			Name:     name,
			Type:     valueField.Type().String(),
			Required: tag.Get("required"),
			Default:  tag.Get("default"),
		})
	}

	return details
}
