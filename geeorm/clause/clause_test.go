package clause

import (
	"reflect"
	"testing"
)

func testselect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT,3)
	clause.Set(SELECT,"user",[]string{"*"})
	clause.Set(WHERE,"name = ?","Tom")
	clause.Set(ORDERBY,"age asc")
	sql,vars := clause.Build(SELECT,WHERE,ORDERBY,LIMIT)
	t.Log(sql,vars)
	if sql != "SELECT * FROM user WHERE name = ? ORDER BY age asc LIMIT ?" {
		t.Fatal("failed to build SQL")
	}

	if !reflect.DeepEqual(vars,[]interface{}{"Tom",3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testselect(t)
	})
}