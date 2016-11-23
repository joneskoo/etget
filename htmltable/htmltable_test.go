package htmltable_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/joneskoo/etget/htmltable"
)

func TestParse(t *testing.T) {
	input := `
<html>
	<body>
		<table>
			<thead>
				<tr>
					<td colspan="20">ASDF</td>
				</tr><tr>
					<td colspan="20">FOO BAR</td>
				</tr><tr>
					<td></td>
					<td style="text-align:left;">H</td>
					<td style="text-align:center;">SYS</td>
					<td style="text-align:center;">S</td>
					<td style="text-align:center;">FOO</td>
				</tr>
			</thead><tbody>
				<tr>
					<td style="text-align:left;">01-01-2016</td>
					<td style="text-align:left;">00&nbsp;-&nbsp;01</td>
					<td style="text-align:right;">16,39</td>
				</tr><tr>
					<td style="text-align:left;">01-01-2016</td>
					<td style="text-align:left;">01&nbsp;-&nbsp;02</td>
					<td style="text-align:right;">16,04</td>
				</tr>
			</tbody>
		</table>
		<table>
			<thead>
				<tr><td>1</td></tr>
			</thead><tbody>
				<tr><td>2</td></tr>
			</tbody>
		</table>
	</body>
</html>
`
	want := []htmltable.Table{
		htmltable.Table{
			Headers: [][]string{
				[]string{"ASDF"},
				[]string{"FOO BAR"},
				[]string{"", "H", "SYS", "S", "FOO"},
			},
			Rows: [][]string{
				[]string{"01-01-2016", "00\u00A0-\u00A001", "16,39"},
				[]string{"01-01-2016", "01\u00A0-\u00A002", "16,04"},
			},
		},
		htmltable.Table{
			Headers: [][]string{
				[]string{"1"},
			},
			Rows: [][]string{
				[]string{"2"},
			},
		},
	}
	tables, err := htmltable.Parse(strings.NewReader(input))
	if err != nil {
		t.Errorf("htmltable.Parse(...): %s", err)
	}
	if !reflect.DeepEqual(tables, want) {
		t.Fatalf("htmltable.Parse(...) = %#v, want %#v", tables, want)
	}
}
