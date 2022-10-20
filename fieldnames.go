package speedygrapher


// mapping of csv column numbers to their respective field names
type FieldnameMap map[int]string // TODO: Support fieldnames from JSON

func NewFieldnameMap() (*FieldnameMap) {
	f := &FieldnameMap{}
	return f
}

// Note:
// Update expects a []string of fieldnames, and then
// imports each string and position into the map
//
// This most useful when paired with something like
// csv.Reader.Read(), which returns []string
func (f FieldnameMap) UpdateFromSlice(s []string) (*FieldnameMap, error) {
	for i, fieldname := range s {
		f[i] = fieldname
	}
	return &f, nil
}