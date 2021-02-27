package session

import (
	"geeorm/log"
	"github.com/stretchr/testify/require"
	"testing"
)

type Account struct {
	ID int `geeorm:"PRIMARY KEY"`
	Password string `geeorm:"size:255"`
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("before insert",a)
	a.ID += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) error {
	log.Info("after query",a)
	a.Password = "***"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	s.DropTable()
	s.CreateTable()
	s.Insert(&Account{1,"123456"},&Account{2,"qwera"})
	u := &Account{}
	err := s.First(u)
	require.Nil(t, err)
	require.Equal(t, 1001,u.ID)
	require.Equal(t, "***",u.Password)
}
