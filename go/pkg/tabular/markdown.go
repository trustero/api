package tabular

import (
	"fmt"
	"github.com/trustero/api/go/pkg/printer"
	"strings"
)

func NewMarkdownPrinter(table *Table) printer.Printer {
	return &tableMarkdownPrinter{
		table: table,
	}
}

type tableMarkdownPrinter struct {
	table *Table
}

func (p *tableMarkdownPrinter) Print() string {
	if len(p.table.Headers) == 0 {
		return ""
	}

	p.table.Prepare()
	isColumnCenterJustified := make([]bool, len(p.table.Headers))
	for i := range isColumnCenterJustified {
		isColumnCenterJustified[i] = true
	}
	var parts = []string{
		printRow(p.table.Headers),
		genColumnJustification(isColumnCenterJustified),
		printBody(p.table.Body),
	}
	return strings.Join(parts, "\n")
}

func printBody(body [][]string) string {
	if len(body) == 0 {
		return ""
	}
	bodyMarkdown := make([]string, len(body))
	for i, row := range body {
		bodyMarkdown[i] = printRow(row)
	}
	return strings.Join(bodyMarkdown, "\n")
}

func printRow(row []string) string {
	return fmt.Sprintf("| %s |", strings.Join(row, " | "))
}

func genColumnJustification(isColumnCenterJustified []bool) string {
	if len(isColumnCenterJustified) == 0 {
		return ""
	}
	var columnJustification []string
	for _, isCentered := range isColumnCenterJustified {
		if isCentered {
			columnJustification = append(columnJustification, ":---:")
			continue
		}
		columnJustification = append(columnJustification, ":---")
	}
	return fmt.Sprintf("|%s|", strings.Join(columnJustification, "|"))
}
