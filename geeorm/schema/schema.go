package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  map[string]string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dst interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dst)).Type()
	schema := &Schema{
		Model:    dst,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = parseTag(v)
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues = make([]interface{},0,len(s.Fields))
	for _,field := range s.Fields {
		fieldValues = append(fieldValues,destValue.FieldByName(field.Name).Interface())
	}

	return fieldValues
}
