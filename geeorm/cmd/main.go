package main

import (
	"fmt"
	"geeorm"
	"geeorm/dialect"
	"geeorm/log"
)

func main() {
	options := make(map[string]string)
	options["charset"] ="utf8mb4"
	options["loc"] = "Local"
	options["parseTime"] = "true"
	dsn := dialect.DSN{
		Username: "zhuzi",
		Password: "123456",
		Database: "blog",
		Options: options,
	}
	engine,err := geeorm.NewEngine("mysql",dsn)
	if err != nil {
		log.Error(err)
	}
	defer engine.Close()
	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS user;").Exec()
	s.Raw("create table user(name text);").Exec()
	result,err := s.Raw("insert into user(`name`) values (?),(?)","Tom","sam").Exec()
	if err != nil {
		log.Error(err)
	}
	count,err := result.RowsAffected()
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("exec success,%d affected\n",count)
}
