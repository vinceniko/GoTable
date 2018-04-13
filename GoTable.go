// Copyright @ Vincent Nikolayev, 2018

// Using the following table for testing:
// +--------+-----+-------+
// | STRING | INT | FLOAT |
// +--------+-----+-------+
// | eff    |   1 |   4.2 |
// | efe    |   3 |  5.32 |
// | efe    |   2 |  1.32 |
// | ffs    |  52 |   2.1 |
// | wg     |  34 |    .8 |
// | ret    |   4 |   9.6 |
// +--------+-----+-------+

package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type converter2D interface {
	convert2D() [][]interface{}
}

func convert2D(c converter2D) [][]interface{} {
	return c.convert2D()
}

// string2D is used for polymorphic transformations of 2D strings slice to 2D interface slices
type string2D struct {
	slice *[][]string
}

func (ss string2D) convert2D() [][]interface{} {
	slice := make([][]interface{}, len(*ss.slice))
	for i, row := range *ss.slice {
		slice[i] = make([]interface{}, len(row))
		for j, cell := range row {
			slice[i][j] = cell
		}
	}
	return slice
}

type interface2D struct {
	slice *[][]interface{}
}

func (i interface2D) convert2D() [][]interface{} {
	return *i.slice
}

type converter1D interface {
	convert1D() []interface{}
}

func convert1D(c converter1D) []interface{} {
	return c.convert1D()
}

// string1D is used for polymorphic transformation of a 1D string slice to a 1D interface slice
type string1D struct {
	slice *[]string
}

func (ss string1D) convert1D() []interface{} {
	slice := make([]interface{}, len(*ss.slice))
	for i, val := range *ss.slice {
		slice[i] = val
	}
	return slice
}

func ConvertToString2D(it [][]interface{}) [][]string {
	slice := make([][]string, len(it))
	for i, row := range it {
		slice[i] = make([]string, len(row))
		for j, cell := range row {
			switch convert := cell.(type) {
			case string:
				slice[i][j] = convert
			case int:
				slice[i][j] = strconv.Itoa(convert)
			case nil:
				slice[i][j] = ""
			}
		}
	}
	return slice
}

func ConvertToString1D(it []interface{}) []string {
	slice := make([]string, len(it))
	for i, val := range it {
		switch convert := val.(type) {
		case string:
			slice[i] = convert
		case int:
			slice[i] = strconv.Itoa(convert)
		case nil:
			slice[i] = ""
		}
	}
	return slice
}

// Table is the main struct which is defined by a header that is a slice of column names, an index of the NamedVector type, and body values.
type Table struct {
	Header MappedSlice
	Index  MappedSlice
	Vals   [][]interface{}
}

type MappedSlice struct {
	Header  interface{}
	Map     map[interface{}][]int // holds indexes of names matched against vals
	Slice   []interface{}
	Length  int
	counter map[interface{}]int
}

func CreateMS(it []interface{}, header interface{}) MappedSlice {
	ms := MappedSlice{Length: 0, Header: header}

	for _, val := range it {
		ms.AddVal(val)
	}

	return ms
}

func CreateHeadMS(it []interface{}) MappedSlice {
	ms := CreateMS(it[1:], it[0])

	return ms
}

func CreateGenMS(axis int, it []interface{}) MappedSlice {
	var ms MappedSlice
	if axis == 0 {
		ms = CreateMS(it, "Index")
	} else if axis == 1 {
		ms = CreateMS(it, "Columns")
	}

	return ms
}

func CreateNumMS(axis int, index int) MappedSlice {
	ms := CreateGenMS(axis, rangeUntil(index))

	return ms
}

func (ms *MappedSlice) AddVal(val interface{}) {
	if ms.Map == nil {
		ms.Map = make(map[interface{}][]int) // initializes map if not already initialized
	}
	if _, ok := ms.Map[val]; ok { // key exists
		ms.Map[val] = append(ms.Map[val], ms.Length)
	} else {
		ms.Map[val] = []int{ms.Length}
	}
	ms.Slice = append(ms.Slice, val)
	ms.Length += 1 // add to length
}

// FromSlice creates a table from a slice. If header is true, the first nested slice is taken as a list of column headers. If index is true, the first element of each slice is taken as the list of index values
func FromSlice(c converter2D, header bool, index bool) *Table {
	vals := convert2D(c) // converts to [][]interface{}

	newHeader := MappedSlice{Header: "Columns"}
	newIndex := MappedSlice{Header: "Index"}
	var newVals [][]interface{}
	if index && header {
		newIndex = CreateHeadMS(GetTranspose(vals, 0))
		newHeader = CreateGenMS(1, vals[0][1:])
		newVals = SliceTranspose(SliceTranspose(vals[1:])[1:])
	} else if header && !index {
		newIndex = CreateNumMS(0, len(vals))
		newHeader = CreateGenMS(1, vals[0])
		newVals = vals[1:]
	} else if !header && index {
		newIndex = CreateGenMS(0, GetTranspose(vals, 0))
		newHeader = CreateNumMS(1, len(vals[0][1:]))
		newVals = SliceTranspose(SliceTranspose(vals)[1:])
	} else {
		newHeader = CreateNumMS(1, len(vals[0]))
		newIndex = CreateNumMS(0, len(vals))
		newVals = vals
	}

	t := &Table{
		Header: newHeader,
		Index:  newIndex,
		Vals:   newVals}
	return t
}

// FromMap creates a Table from a map along a given axis ie. keys become headers or index vals
func FromMap(axis int, m map[interface{}]interface{}) *Table {
	var t *Table

	mLength := len(m)

	if axis == 0 {

		index := make([]interface{}, mLength)
		vals := make([][]interface{}, mLength)
		var i int
		for key, row := range m { // seperate keys into own slice
			index[i] = key.(interface{})
			vals[i] = row.([]interface{})
			i++
		}

		slice := mergeIndex2D(index, vals)

		t = FromSlice(interface2D{&slice}, false, true)
	} else if axis == 1 {
		index := make([]interface{}, mLength)
		vals := make([][]interface{}, mLength)
		var i int
		for key, row := range m { // seperate keys and vals into own slices
			index[i] = key.(interface{})
			vals[i] = row.([]interface{})
			i++
		}
		slice := mergeIndex2D(index, vals)
		slice = SliceTranspose(slice)

		t = FromSlice(interface2D{&slice}, true, false)
	}
	return t
}

// FromCSVFile creates Table object from a csv file
func FromCSVFile(path string, Header bool, Index bool) *Table {
	file, err := os.Open(path)
	if err != nil {
		panic("panic")
	}

	reader := csv.NewReader(file)

	data, err := reader.ReadAll()
	if err != nil {
		panic("panic")
	}

	defer file.Close()

	t := FromSlice(string2D{&data}, Header, Index)

	return t
}

// ResetIndex resets the index to the sequential form and replaces Table.Index.Name with "Index"
// +-------+--------+-----+-------+
// | INDEX | STRING | INT | FLOAT |
// +-------+--------+-----+-------+
// |     0 | eff    |   1 |   4.2 |
// |     1 | efe    |   3 |  5.32 |
// |     2 | efe    |   2 |  1.32 |
// |     3 | ffs    |  52 |   2.1 |
// |     4 | wg     |  34 |    .8 |
// |     5 | ret    |   4 |   9.6 |
// +-------+--------+-----+-------+
func (t *Table) ResetIndex() {
	t.mergeBoth()
	numRows := t.Index.Length
	t.Index = CreateNumMS(0, numRows)
}

func (t *Table) SetIndex(column interface{}) {
	var col string

	// standardize lookup of column to name (not index)
	switch val := column.(type) {
	case string:
		col = val
	case int: // set by index of column
		col = t.Header.Slice[val].(string)
	}
	ms := CreateMS(t.GetCols(column)[0], col)
	t.mergeBoth() // transfer current index into cols
	t.Index = ms  // set new index
	t.DropCol(col)
}

func (t *Table) DropCol(column interface{}) {
	outHeader := sliceWOelement(t.Header.Slice, column) // remove column name
	*t = *t.GenSliceLoc(1, outHeader...)
}

// GetCols returns column values
func (t *Table) GetCols(columns ...interface{}) [][]interface{} {
	return SliceTranspose(t.GenSliceLoc(1, columns...).Vals)
}

// sliceWOelement returns a slice without a specific element
func sliceWOelement(slice []interface{}, value interface{}) []interface{} {
	outSlice := make([]interface{}, len(slice)-1)
	var i int
	for index, elem := range slice {
		switch val := value.(type) {
		case string:
			if elem != val { // match element
				outSlice[i] = elem
				i++
			}
		case int: // do not return a specific index
			if index != val {
				outSlice[i] = elem
				i++
			}
		}
	}
	return outSlice
}

// rangeUntil creates an ascending slice of interfaces of strings (converted from int for later printing) until Index paramater
func rangeUntil(Index int) []interface{} {
	s := make([]interface{}, Index)
	for i := 0; i < Index; i++ {
		s[i] = strconv.Itoa(i)
	}
	return s
}

// GetTranspose returns a transposed slice along a given Index in a set of row slices. ie. if the slices are oriented by rows, it returns a column at a given column Index
func GetTranspose(s [][]interface{}, index int) []interface{} {
	slice := make([]interface{}, len(s))
	for i, row := range s {
		slice[i] = row[index]
	}
	return slice
}

// SliceTranspose transposes a 2D slice in its entirety using GetTranspose
func SliceTranspose(s [][]interface{}) [][]interface{} {
	slice := make([][]interface{}, len(s[0]))
	for i := range s[0] {
		slice[i] = GetTranspose(s, i)
	}
	return slice
}

// Transpose transposes the entire table. indexHeader param is used to set the new indexHeader
// +---------+-----+------+------+-----+----+-----+
// | COLUMNS | EFF | EFE  | EFE  | FFS | WG | RET |
// +---------+-----+------+------+-----+----+-----+
// | Int     |   1 |    3 |    2 |  52 | 34 |   4 |
// | Float   | 4.2 | 5.32 | 1.32 | 2.1 | .8 | 9.6 |
// +---------+-----+------+------+-----+----+-----+
func (t *Table) Transpose() *Table {
	t0 := Table{}
	t0.Vals = SliceTranspose(t.Vals)
	oldHeader := t.Header
	t0.Header = t.Index
	t0.Index = oldHeader

	return &t0
}

// mergeIndex2D is used to merge Table.Index.vals with Table.Vals for resetting the Table index
func mergeIndex2D(index []interface{}, vals [][]interface{}) [][]interface{} {
	v := make([][]interface{}, len(vals))
	for i, row := range vals {
		v1 := make([]interface{}, len(row)+1)
		v1[0] = index[i]
		for j, cell := range row {
			v1[j+1] = cell
		}
		v[i] = v1
	}
	return v
}

// mergeIndex1D is used to merge Table.Index.Header with Table.Header for resetting the Table index
func mergeIndex1D(index interface{}, vals []interface{}) []interface{} {
	v := make([]interface{}, len(vals)+1)
	v[0] = index
	for i, val := range vals {
		v[i+1] = val
	}
	return v
}

// mergeBoth calls mergeIndex1D and mergeIndex2D and reassigns Table fields inplace. Table.Index and Table.Header remain unchanged
func (t *Table) mergeBoth() {
	newHeader := mergeIndex1D(t.Index.Header, t.Header.Slice)
	t.Header = CreateGenMS(1, newHeader)

	newVals := mergeIndex2D(t.Index.Slice, t.Vals)
	t.Vals = newVals
}

// PrintTable prints the table using non-std Ascii Table package
func (t *Table) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(
		ConvertToString1D(
			mergeIndex1D(t.Index.Header, t.Header.Slice)))
	table.AppendBulk(
		ConvertToString2D(
			mergeIndex2D(t.Index.Slice, t.Vals))) // Add Bulk Data
	table.Render()
}

func (t *Table) SliceLoc(axis int, names ...string) *Table {
	t0 := Table{}

	var m map[interface{}][]int
	var vals [][]interface{} = t.Vals
	var outVals [][]interface{}
	var outNames []interface{}

	for _, name := range names {
		if axis == 0 {
			m = t.Index.Map

			if indices, ok := m[name]; ok {
				for _, index := range indices {
					outVals = append(outVals, vals[index])
					outNames = append(outNames, name)
				}
			} else {
				panic("panic")
			}

		} else if axis == 1 {
			m = t.Header.Map

			if indices, ok := m[name]; ok {
				for _, index := range indices {
					outVals = append(outVals, GetTranspose(vals, index))
					outNames = append(outNames, name)
				}
			} else {
				panic("panic")
			}
		}
	}

	if axis == 0 {
		ms := CreateMS(outNames, t.Index.Header)
		t0.Index = ms
		t0.Header = t.Header
		t0.Vals = outVals
	} else if axis == 1 {
		ms := CreateMS(outNames, t.Header.Header)
		t0.Index = t.Index
		ms.Header = t.Header.Header
		t0.Header = ms
		t0.Vals = SliceTranspose(outVals)
	}

	return &t0
}

func (t *Table) Loc(rows []string, cols []string) *Table {
	t0 := t
	if len(rows) > 0 {
		t0 = t0.SliceLoc(0, rows...)
	}
	if len(cols) > 0 {
		t0 = t0.SliceLoc(1, cols...)
	}

	return t0
}

// SliceILoc returns an interface slice for a selected subset of indexed columns or Vals found in a Table object--basically a modified slice of Table.values matched against the correspinding Table.Index or Table.Header field (which is selected with the axis parameter). Panics if index is out of the bounds of the axis
func (t *Table) SliceILoc(axis int, indices ...int) *Table {
	t0 := Table{}

	var vals [][]interface{} = t.Vals
	outVals := make([][]interface{}, len(indices))
	outNames := make([]interface{}, len(indices))

	for i, index := range indices {
		if axis == 0 {
			outVals[i] = vals[index]
			outNames[i] = t.Index.Slice[index]
		} else if axis == 1 {
			vals = SliceTranspose(vals)
			outVals[i] = vals[index]
			outNames[i] = t.Header.Slice[index]
		}
	}

	if axis == 0 {
		ms := CreateMS(outNames, t.Index.Header)
		t0.Index = ms
		t0.Header = t.Header
		t0.Vals = outVals
	} else if axis == 1 {
		ms := CreateMS(outNames, t.Header.Header)
		t0.Index = t.Index
		ms.Header = t.Header.Header
		t0.Header = ms
		t0.Vals = SliceTranspose(outVals)
	}

	return &t0
}

// ILoc uses SliceILoc to find a selected subsections of indexed rows and columns on both axes (using Table.Index.Vals and Table.Header) and returns a new Table of the subsection
func (t *Table) ILoc(rows []int, cols []int) *Table {
	t0 := t
	if len(rows) > 0 {
		t0 = t0.SliceILoc(0, rows...)
	}
	if len(cols) > 0 {
		t0 = t0.SliceILoc(1, cols...)
	}

	return t0

}

func (t *Table) GenSliceLoc(axis int, values ...interface{}) *Table {
	t0 := Table{}

	var m map[interface{}][]int
	var vals [][]interface{} = t.Vals
	var outVals [][]interface{}
	var outNames []interface{}

	for _, val := range values {
		switch v := val.(type) {
		case string:
			if axis == 0 {
				m = t.Index.Map

				if indices, ok := m[v]; ok {
					for _, index := range indices {
						outVals = append(outVals, vals[index])
						outNames = append(outNames, v)
					}
				} else {
					panic("panic")
				}

			} else if axis == 1 {
				m = t.Header.Map

				if indices, ok := m[v]; ok {
					for _, index := range indices {
						outVals = append(outVals, GetTranspose(vals, index))
						outNames = append(outNames, v)
					}
				} else {
					panic("panic")
				}
			}
		case int:
			if axis == 0 {
				outVals = append(outVals, vals[v])
				outNames = append(outNames, t.Index.Slice[v])
			} else if axis == 1 {
				vals = SliceTranspose(vals)
				outVals = append(outVals, vals[v])
				outNames = append(outNames, t.Header.Slice[v])
			}
		}
	}

	if axis == 0 {
		ms := CreateMS(outNames, t.Index.Header)
		t0.Index = ms
		t0.Header = t.Header
		t0.Vals = outVals
	} else if axis == 1 {
		ms := CreateMS(outNames, t.Header.Header)
		t0.Index = t.Index
		ms.Header = t.Header.Header
		t0.Header = ms
		t0.Vals = SliceTranspose(outVals)
	}

	return &t0
}

func (t *Table) AddSlice(axis int, header interface{}, slice []interface{}) {
	var ms MappedSlice
	if axis == 0 {
		ms = t.Index
		ms.AddVal(header)
		t.Vals = append(t.Vals, slice)
		t.Index = ms
	} else if axis == 1 {
		ms = t.Header
		ms.AddVal(header)
		t.Vals = SliceTranspose(append(SliceTranspose(t.Vals), slice))
		t.Header = ms
	}
}

// // PairedSliceLoc is similar to the other SliceLoc functions but returns the passed in names and the results. PairedSliceLoc does not panic but instead returns empty slices where it could not find a name.
// [efe efe trer hello] [[3 5.32] [2 1.32] [<nil> <nil>] [<nil> <nil>]]
func (t *Table) PairedSliceLoc(axis int, vals ...interface{}) ([]interface{}, [][]interface{}) {
	slice := t.Vals
	var searchMap map[interface{}][]int
	if axis == 0 {
		searchMap = t.Index.Map
	} else if axis == 1 {
		searchMap = t.Header.Map
	} else {
		panic("panic")
	}

	// needs appends because unless we itterate through all indices in map, we wont know how big to make the returned slices
	outSlice := make([][]interface{}, 0)
	outvals := make([]interface{}, 0)
	for _, name := range vals {
		indices, err := searchMap[name]
		if err {
			for _, index := range indices {
				if axis == 0 {
					outSlice = append(outSlice, slice[index])
				} else if axis == 1 {
					outSlice = append(outSlice, GetTranspose(slice, index))
				}
				outvals = append(outvals, name)
			}
		} else { // not found in map
			outSlice = append(outSlice, make([]interface{}, len(slice[0]))) // append empty slice
			outvals = append(outvals, name)
		}
	}

	return outvals, outSlice // extract nested list
}

// Index returns Index of a specific string element in a slice
func Index(search interface{}, slice []interface{}) (int, error) {
	for i, val := range slice {
		if search == val {
			return i, nil
		}
	}
	return -1, errors.New("Search not found")
}

func createDupeMap(m map[interface{}]interface{}, keys []interface{}, slices [][]interface{}, length int, totalLength int, index *int) ([]interface{}, map[interface{}]interface{}) {
	counter := make(map[interface{}]int)
	newKeys := make([]interface{}, len(keys))

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		var k int
		if _, ok := m[key]; !ok { // if not in map
			m[key] = make([]interface{}, totalLength) // make nested slice which is the length of the axis
			for j := *index; j < *index+length; j++ { // access the correct nested slice
				m[key].([]interface{})[j] = slices[i][k] // change each element of the nested slice
				k++
			}
			counter[key] = 1
		} else if counter[key] == 0 { // if key exists in map but we are at a new table and counter has reset
			for j := *index; j < *index+length; j++ {
				m[key].([]interface{})[j] = slices[i][k]
				k++
			}
			counter[key] = 1
		} else { // key exists and counter has not been reset
			key = key.(string) + "_" + strconv.Itoa(counter[key]) // create dupe key name
			if _, ok := m[key]; !ok {                             // create slice for key if its not already there
				m[key] = make([]interface{}, totalLength)
			}
			for j := *index; j < *index+length; j++ { // change each element of nested list
				m[key].([]interface{})[j] = slices[i][k]
				k++
			}
		}
		k = 0 // reset the index position of slice
		newKeys[i] = key
	}
	*index += length
	return newKeys, m
}

func appendMapSlices(m map[interface{}]interface{}, key interface{}, slices []interface{}) map[interface{}]interface{} {
	if _, ok := m[key]; !ok {
		m[key] = slices
	} else {
		val := m[key].([]interface{})
		m[key] = append(val, slices...)
	}
	return m
}

// Appends nills to beginning of slice to reach the total length desired
func spacer(length int, slice []interface{}) []interface{} {
	s := make([]interface{}, length)

	lenSlice := len(slice)
	j := lenSlice - 1
	for i := length - 1; i > -1; i-- {
		if j > -1 {
			s[i] = slice[j]
			j--
		} else {
			s[i] = nil
		}
	}
	return s
}

func getuniqs(slice []interface{}) []interface{} {
	uniqs := make([]interface{}, 0)
	for _, val := range slice {
		if _, err := Index(val, uniqs); err != nil {
			uniqs = append(uniqs, val)
		}
	}
	return uniqs
}

// Concat concantenates multiple tables along a given axis and merges common rows and columns. If there are duplicates which interfere with merging, the first dupe is merged with and the subsequent one is renamed in a sequential fashion i.e. name -> name_1
// +--------+-----+-------+
// | STRING | INT | FLOAT |
// +--------+-----+-------+
// | eff    |   1 |   4.2 |
// | efe    |   3 |  5.32 |
// | efe    |   2 |  1.32 |
// | ffs    |  52 |   2.1 |
// | wg     |  34 |    .8 |
// | ret    |   4 |   9.6 |
// +--------+-----+-------+
// 			  (+)
// +--------+-----+-------+
// | STRING | INT | FLOAT |
// +--------+-----+-------+
// | eff    |   2 |  34.3 |
// | efe    |   8 |   7.2 |
// | efe    |   2 |   6.2 |
// | ffs    |   4 |  7.47 |
// | wg     |   5 |   7.5 |
// | gr     |   8 |  56.7 |
// | vin    |   9 |  1.23 |
// +--------+-----+-------+
// 			  (=)
// axis == 0:
// --------+-----+-------+-----+-------+
// | STRING | INT | FLOAT | INT | FLOAT |
// +--------+-----+-------+-----+-------+
// | eff    |   1 |   4.2 |   2 |  34.3 |
// | efe    |   3 |  5.32 |   8 |   7.2 |
// | efe_1  |   2 |  1.32 |   2 |   6.2 |
// | ffs    |  52 |   2.1 |   4 |  7.47 |
// | wg     |  34 |    .8 |   5 |   7.5 |
// | ret    |   4 |   9.6 |     |       |
// | gr     |     |       |   8 |  56.7 |
// | vin    |     |       |   9 |  1.23 |
// +--------+-----+-------+-----+-------+
// axis == 1:
// +--------+-----+-------+
// | STRING | INT | FLOAT |
// +--------+-----+-------+
// | eff    |   1 |   4.2 |
// | efe    |   3 |  5.32 |
// | efe    |   2 |  1.32 |
// | ffs    |  52 |   2.1 |
// | wg     |  34 |    .8 |
// | ret    |   4 |   9.6 |
// | eff    |   2 |  34.3 |
// | efe    |   8 |   7.2 |
// | efe    |   2 |   6.2 |
// | ffs    |   4 |  7.47 |
// | wg     |   5 |   7.5 |
// | gr     |   8 |  56.7 |
// | vin    |   9 |  1.23 |
// +--------+-----+-------+
func Concat(axis int, tables ...*Table) *Table {
	t := tables[0]
	m := make(map[interface{}]interface{})

	if axis == 0 {
		var totalLength int
		var newHeaderVals []interface{}
		for _, table := range tables {
			totalLength += table.Header.Length
			newHeaderVals = append(newHeaderVals, table.Header.Slice...) // append all header names together
		}

		indexVals := make([]interface{}, 0)
		var newNames []interface{}

		currIndex := 0
		indexptr := &currIndex

		for _, table := range tables {
			var names []interface{}
			var slices [][]interface{}
			uniqs := getuniqs(table.Index.Slice)
			names, slices = table.PairedSliceLoc(0, uniqs...)

			newNames, m = createDupeMap(m, names, slices, table.Header.Length, totalLength, indexptr)
			indexVals = append(indexVals, newNames...)
		}

		indexVals = getuniqs(indexVals)

		t = FromMap(0, m)
		t.Index.Header = tables[0].Index.Header
		t.Header = CreateMS(newHeaderVals, tables[0].Header.Header)
		t = t.GenSliceLoc(0, indexVals...) // needed to rearrange into original order after map disrupts order
	} else if axis == 1 {

		var totalLength int
		var newIndexVals []interface{}
		for _, table := range tables {
			totalLength += table.Index.Length
			newIndexVals = append(newIndexVals, table.Index.Slice...) // append all index names together
		}

		headerVals := make([]interface{}, 0)
		var newNames []interface{}

		currIndex := 0
		indexptr := &currIndex
		for _, table := range tables {
			var names []interface{}
			var slices [][]interface{}
			uniqs := getuniqs(table.Header.Slice)
			names, slices = table.PairedSliceLoc(1, uniqs...)

			newNames, m = createDupeMap(m, names, slices, table.Index.Length, totalLength, indexptr)
			headerVals = append(headerVals, newNames...)
		}

		headerVals = getuniqs(headerVals)

		t = FromMap(1, m)
		t.Index = CreateMS(newIndexVals, tables[0].Index.Header)
		t = t.GenSliceLoc(1, headerVals...) // needed to rearrange into original order after map disrupts order
	}
	return t
}

func (t *Table) getAxisLabels(axis Axis) (headerName interface{}, headerSlice []interface{}) {
	if axis == 0 {
		headerName = t.Index.Header
		headerSlice = t.Index.Slice
	} else if axis == 1 {
		headerName = t.Header.Header
		headerSlice = t.Header.Slice
	}
	return
}

func (t *Table) getOrientation(axis Axis) *Table {
	if axis == 1 {
		t = t.Transpose()
	} else if axis != 0 {
		panic("panic")
	}
	return t
}

func (t *Table) getAxisAll(axis Axis) (headerName interface{}, headerSlice []interface{}, vals [][]interface{}) {
	headerName, headerSlice = t.getAxisLabels(axis)
	vals = t.getOrientation(axis).Vals

	return headerName, headerSlice, vals
}

type Axis uint8

// Axis.checkError checks to see whether axis is an int other than 0 and 1
func (a *Axis) checkError() {
	if (*a != 0) && (*a != 1) {
		panic("panic")
	}
}

// Axis.Opposite changes 0 to 1 and 1 to 0
func (a *Axis) Opposite() {
	a.checkError()
	if *a == 0 {
		*a = 1
	} else {
		*a = 0
	}
}

// ToMap converts a Table to a map along a given axis (discards other axis)
func (t *Table) ToMap(axis Axis) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})

	var headerName interface{}
	var headerSlice []interface{}
	var vals [][]interface{}

	headerName, headerSlice, vals = t.getAxisAll(axis)

	// index := t.Index.Slice
	for i := 0; i < len(headerSlice); i++ {
		m[headerSlice[i]] = vals[i]
	}
	axis.Opposite()
	headerName, headerSlice = t.getAxisLabels(axis)
	m[headerName] = headerSlice

	return m
}

func main() {
	testFor := FromCSVFile("Data/test.csv", true, true)
	testFor.PrintTable()
}
