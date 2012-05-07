package skiplist

import "testing"
import "fmt"
import "math/rand"

func (s SkipList) printRepr() {

	for node := s.header; !node.isEnd(); node = node.forward[0] {
		fmt.Printf("%v: %v (level %d)\n", node.key, node.value, node.level())
		for i, link := range node.forward {
			fmt.Printf("\t%d: -> %v\n", i, link.key)
		}
	}
	fmt.Println()
}

func TestInitialization(t *testing.T) {
	s := New(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	if !s.lessThan(1, 2) {
		t.Errorf("Less than doesn't work correctly.")
	}
}

func TestIsEnd(t *testing.T) {
	s := NewIntKey()
	if !s.end.isEnd() {
	 	t.Errorf("isEnd() is not true for s.end.")
	}

	if s.header.isEnd() {
		t.Errorf("isEnd() is true for s.header.")
	}

	s.Set(0, 0)
	node := s.header.Next()
	if node.Key() != 0 {
		t.Fatalf("We got the wrong node: %v.", node)
	}

	if node.isEnd() {
		t.Errorf("isEnd() should not be true for %v.", node)
	}

	if node == s.end {
		t.Errorf("%v should not be equal to s.end.", node)
	}

	if node.Next() != s.end {
		t.Errorf("node.next should not be equal to s.end (was %v).", node, node.Next())
	}

}

func (s SkipList) check(t *testing.T, key, wanted int) bool {
	if got := s.Get(key); got != wanted {
		t.Errorf("Wanted %v, got %v.", wanted, got)
		return true
	}
	return false
}

func TestSet(t *testing.T) {
	s := NewIntKey()
	if l := s.Len(); l != 0 {
		t.Errorf("Len is not 0, it is %v", l)
	}

	s.Set(0, 0)
	s.Set(1, 1)
	if l := s.Len(); l != 2 {
		t.Errorf("Len is not 2, it is %v", l)
	}
	if s.check(t, 0, 0) {
		t.Errorf("%v", s.header.Next())
	}
	s.check(t, 1, 1)

}

func TestDelete(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	for i := 0; i < 10; i+=2 {
		s.Delete(i)
	}

	for i := 0; i < 10; i+=2 {
		if s.Get(i) != nil {
			t.Errorf("%d should not be present in s", i)
		}
	}
	if t.Failed() {
		s.printRepr()
	}
	
}

func TestLen(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	if length := s.Len(); length != 10 {
		t.Errorf("Length should be equal to 10, not %v.", length)
		s.printRepr()
	}
}

func TestIterator(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 20; i++ {
		s.Set(i, i)
	}

	seen := 0
	var lastKey int
	for i := s.Iter(); i.HasNext(); i = i.Next() {
		seen++
		lastKey = i.Key().(int)
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}
	}

	if seen != s.Len() {
		t.Errorf("Not all the items in s where iterated through (seen %d, should have seen %d). Last one seen was %d.", seen, s.Len(), lastKey)
	}
}

func TestSomeMore(t *testing.T) {
	s := NewIntKey()
	insertions := [...]int{4, 1, 2, 9, 10, 7, 3}
	for _, i := range insertions {
		s.Set(i, i)
	}
	for _, i := range insertions {
		s.check(t, i, i)
	}

}

func LookupBenchmark(b *testing.B, n int) {
	b.StopTimer()
	s := NewIntKey()
	for i := 0; i < n; i++ {
		insert := rand.Int()
		s.Set(insert, insert)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get(rand.Int())
	}
}


func BenchmarkLookup16(b *testing.B) {
	LookupBenchmark(b, 16)
}


func BenchmarkLookup256(b *testing.B) {
	LookupBenchmark(b, 256)
}


func BenchmarkLookup65536(b *testing.B) {
	LookupBenchmark(b, 65536)
}
