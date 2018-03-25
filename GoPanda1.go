package main

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	Header interface{}
	Map    map[interface{}][]int // holds indexes of names matched against vals
	Slice  []interface{}
	Length int
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
func (t *Table) ResetIndex() {
	t.mergeBoth()
	numRows := t.Index.Length
	t.Index = CreateNumMS(0, numRows)
}

// ResetHeader resets the header to the sequential form
func (t *Table) ResetHeader() {
	numCols := t.Header.Length
	t.Header = CreateNumMS(0, numCols)
}

func (t *Table) SetIndex(column interface{}) {
	var col string

	// standardize lookup of column to name (not index)
	switch val := column.(type) {
	case string:
		col = val
	case int:
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
		case int:
			if index != val { // match index
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
		fmt.Print(t.Vals)
		t.Index = ms
	} else if axis == 1 {
		ms = t.Header
		ms.AddVal(header)
		t.Vals = SliceTranspose(append(SliceTranspose(t.Vals), slice))
		t.Header = ms
	}
}

// // PairedSliceLoc is similar to the other SliceLoc functions but returns the passed in names and the results. PairedSliceLoc does not panic but instead returns empty slices where it could not find a name
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

	outSlice := make([][]interface{}, 0)
	outvals := make([]interface{}, 0)
	var i int
	for _, name := range vals {
		indices, ok := searchMap[name]
		if ok == true {
			for _, index := range indices {
				if axis == 0 {
					outSlice = append(outSlice, slice[index])
				} else if axis == 1 {
					outSlice = append(outSlice, GetTranspose(slice, index))
				}
				outvals = append(outvals, name)
			}
		} else {
			outSlice = append(outSlice, make([]interface{}, len(slice[0])))
			outvals = append(outvals, name)
		}
		i++
	}

	return outvals[:i], outSlice[:i] // extract nested list
}

//
// // GenLoc is a generic implementation of Loc/ILoc using slice of interfaces (string or int)
// func (t *Table) GenLoc(cols []interface{}, rows []interface{}) *Table {
// 	t0 := t
//
// 	outHeader := t.Header
// 	outIndex := make([]interface{}, len(rows))
// 	if len(cols) > 0 {
// 		outHeader = make([]interface{}, len(cols))
// 		for i, col := range cols {
// 			switch val := col.(type) {
// 			case int:
// 				outHeader[i] = t.Header[val]
// 			case string:
// 				outHeader[i] = val
// 			}
// 		}
//
// 		_, vals := t._GenSliceLoc(1, cols...)
// 		t0 = &Table{
// 			Index:  t.Index,
// 			Header: outHeader,
// 			Vals:   vals}
// 	}
//
// 	if len(rows) > 0 {
// 		index := t.Index.Vals
// 		for i, in := range rows {
// 			switch val := in.(type) {
// 			case int:
// 				outIndex[i] = index[val]
// 			case string:
// 				outIndex[i] = val
// 			}
// 		}
// 		_, vals := t0._GenSliceLoc(0, rows...)
// 		t0 = &Table{
// 			Index:  NamedVector{t.Index.Header, outIndex},
// 			Header: outHeader,
// 			Vals:   SliceTranspose(vals)}
// 	}
//
// 	return t0
// }

// Index returns Index of a specific string element in a slice. Used for _SliceLoc to find Index position of name in Data.name field
func Index(search interface{}, slice []interface{}) (int, error) {
	for i, val := range slice {
		if search == val {
			return i, nil
		}
	}
	return -1, errors.New("Search not found")
}

// Concat concantenates multiple tables along a given axis and merges common rows and columns properly
func Concat(axis int, tables ...*Table) *Table {
	t := tables[0]

	if axis == 0 {
		var indexVals []interface{}
		for _, table := range tables { // get unique index values across all tables
			for _, val := range (*table).Index.Slice {
				if _, ok := Index(val, indexVals); ok != nil {
					indexVals = append(indexVals, val)
				}
			}
		}

		// axis = 0
		m := make(map[interface{}]interface{}, len(indexVals))
		for _, table := range tables {
			names, slices := table.PairedSliceLoc(0, indexVals...)
			for i := 0; i < len(names); i++ {
				if _, ok := m[names[i]]; ok == false {
					m[names[i]] = slices[i]
				} else {
					val := m[names[i]].([]interface{})
					m[names[i]] = append(val, slices[i]...)
				}
			}
		}

		var newHeaderVals []interface{}
		for _, table := range tables {
			for _, h := range table.Header.Slice {
				newHeaderVals = append(newHeaderVals, h)
			}
		}

		t = FromMap(0, m)
		t.Index.Header = tables[0].Index.Header
		t.Header = CreateMS(newHeaderVals, tables[0].Header.Header)
		t = t.GenSliceLoc(0, indexVals...) // needed to rearrange into original order after map disrupts order
	} else if axis == 1 {

		var headerVals []interface{}
		for _, table := range tables { // get unique index values across all tables
			for _, val := range (*table).Header.Slice {
				if _, ok := Index(val, headerVals); ok != nil {
					headerVals = append(headerVals, val)
				}
			}
		}

		// axis = 1
		m := make(map[interface{}]interface{}, len(headerVals))
		for _, table := range tables {
			names, slices := table.PairedSliceLoc(1, headerVals...)
			for i := 0; i < len(names); i++ {
				if _, ok := m[names[i]]; ok == false {
					m[names[i]] = slices[i]
				} else {
					val := m[names[i]].([]interface{})
					m[names[i]] = append(val, slices[i]...)
				}
			}
		}

		var newIndexVals []interface{}
		for _, table := range tables {
			for _, h := range table.Index.Slice {
				newIndexVals = append(newIndexVals, h)
			}
		}

		t = FromMap(1, m)
		t.Index = CreateMS(newIndexVals, tables[0].Index.Header)
		t = t.GenSliceLoc(1, headerVals...) // needed to rearrange into original order after map disrupts order
	}
	return t
}

// ToMap converts a Table to a map along a given axis (discards other axis)
func (t *Table) ToMap(axis int) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	if axis == 0 {
		index := t.Index.Slice
		for i := 0; i < len(index); i++ {
			m[index[i]] = t.Vals[i]
		}
		m["Header"] = t.Header
	} else if axis == 1 {
		header := t.Header.Slice
		for i := 0; i < len(header); i++ {
			m[header[i]] = GetTranspose(t.Vals, i)
		}
		m[t.Index.Header] = t.Index.Slice
	}
	return m
}

func main() {
	testFor := FromCSVFile("Data/test.csv", true, true)
	testFor.PrintTable()
}
