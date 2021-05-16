package typeTrans

import (
	"fmt"
	"reflect"
	"time"
)

type Trans interface {
	DataTypeof(typ reflect.Value) string
}

func GetTrans() Trans {
	return &mySql{}
}

type mySql struct {
}

var _ Trans = (*mySql)(nil)

func (m *mySql) DataTypeof(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool, reflect.Int8:
		return "tinyint"
	case reflect.Int, reflect.Int64:
		return "bigint"
	case reflect.Int32, reflect.Uint16:
		return "int"
	case reflect.Int16, reflect.Uint8:
		return "smallint"
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return "bigint"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "set"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datatime"
		}
	}
	panic(fmt.Sprintf("数据类型转换错误 %s (%s)", typ.Type().Name(), typ.Kind()))
	return ""
}
