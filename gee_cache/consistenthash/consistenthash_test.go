package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i,_ := strconv.Atoi(string(key))
		return uint32(i)
	})

	hash.Add("6","4","2")

	testcases := map[string]string{
		"2":"2",
		"11":"2",
		"23":"4",
		"27":"2",
	}

	for k,v := range testcases {
		if r := hash.Get(k); r != v {
			t.Errorf("should %s,get %s",v,r)
		}
	}

	// hash.Add("8")
}
