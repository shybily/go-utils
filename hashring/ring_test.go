package hashring

import (
	"testing"
)

func TestRing_Locate(t *testing.T) {
	nodes := []string{
		"a", "b", "c",
	}

	ring := NewRing(1024, nodes)

	k := []string{
		"asdasd",
		"dsdfds",
		"kgfiwer",
		"oqwel",
	}
	for _, s := range k {
		n := ring.Locate(s)
		println(n)
	}
}
