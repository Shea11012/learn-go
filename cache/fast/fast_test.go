package fast

import (
	"fmt"
	"testing"
)

func TestHasher(t *testing.T) {
	hasher := newDefaultHasher()
	key := hasher.Sum64("k1")
	var mask uint64 = 16
	fmt.Println(key % mask)
	fmt.Println(key & (mask - 1))
}
