package dialect

import (
	"fmt"
	"reflect"
	"time"
)

func init() {
	RegisterDialect("mysql",&mysql{})
}

type mysql struct {
	schema string
}

func (m *mysql) DataTypeOf(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "varchar"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _,ok := value.Interface().(time.Time); ok {
			return "datetime"
		}
	}

	panic(fmt.Sprintf("invalid sql type %s (%s)",value.Type().Name(),value.Kind()))
}

func (m *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{m.schema,tableName}
	return "SELECT TABLE_NAME FROM information_schema.`TABLES` WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = ? AND TABLE_NAME = ?;",args
}

func (m *mysql) SetSchema(schema string) {
	m.schema = schema
}

var _ Dialect = (*mysql)(nil)
