package goset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSet(t *testing.T) {
	assert.NotNil(t, BuildSet())
}

func TestBuildFromRegularSlice(t *testing.T) {
	a := []string{"321", "dsafgdsg", "erwe", "utryrt", "6844"}
	s := BuildFromRegularSlice(a)
	assert.Equal(t, 5, s.GetLength())

	b := []int{1, 2, 3, 4, 5}
	s = BuildFromRegularSlice(b)
	assert.Equal(t, 5, s.GetLength())

	c := []bool{false, false, false, true, true}
	s = BuildFromRegularSlice(c)
	assert.Equal(t, 5, s.GetLength())

	d := []float64{1.0, 22.5, 3.45, 43.486, 5.456678658}
	s = BuildFromRegularSlice(d)
	assert.Equal(t, 5, s.GetLength())

	e := []interface{}{1.0, "22.5", true, 43.486, struct{ Name string }{Name: "This is stuct"}}
	s = BuildFromRegularSlice(e)
	assert.Equal(t, 5, s.GetLength())
}

func TestAddSetIsExist(t *testing.T) {
	s := BuildSet()
	for i := 1; i < 1000; i++ {
		go func(idx int, s *Set) {
			s.Add(idx)
			assert.Equal(t, true, s.IsExist(idx))
		}(i, &s)
	}
	for i := 1; i < 1000; i++ {
		go func(idx int, s *Set) {
			s.Add(fmt.Sprint(idx))
			assert.Equal(t, true, s.IsExist(fmt.Sprint(idx)))
		}(i, &s)
	}
}
func TestGetLengthIsEmpty(t *testing.T) {
	s := BuildSet()
	assert.Equal(t, 0, s.GetLength())
	assert.Equal(t, true, s.IsEmpty())
	N := 1000
	for i := 0; i < N; i++ {
		s.Add(i)
	}
	assert.Equal(t, N, s.GetLength())
}

func TestRemove(t *testing.T) {
	s := BuildSet()
	s.Remove(0) // no action
	N := 10000
	for i := 0; i < N; i++ {
		s.Add(i)
		s.Add(fmt.Sprint(i))
	}
	for i := 1; i < N; i++ {
		go func(idx int, s *Set) {
			s.Remove(idx)
			assert.Equal(t, false, s.IsExist(idx))
			assert.Equal(t, true, s.IsExist(fmt.Sprint(idx)))
		}(i, &s)
	}
}

func TestHasSubSet(t *testing.T) {
	s := BuildSet()
	a := BuildSet()
	s.Add(1)
	s.Add(2)
	s.Add(4)
	s.Add(5)
	assert.Equal(t, true, s.HasSubSet(a))
	a.Add(2)
	a.Add(5)
	assert.Equal(t, true, s.HasSubSet(a))
	assert.Equal(t, false, a.HasSubSet(s))
	a.Add(8)
	assert.Equal(t, false, s.HasSubSet(a))
	s.Add("sdfsdflijsdfolj")
	b := BuildSet()
	b.Add("sdfsdflijsdfolj")
	assert.Equal(t, true, s.HasSubSet(b))
}

func TestToSlice(t *testing.T) {
	s := BuildSet()
	for i := 0; i < 50; i++ {
		s.Add(i)
	}
	arr := s.ToIntSlice()
	assert.Equal(t, 50, len(arr))

	s = BuildSet()
	for i := 0; i < 50; i++ {
		s.Add(float64(i) * 0.01)
	}
	arr1 := s.ToFloatSlice()
	assert.Equal(t, 50, len(arr1))

	s = BuildSet()
	for i := 0; i < 50; i++ {
		s.Add(fmt.Sprint(i) + "w")
	}
	arr2 := s.ToStrSlice()
	assert.Equal(t, 50, len(arr2))

	s = BuildSet()
	for i := 0; i < 50; i++ {
		s.Add(i)
	}
	arr3 := s.ToInterfaceSlice()
	assert.Equal(t, 50, len(arr3))
}

func TestClear(t *testing.T) {
	s := BuildSet()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	assert.Equal(t, 10, s.GetLength())
	s.Clear()
	assert.Equal(t, true, s.IsEmpty())
}

func TestUnion(t *testing.T) {
	a := BuildSet()
	b := BuildSet()
	for i := 0; i < 10; i++ {
		a.Add(i)
	}
	for i := 10; i < 20; i++ {
		b.Add(i)
	}
	c := Union(a, b)
	for i := 0; i < 20; i++ {
		assert.Equal(t, true, c.IsExist(i))
	}
	assert.Equal(t, false, c.IsExist(20))
}

func TestIntersection(t *testing.T) {
	a := BuildSet()
	b := BuildSet()
	for i := 0; i < 10; i++ {
		a.Add(i)
	}
	for i := 10; i < 20; i++ {
		b.Add(i)
	}
	c := Intersection(a, b)
	assert.Equal(t, 0, c.GetLength())
	for i := 5; i < 10; i++ {
		c.Add(i)
	}
	c = Intersection(a, c)
	assert.Equal(t, 5, c.GetLength())
}

func TestExclude(t *testing.T) {
	a := BuildSet()
	b := BuildSet()
	for i := 0; i < 100; i++ {
		a.Add(i)
	}
	for i := 0; i < 95; i++ {
		b.Add(i)
	}
	c := Exclude(a, b)
	assert.Equal(t, 5, c.GetLength())
}
