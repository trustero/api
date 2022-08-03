// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package tabular_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustero/api/go/pkg/tabular"
	"testing"
)

func TestEmptyTable(t *testing.T) {
	table := &tabular.Table{
		Headers: make([]string, 0),
		Body:    make([][]string, 0),
	}
	assert.Equal(t, "", tabular.NewMarkdownPrinter(table).Print())
}

func TestOnlyOneHeaderNoBody(t *testing.T) {
	table := &tabular.Table{
		Headers: []string{"foo"},
		Body:    make([][]string, 0),
	}

	expected := `| foo |
|:---:|
`
	printer := tabular.NewMarkdownPrinter(table)
	assert.Equal(t, expected, printer.Print())
}
func TestOnlyNHeaderNoBody(t *testing.T) {
	table := &tabular.Table{
		Headers: []string{"foo", "bar", "baz"},
		Body:    make([][]string, 0),
	}

	const expected = `| foo | bar | baz |
|:---:|:---:|:---:|
`
	assert.Equal(t, expected, tabular.NewMarkdownPrinter(table).Print())
}

func TestOnlyOneHeaderOneRow(t *testing.T) {
	table := &tabular.Table{
		Headers: []string{"foo"},
		Body: [][]string{
			{"foo-a"},
		},
	}
	expected := `| foo |
|:---:|
| foo-a |`
	assert.Equal(t, expected, tabular.NewMarkdownPrinter(table).Print())
}

func TestNHeadersNRow(t *testing.T) {
	table := &tabular.Table{
		Headers: []string{"foo", "bar", "baz"},
		Body: [][]string{
			{"foo-a", "bar-a", "baz-a"},
			{"foo-b", "-", "-"},
			{"foo-c", "bar-c", "-"},
		},
	}

	expected := `| foo | bar | baz |
|:---:|:---:|:---:|
| foo-a | bar-a | baz-a |
| foo-b | - | - |
| foo-c | bar-c | - |`
	assert.Equal(t, expected, tabular.NewMarkdownPrinter(table).Print())
}
