GoTables
===============================

version number: 0.0.1
author: Vincent Nikolayev

Overview
--------

A Go library enabling operations on table-like structures. Influenced by the Python Pandas Package.

Requires
-------
* github.com/olekukonko/tablewriter

Current Features:
-----------------
* Select columns and/or rows using names
  * Fast name lookups using map like-structures
* Concatenate multiple tables together
* Allowable duplicate keys for index and header names
* Built using interfaces
* Table transposition
* Graphic printability using ascii tables
* Table creation from .csv, slices, and maps
* Table writing to maps

To Do:
-----------------
* Implement to slices
* Clean up and document code
* Others
  * If I or someone else has ideas to implement

Example:
--------
```Go
table1 := FromCSVFile("Data/test.csv", true, true)
table1.PrintTable()

//+--------+-----+-------+
// | STRING | INT | FLOAT |
// +--------+-----+-------+
// | eff    |   1 |   4.2 |
// | efe    |   3 |  5.32 |
// | efe    |   2 |  1.32 |
// | ffs    |  52 |   2.1 |
// | wg     |  34 |    .8 |
// | ret    |   4 |   9.6 |
// +--------+-----+-------+

table2 := FromCSVFile("Data/test1.csv", true, true)
table2.PrintTable()// +--------+-----+-------+

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

tableOut := Concat(axis, table1, table2)
tableOut.PrintTable()

// axis = 0:
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
// axis = 1:
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
