package main

import (
	"testing"
)

func TestStringSet(t *testing.T) {
	ss := NewStringSet()

	expectedList := []struct{
		input string
		result bool 
	}{
		{"one", true, },
		{"two", false, },
	}

	ss.Add("one")

	for _, expected := range expectedList {
		if ss.Has(expected.input) != expected.result {
			t.Errorf("For key=%s, Expected Has() to return %v, but %v instead", expected.input, ss.Has(expected.input), expected.result)
		}
	}
}

func TestEdgeSet(t *testing.T) {
	es := NewEdgeSet()

	expectedList := []struct{
		input AnnotatedEdge
		result bool 
	}{
		{ AnnotatedEdge{"Src1","Dst1","Label1"}, true, },
		{ AnnotatedEdge{"Src2","Dst2","Label2"}, false, },
	}

	es.Add(expectedList[0].input)

	for _, expected := range expectedList {
		if es.Has(expected.input) != expected.result {
			t.Errorf("For key=%s, Expected Has() to return %v, but %v instead", expected.input, es.Has(expected.input), expected.result)
		}
	}
}