package session

import (
	"testing"
	"github.com/stretchr/testify/require"
)

var (
	user1 = &User{
		Name: "Tom",
		Age:  18,
	}

	user2 = &User{
		Name: "Sam",
		Age: 25,
	}

	user3 = &User{
		Name: "Jack",
		Age: 25,
	}
)

func testRecordinit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err := s.DropTable()
	require.Nil(t, err)
	err = s.CreateTable()
	require.Nil(t, err)
	_,err = s.Insert(user1,user2)
	require.Nil(t, err)
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordinit(t)
	affected,err := s.Insert(user3)
	require.Nil(t, err)
	require.Equal(t, int64(1),affected)
}

func TestSession_Find(t *testing.T) {
	s := testRecordinit(t)
	var users []User
	err := s.Find(&users)
	require.Nil(t, err)
	require.Equal(t, 2,len(users))
}

func TestSession_Limit(t *testing.T) {
	s := testRecordinit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	require.Nil(t, err)
	require.Equal(t, 1,len(users))
}

func TestSession_Update(t *testing.T) {
	s := testRecordinit(t)
	affected,_ := s.Where("Name = ?","Tom").Update("Age",30)
	u := &User{}
	_ = s.OrderBy("Age desc").First(u)

	require.Equal(t, int64(1),affected)
	require.Equal(t, uint8(30),u.Age)
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordinit(t)
	affected,_ := s.Where("Name = ?","Tom").Delete()
	count,_ := s.Count()
	require.Equal(t, int64(1),affected)
	require.Equal(t, int64(1),count)
}
