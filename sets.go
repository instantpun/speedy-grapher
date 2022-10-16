package main

import (
    "sync"
)

type StringSetter interface {
    Add(elem string)
    Has(elem string) bool
}

type StringSet map[string]struct{}

type StringSetAsync struct {
	Set StringSet
	sync.RWMutex
}

type AnnotatedEdge struct { Src, Dst, Label string }

type EdgeSet map[AnnotatedEdge]struct{}

type EdgeSetAsync struct {
	Set EdgeSet
	sync.RWMutex
}

func NewStringSet() (*StringSet) {
	s := &StringSet{}
	return s
}

func NewStringSetAsync() (*StringSetAsync) {
	s := &StringSetAsync{}
	s.Set = *NewStringSet()
	return s
}

func NewEdgeSet() (*EdgeSet) {
	s := &EdgeSet{}
	return s	
}

func NewEdgeSetAsync() (*EdgeSetAsync) {
	s := &EdgeSetAsync{}
	s.Set = *NewEdgeSet()
	return s
}

func (s *StringSet) Add(elem string) {
	if s == nil {
		panic("Error within StringSet.Add(): *StringSet is nil pointer")
	}

	(*s)[elem] = struct{}{} // only need to assign the key, the value is irrelevant
	return
}

func (s *StringSet) Has(elem string) bool {
	if s == nil {
		panic("Error within StringSet.Add(): *StringSet is nil pointer")
	}

	_, ok := (*s)[elem]
	return ok
}

func (s *StringSetAsync) Add(elem string) {
	if s == nil {
		panic("Error within StringSetAsync.Add(): *StringSetAsync is nil pointer")
	}

	s.Lock()
	defer s.Unlock()

	s.Set[elem] = struct{}{} // only need to assign the key, the value is irrelevant
	return
}

func (s *StringSetAsync) Has(elem string) bool {
	if s == nil {
		panic("Error within StringSetAsync.Add(): *StringSetAsync is nil pointer")
	}

	s.RLock()
	defer s.RUnlock()
	_, ok := s.Set[elem]
	return ok
}


func (s *EdgeSet) Add(elem AnnotatedEdge) {
	if s == nil {
		panic("Error within EdgeSet.Add(): *EdgeSet is nil pointer")
	}

	(*s)[elem] = struct{}{} // only need to assign the key, the value is irrelevant
	return
}

func (s *EdgeSet) Has(elem AnnotatedEdge) bool {
	if s == nil {
		panic("Error within EdgeSet.Add(): *EdgeSet is nil pointer")
	}

	_, ok := (*s)[elem]
	return ok
}

func (s *EdgeSetAsync) Add(elem AnnotatedEdge) {
	if s == nil {
		panic("Error within EdgeSetAsync.Add(): *EdgeSetAsync is nil pointer")
	}

	s.Lock()
	defer s.Unlock()

	s.Set[elem] = struct{}{} // only need to assign the key, the value is irrelevant
	return
}

func (s *EdgeSetAsync) Has(elem AnnotatedEdge) bool {
	if s == nil {
		panic("Error within EdgeSetAsync.Add(): *EdgeSetAsync is nil pointer")
	}

	s.RLock()
	defer s.RUnlock()
	_, ok := s.Set[elem]
	return ok
}
