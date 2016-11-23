package tableparser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/joneskoo/etget/tableparser"
)

func TestHeader(t *testing.T) {
	want := [][]string{
		[]string{"Elspot Prices in EUR/MWh"},
		[]string{"Data was last updated 21-11-2016"},
		[]string{"", "Hours", "SYS", "SE1", "SE2", "SE3", "SE4", "FI", "DK1", "DK2", "Oslo", "Kr.sand", "Bergen", "Molde", "Tr.heim", "Tromsø", "EE", "LV", "LT", "FRE"},
	}
	rows, err := tableparser.Header(strings.NewReader(testHTML))
	if err != nil {
		t.Errorf("tableparser.Header(...): %s", err)
	}
	if len(rows) != len(want) {
		t.Errorf("len(rows) = %d, want %d", len(rows), len(want))
		t.FailNow()
	}
	for i, row := range want {
		if !reflect.DeepEqual(row, rows[i]) {
			t.Errorf("rows[%d] = %v, want %v", i, row, want[i])
		}
	}
}

func TestData(t *testing.T) {
	want := [][]string{
		[]string{"01-01-2016", "00\u00A0-\u00A001", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "28,11", "28,11", ""},
		[]string{"01-01-2016", "01\u00A0-\u00A002", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", ""},
	}
	rows, err := tableparser.Data(strings.NewReader(testHTML))
	if err != nil {
		t.Errorf("tableparser.Data(...): %s", err)
	}
	if len(rows) != len(want) {
		t.Errorf("len(rows) = %d, want %d", len(rows), len(want))
		t.FailNow()
	}
	for i, row := range want {
		if !reflect.DeepEqual(row, rows[i]) {
			t.Errorf("rows[%d] = %v, want %v", i, row, want[i])
		}
	}
}

const testHTML = `
<html>
	<body>
		<table>
			<thead>
				<tr>
					<td colspan="20">Elspot Prices in EUR/MWh</td>
				</tr><tr>
					<td colspan="20">Data was last updated 21-11-2016</td>
				</tr><tr>
					<td></td><td style="text-align:left;">Hours</td><td style="text-align:center;">SYS</td><td style="text-align:center;">SE1</td><td style="text-align:center;">SE2</td><td style="text-align:center;">SE3</td><td style="text-align:center;">SE4</td><td style="text-align:center;">FI</td><td style="text-align:center;">DK1</td><td style="text-align:center;">DK2</td><td style="text-align:center;">Oslo</td><td style="text-align:center;">Kr.sand</td><td style="text-align:center;">Bergen</td><td style="text-align:center;">Molde</td><td style="text-align:center;">Tr.heim</td><td style="text-align:center;">Tromsø</td><td style="text-align:center;">EE</td><td style="text-align:center;">LV</td><td style="text-align:center;">LT</td><td style="text-align:center;">FRE</td>
				</tr>
			</thead><tbody>
				<tr>
					<td style="text-align:left;">01-01-2016</td><td style="text-align:left;">00&nbsp;-&nbsp;01</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">28,11</td><td style="text-align:right;">28,11</td><td style="text-align:right;"></td>
				</tr><tr>
					<td style="text-align:left;">01-01-2016</td><td style="text-align:left;">01&nbsp;-&nbsp;02</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;"></td>
				</tr>
			</tbody>
		</table>
	</body>
</html>
`

func TestTwoTablesData(t *testing.T) {
	// If there are many tables, parser should return the first table in
	// the document and ignore everything else.
	want := [][]string{
		[]string{"01-01-2016", "00\u00A0-\u00A001", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "16,39", "28,11", "28,11", ""},
		[]string{"01-01-2016", "01\u00A0-\u00A002", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", "16,04", ""},
	}
	rows, err := tableparser.Data(strings.NewReader(testHTMLTwoTables))
	if err != nil {
		t.Errorf("tableparser.Data(...): %s", err)
	}
	if len(rows) != len(want) {
		t.Errorf("len(rows) = %d, want %d", len(rows), len(want))
		t.FailNow()
	}
	for i, row := range want {
		if !reflect.DeepEqual(row, rows[i]) {
			t.Errorf("rows[%d] = %v, want %v", i, row, want[i])
		}
	}
}

const testHTMLTwoTables = `
<html>
	<body>
		<table>
			<thead>
				<tr>
					<td colspan="20">Elspot Prices in EUR/MWh</td>
				</tr><tr>
					<td colspan="20">Data was last updated 21-11-2016</td>
				</tr><tr>
					<td></td><td style="text-align:left;">Hours</td><td style="text-align:center;">SYS</td><td style="text-align:center;">SE1</td><td style="text-align:center;">SE2</td><td style="text-align:center;">SE3</td><td style="text-align:center;">SE4</td><td style="text-align:center;">FI</td><td style="text-align:center;">DK1</td><td style="text-align:center;">DK2</td><td style="text-align:center;">Oslo</td><td style="text-align:center;">Kr.sand</td><td style="text-align:center;">Bergen</td><td style="text-align:center;">Molde</td><td style="text-align:center;">Tr.heim</td><td style="text-align:center;">Tromsø</td><td style="text-align:center;">EE</td><td style="text-align:center;">LV</td><td style="text-align:center;">LT</td><td style="text-align:center;">FRE</td>
				</tr>
			</thead><tbody>
				<tr>
					<td style="text-align:left;">01-01-2016</td><td style="text-align:left;">00&nbsp;-&nbsp;01</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">28,11</td><td style="text-align:right;">28,11</td><td style="text-align:right;"></td>
				</tr><tr>
					<td style="text-align:left;">01-01-2016</td><td style="text-align:left;">01&nbsp;-&nbsp;02</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;">16,04</td><td style="text-align:right;"></td>
				</tr>
			</tbody>
		</table>
		<table>
			<thead>
				<tr>
					<td colspan="20">This is ignored</td>
				</tr><tr>
					<td colspan="20">Data was last updated 21-11-2016</td>
				</tr><tr>
					<td></td><td style="text-align:left;">Hours</td><td style="text-align:center;">SYS</td><td style="text-align:center;">SE1</td><td style="text-align:center;">SE2</td><td style="text-align:center;">SE3</td><td style="text-align:center;">SE4</td><td style="text-align:center;">FI</td><td style="text-align:center;">DK1</td><td style="text-align:center;">DK2</td><td style="text-align:center;">Oslo</td><td style="text-align:center;">Kr.sand</td><td style="text-align:center;">Bergen</td><td style="text-align:center;">Molde</td><td style="text-align:center;">Tr.heim</td><td style="text-align:center;">Tromsø</td><td style="text-align:center;">EE</td><td style="text-align:center;">LV</td><td style="text-align:center;">LT</td><td style="text-align:center;">FRE</td>
				</tr>
			</thead><tbody>
				<tr>
					<td style="text-align:left;">01-01-2016</td><td style="text-align:left;">00&nbsp;-&nbsp;01</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">16,39</td><td style="text-align:right;">28,11</td><td style="text-align:right;">28,11</td><td style="text-align:right;"></td>
				</tr>
			</tbody>
		</table>
	</body>
</html>
`
