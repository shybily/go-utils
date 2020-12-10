package matchs

import (
	"sync"
	"testing"
)

func TestNewAcMatcher(t *testing.T) {
	w := []string{
		"a", "b",
	}
	m := NewAcMatcher(w)
	c := []string{
		"aaa", "bbb", "aabab", "vc", "ccsad", "aasd",
	}
	r := make([]interface{}, len(c))
	wg := sync.WaitGroup{}
	for i := range c {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r[i] = m.Match(c[i])
		}(i)
	}
	wg.Wait()
	println(r)
}
