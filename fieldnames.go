package main

// mapping of csv column numbers to their respective field names
type FieldnameMap map[int]string

func NewFieldnameMap() (*FieldnameMap) {
	f := &FieldnameMap{}
	return f
}

// csv.Reader.Read() returns []string
// This is a helper method to extract fields names from the slice
func (f FieldnameMap) Update(s []string) (*FieldnameMap, error) {
	for i, fieldname := range s {
		f[i] = fieldname
	}
	return &f, nil
}