package main

import (
	"fmt"
	"testing"
)

// import from slice, map?
// export to slice?, map?
// concat?

func TestCreateBoth(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)
	testFor.PrintTable()
}

func TestCreateHeader(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, false)
	testFor.PrintTable()
}

func TestCreateIndex(t *testing.T) {
	testFor := FromCSVFile("test.csv", false, true)
	testFor.PrintTable()
}

func TestCreateNone(t *testing.T) {
	testFor := FromCSVFile("test.csv", false, false)
	testFor.PrintTable()
}

func TestResetIndex(t *testing.T) {
	table1 := FromCSVFile("test.csv", true, true)
	table1.ResetIndex()
	table1.PrintTable()
}

func TestResetHeader(t *testing.T) {
	table1 := FromCSVFile("test.csv", true, true)
	table1.ResetHeader()
	table1.PrintTable()
}

func TestSliceLoc0(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	test := testFor.SliceLoc(0, "efe")
	test.PrintTable()
}

func TestSliceLoc1(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	test := testFor.SliceLoc(1, "Float")
	test.PrintTable()
}

func TestSliceILoc0(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	test := testFor.SliceILoc(0, 0)
	test.PrintTable()
}

func TestILoc(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	test := testFor.ILoc([]int{0, 1}, []int{0, 1})
	test.PrintTable()
}

func TestTranspose(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)
	test1 := testFor.Transpose()
	test1.PrintTable()
}

func TestSetIndex(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)
	testFor.SetIndex(0)
	testFor.PrintTable()
}

func TestGenSliceLoc(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	test := testFor.GenSliceLoc(1, 0)
	test.PrintTable()
}

func TestFromMap(t *testing.T) {
	m := make(map[interface{}]interface{})
	m["Test1"] = []interface{}{"1Test", "2Test", "3Test"}
	m["Test2"] = []interface{}{"4Test", "5Test", "6Test"}
	m["Test3"] = []interface{}{"7Test", "8Test", "9Test"}

	tableOut := FromMap(1, m)
	tableOut.PrintTable()
}

func TestAddSlice(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)
	testFor.AddSlice(0, "vin", []interface{}{0, 1})
	testFor.PrintTable()
}

func TestPairedSliceLoc(t *testing.T) {
	testFor := FromCSVFile("test.csv", true, true)

	fmt.Println(testFor.PairedSliceLoc(0, "efe", "trer", "hello"))
}

func TestConcat0(t *testing.T) {
	table1 := FromCSVFile("test.csv", true, true)
	table2 := FromCSVFile("test1.csv", true, true)

	Concat(0, table1, table2).PrintTable()
}

func TestConcat1(t *testing.T) {
	table1 := FromCSVFile("test.csv", true, true)
	table2 := FromCSVFile("test1.csv", true, true)

	Concat(1, table1, table2).PrintTable()
}
