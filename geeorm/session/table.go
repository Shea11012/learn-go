package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value,s.dialect)
	}

	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}

	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns = make([]string,0,len(table.Fields))
	for _,field := range table.Fields {
		f := *field
		columns = append(columns,schema.Column(&f))
	}

	desc := strings.Join(columns,",")
	_,err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)",table.Name,desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_,err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s",s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql,values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql,values...).QueryRaw()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == strings.ToLower(s.RefTable().Name)
}
