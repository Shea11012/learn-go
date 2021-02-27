package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string,dsn dialect.DSN) (*Engine,error) {
	db,err := sql.Open(driver,dsn.String())
	if err != nil {
		log.Error(err)
		return nil,err
	}

	if err := db.Ping();err != nil {
		log.Error(err)
		return nil,err
	}

	dial,ok := dialect.GetDialect(driver)
	if !ok {
		return nil,fmt.Errorf("dialect %s not found",driver)
	}
	dial.SetSchema(dsn.Database)

	engine := &Engine{
		db:db,
		dialect: dial,
	}
	log.Info("connect database success")
	return engine,nil
}

func (e *Engine) Close() {
	if err := e.db.Close();err != nil {
		log.Error("failed to close database")
	}

	log.Info("close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db,e.dialect)
}

type TxFunc func(s *session.Session) (interface{},error)

func (e *Engine) Transaction(f TxFunc) (result interface{},err error) {
	s := e.NewSession()
	if err := s.Begin();err != nil {
		return nil,err
	}

	defer func() {
		if p:= recover();p!= nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			// 如果提交失败，再次回滚
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()

	return f(s)
}