package session

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var (
	TestDB *sql.DB
	TestDial,_ = dialect.GetDialect("mysql")
)

type User struct {
	Name string `geeorm:"PRIMARY KEY,size:255"`
	Age uint8
}

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&loc=Local&parseTime=true","zhuzi","123456","blog")
	var err error
	TestDB,err = sql.Open("mysql",dsn)
	if err != nil {
		log.Fatal("connect error",err)
	}

	if err = TestDB.Ping(); err != nil {
		log.Fatal("ping error",err)
	}

	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	TestDial.SetSchema("blog")
	return &Session{
		db: TestDB,
		dialect: TestDial,
	}
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("failed to create table user")
	}
}
