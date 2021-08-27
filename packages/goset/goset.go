package goset

import (
	"sync"
)

type Set struct {
	elem map[interface{}]struct{}
	*sync.RWMutex
}

/* Return an empty set */
func BuildSet() Set {
	return Set{
		elem:    map[interface{}]struct{}{},
		RWMutex: &sync.RWMutex{},
	}
}

/*
Create a set by a given slice.

Support: int, float64, string, bool, interface{}
*/
func BuildFromRegularSlice(a_slice interface{}) Set {
	s := Set{
		elem:    map[interface{}]struct{}{},
		RWMutex: &sync.RWMutex{},
	}
	switch a_slice := a_slice.(type) {
	case []int:
		for _, v := range a_slice {
			s.elem[v] = struct{}{}
		}
	case []float64:
		for _, v := range a_slice {
			s.elem[v] = struct{}{}
		}
	case []string:
		for _, v := range a_slice {
			s.elem[v] = struct{}{}
		}
	case []bool:
		for _, v := range a_slice {
			s.elem[v] = struct{}{}
		}
	case []interface{}:
		for _, v := range a_slice {
			s.elem[v] = struct{}{}
		}
	}
	return s
}

/* Add value(s) to set.*/
func (s *Set) Add(values interface{}) {
	s.Lock()
	defer s.Unlock()
	s.elem[values] = struct{}{}
}

/* Check value exists in set.*/
func (s *Set) IsExist(value interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, exist := s.elem[value]
	return exist
}

/* Return Set length .*/
func (s *Set) GetLength() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.elem)
}

/* Return IsEmpty or not .*/
func (s *Set) IsEmpty() bool {
	s.RLock()
	defer s.RUnlock()
	return len(s.elem) == 0
}

/* Delete given value(s) from map. There will be no action if value not exists in the set.*/
func (s *Set) Remove(values interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.elem, values)
}

/*Check if a is a subset of current set or not.*/
func (s *Set) HasSubSet(a Set) bool {
	s.RLock()
	a.RLock()
	defer s.RUnlock()
	defer a.RUnlock()
	for value := range a.elem {
		if _, exist := s.elem[value]; !exist {
			return false
		}
	}
	return true
}

/* Turn Set into slice */
func (s *Set) ToIntSlice() []int {
	s.RLock()
	defer s.RUnlock()
	arr := make([]int, 0, len(s.elem))
	for k := range s.elem {
		arr = append(arr, k.(int))
	}
	return arr
}

/* Turn Set into slice */
func (s *Set) ToFloatSlice() []float64 {
	s.RLock()
	defer s.RUnlock()
	arr := make([]float64, 0, len(s.elem))
	for k := range s.elem {
		arr = append(arr, k.(float64))
	}
	return arr
}

/* Turn Set into slice */
func (s *Set) ToStrSlice() []string {
	s.RLock()
	defer s.RUnlock()
	arr := []string{}
	for k := range s.elem {
		arr = append(arr, k.(string))
	}
	return arr
}

/* Turn Set into slice */
func (s *Set) ToInterfaceSlice() []interface{} {
	s.RLock()
	defer s.RUnlock()
	arr := make([]interface{}, 0, len(s.elem))
	for k := range s.elem {
		arr = append(arr, k)
	}
	return arr
}

/* delete all from slice */
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	for v := range s.elem {
		delete(s.elem, v)
	}
}

/* Return a + b (set operation) */
func Union(a Set, b Set) Set {
	s_union := Set{
		elem:    map[interface{}]struct{}{},
		RWMutex: &sync.RWMutex{},
	}
	for k := range a.elem {
		s_union.elem[k] = struct{}{}
	}
	for k := range b.elem {
		s_union.elem[k] = struct{}{}
	}
	return s_union
}

/* Return a & b (set operation) */
func Intersection(a Set, b Set) Set {
	s_intersect := Set{
		elem:    map[interface{}]struct{}{},
		RWMutex: &sync.RWMutex{},
	}
	if len(a.elem) > len(b.elem) {
		a, b = b, a
	}
	for k := range a.elem {
		if _, b_exist := b.elem[k]; b_exist {
			s_intersect.elem[k] = struct{}{}
		}
	}
	return s_intersect
}

/* Return a - b (set operation)*/
func Exclude(a Set, b Set) Set {
	s_exclude := Set{
		elem:    map[interface{}]struct{}{},
		RWMutex: &sync.RWMutex{},
	}
	for k := range a.elem {
		if _, b_exist := b.elem[k]; !b_exist {
			s_exclude.elem[k] = struct{}{}
		}
	}
	return s_exclude
}
