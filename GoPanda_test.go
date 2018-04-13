package main

import (
	"fmt"
	"testing"
)

var file1 string = "Data/test.csv"
var file2 string = "Data/test1.csv"

func TestCreate(t *testing.T) {
	// BOTH
	table1 := FromCSVFile(file1, true, true)
	table1.PrintTable()

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

	// HEADER
	table1 = FromCSVFile(file1, true, false)
	table1.PrintTable()

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

	// INDEX
	table1 = FromCSVFile(file1, false, true)
	table1.PrintTable()

	// +--------+-----+-------+
	// | INDEX  |  0  |   1   |
	// +--------+-----+-------+
	// | String | Int | Float |
	// | eff    |   1 |   4.2 |
	// | efe    |   3 |  5.32 |
	// | efe    |   2 |  1.32 |
	// | ffs    |  52 |   2.1 |
	// | wg     |  34 |    .8 |
	// | ret    |   4 |   9.6 |
	// +--------+-----+-------+

	// NONE
	table1 = FromCSVFile(file1, false, false)
	table1.PrintTable()

	// +-------+--------+-----+-------+
	// | INDEX |   0    |  1  |   2   |
	// +-------+--------+-----+-------+
	// |     0 | String | Int | Float |
	// |     1 | eff    |   1 |   4.2 |
	// |     2 | efe    |   3 |  5.32 |
	// |     3 | efe    |   2 |  1.32 |
	// |     4 | ffs    |  52 |   2.1 |
	// |     5 | wg     |  34 |    .8 |
	// |     6 | ret    |   4 |   9.6 |
	// +-------+--------+-----+-------+
}

func TestResetIndex(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	table1.ResetIndex()
	table1.PrintTable()

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
}

func TestSliceLoc0(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)

	test := table1.SliceLoc(0, "efe")
	test.PrintTable()

	// +--------+-----+-------+
	// | STRING | INT | FLOAT |
	// +--------+-----+-------+
	// | efe    |   3 |  5.32 |
	// | efe    |   2 |  1.32 |
	// +--------+-----+-------+

	test = table1.SliceLoc(1, "Float")
	test.PrintTable()

	// +--------+-------+
	// | STRING | FLOAT |
	// +--------+-------+
	// | eff    |   4.2 |
	// | efe    |  5.32 |
	// | efe    |  1.32 |
	// | ffs    |   2.1 |
	// | wg     |    .8 |
	// | ret    |   9.6 |
	// +--------+-------+
}

func TestSliceILoc0(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)

	test := table1.SliceILoc(0, 0)
	test.PrintTable()
	// +--------+-----+-------+
	// | STRING | INT | FLOAT |
	// +--------+-----+-------+
	// | eff    |   1 |   4.2 |
	// +--------+-----+-------+
}

func TestILoc(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)

	test := table1.ILoc([]int{0, 1}, []int{0, 1})
	test.PrintTable()
	// +--------+-----+-------+
	// | STRING | INT | FLOAT |
	// +--------+-----+-------+
	// | eff    |   1 |     3 |
	// | efe    |   3 |  5.32 |
	// +--------+-----+-------+
}

func TestTranspose(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	test1 := table1.Transpose()
	test1.PrintTable()

	// +---------+-----+------+------+-----+----+-----+
	// | COLUMNS | EFF | EFE  | EFE  | FFS | WG | RET |
	// +---------+-----+------+------+-----+----+-----+
	// | Int     |   1 |    3 |    2 |  52 | 34 |   4 |
	// | Float   | 4.2 | 5.32 | 1.32 | 2.1 | .8 | 9.6 |
	// +---------+-----+------+------+-----+----+-----+
}

func TestSetIndex(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	table1.SetIndex(0)
	table1.PrintTable()

	// +-----+--------+-------+
	// | INT | STRING | FLOAT |
	// +-----+--------+-------+
	// |   1 | eff    |   4.2 |
	// |   3 | efe    |  5.32 |
	// |   2 | efe    |  1.32 |
	// |  52 | ffs    |   2.1 |
	// |  34 | wg     |    .8 |
	// |   4 | ret    |   9.6 |
	// +-----+--------+-------+
}

func TestGenSliceLoc(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)

	test := table1.GenSliceLoc(1, 0)
	test.PrintTable()

	// +--------+-----+
	// | STRING | INT |
	// +--------+-----+
	// | eff    |   1 |
	// | efe    |   3 |
	// | efe    |   2 |
	// | ffs    |  52 |
	// | wg     |  34 |
	// | ret    |   4 |
	// +--------+-----+
}

func TestFromMap(t *testing.T) {
	m := make(map[interface{}]interface{})
	m["Test1"] = []interface{}{"1Test", "2Test", "3Test"}
	m["Test2"] = []interface{}{"4Test", "5Test", "6Test"}
	m["Test3"] = []interface{}{"7Test", "8Test", "9Test"}

	tableOut := FromMap(1, m)
	tableOut.PrintTable()

	// +-------+-------+-------+-------+
	// | INDEX | TEST1 | TEST2 | TEST3 |
	// +-------+-------+-------+-------+
	// |     0 | 1Test | 4Test | 7Test |
	// |     1 | 2Test | 5Test | 8Test |
	// |     2 | 3Test | 6Test | 9Test |
	// +-------+-------+-------+-------+
}

func TestAddSlice(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	table1.AddSlice(0, "check", []interface{}{0, 1})
	table1.PrintTable()

	// +--------+-----+-------+
	// | STRING | INT | FLOAT |
	// +--------+-----+-------+
	// | eff    |   1 |   4.2 |
	// | efe    |   3 |  5.32 |
	// | efe    |   2 |  1.32 |
	// | ffs    |  52 |   2.1 |
	// | wg     |  34 |    .8 |
	// | ret    |   4 |   9.6 |
	// | check  |   0 |     1 |
	// +--------+-----+-------+
}

func TestPairedSliceLoc(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)

	fmt.Println(table1.PairedSliceLoc(0, "efe", "trer", "hello"))
	// [efe efe trer hello] [[3 5.32] [2 1.32] [<nil> <nil>] [<nil> <nil>]]
}

func TestUniqs(t *testing.T) {
	testSlice := make([]interface{}, 2)
	for i := range testSlice {
		testSlice[i] = 1
	}
	testSlice = getuniqs(testSlice)
	fmt.Println("Uniqs", testSlice)
	// Uniqs [1]
}

func TestConcat(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	table2 := FromCSVFile(file2, true, true)

	table2.PrintTable()

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

	tableOut := Concat(0, table1, table2)
	tableOut.PrintTable()

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

	tableOut = Concat(1, table1, table2)
	tableOut.PrintTable()

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

}

func TestToMap(t *testing.T) {
	table1 := FromCSVFile(file1, true, true)
	fmt.Println(table1.ToMap(0))
	// map[wg:[34 .8] ret:[4 9.6] Columns:[Int Float] eff:[1 4.2] efe:[2 1.32] ffs:[52 2.1]]
	fmt.Println(table1.ToMap(1))
	// map[Int:[1 3 2 52 34 4] Float:[4.2 5.32 1.32 2.1 .8 9.6] String:[eff efe efe ffs wg ret]]
}
