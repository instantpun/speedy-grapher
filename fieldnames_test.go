package speedygrapher

import (
	"testing"
)

func TestFieldnameMap(t *testing.T) {
	input := []string{ "field1", "field2", "field3", }
	expected := FieldnameMap{
		0: "field1", 
		1: "field2", 
		2: "field3",
	}

	fmap := NewFieldnameMap()
	result, err := fmap.UpdateFromSlice(input)
	if err != nil {
		t.Errorf("Received error from FieldnameMap.Update(): %v\n", err)
	}

	if len(*result) != len(expected) {
		t.Error("The length of input Slice should match length of result FieldnameMap.\n")
		panic("Exiting!")
	}

	for k, v := range expected {
		resVal, ok := (*result)[k]
		if !ok {
			t.Errorf("Missing key: Expected %s but got %s.\n", v, resVal)
		}
		if resVal != v {
			t.Errorf("Mismatched key: Expected %s but got %s.\n", v, resVal)
		}
	}

}