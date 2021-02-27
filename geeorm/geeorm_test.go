package geeorm

import (
	"errors"
	"geeorm/dialect"
	"geeorm/session"
	"github.com/stretchr/testify/require"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	options := make(map[string]string)
	options["loc"] = "Local"
	options["charset"] = "utf8mb4"
	options["parseTime"] = "true"
	dsn := dialect.DSN{
		Username: "zhuzi",
		Password: "123456",
		Database: "blog",
		Options:options,
	}
	engine,err := NewEngine("mysql",dsn)
	if err != nil {
		t.Fatal("failed to connect",err)
	}

	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY,size:255"`
	Age int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})

	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_,err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_,err := s.Insert(&User{"Tom",18})
		if err != nil {
			return nil,err
		}

		return nil,nil
	})

	u := &User{}
	_ = s.First(u)
	require.Nil(t, err)
	require.Equal(t, "Tom",u.Name)
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_,err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_,err := s.Insert(&User{"Tom",18})
		if err != nil {
			return nil,err
		}

		return nil,errors.New("Error")
	})

	require.NotNil(t, err)
	require.True(t, s.HasTable())
}
