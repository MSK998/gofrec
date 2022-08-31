package gofrec

import (
	"errors"
	"reflect"
	"strconv"
)

func DynamicType(t reflect.Type, fieldIndex int, v *reflect.Value, data string) error {

	switch t.Field(fieldIndex).Type.String() {
	case "string":
		v.Elem().Field(fieldIndex).SetString(data)
		return nil
	case "int", "int8", "int16", "int32", "int64":
		cVal, _ := strconv.Atoi(data)
		v.Elem().Field(fieldIndex).SetInt(int64(cVal))
		return nil
	case "float32", "float64":
		cVal, _ := strconv.ParseFloat(data, 64)
		v.Elem().Field(fieldIndex).SetFloat(cVal)
		return nil
	case "bool":
		cVal, _ := strconv.ParseBool(data)
		v.Elem().Field(fieldIndex).SetBool(cVal)
		return nil
	default:
		return errors.New("can't convert type")
	}
}