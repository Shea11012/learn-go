package dialect

import (
	"fmt"
	"net/url"
	"reflect"
)

type Dialect interface {
	DataTypeOf(value reflect.Value) string
	TableExistSQL(string) (string,[]interface{})
	SetSchema(schema string)
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	d,ok := dialectsMap[name]
	return d,ok
}

type DSN struct {
	Username string
	Password string
	Database string
	Options map[string]string
}

func (d DSN) String() string {
	urlV := url.Values{}
	for k,v := range d.Options {
		urlV.Add(k,v)
	}

	return fmt.Sprintf("%s:%s@/%s?%s",d.Username,d.Password,d.Database,urlV.Encode())
}
