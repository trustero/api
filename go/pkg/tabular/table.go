package tabular

// BlankHeaderSentinel This is used as a placeholder value for when we want to insert an actual blank header value.
var BlankHeaderSentinel = "TRUSTERO_BLANK_HEADER_SENTINEL"

type Table struct {
	Headers     []string
	Body        [][]string
	ProblemRows map[int]bool
}

func (t *Table) Prepare() {
	t.annotatedTable()

	t.excludeIgnored()
}

func (t *Table) excludeIgnored() {
	ignoredIndexes := map[int]bool{}
	j := 0
	for i, h := range t.Headers {
		if h == "" {
			ignoredIndexes[i] = true
			continue
		} else if h == BlankHeaderSentinel {
			t.Headers[j] = ""
		} else {
			t.Headers[j] = h
		}
		j++
	}
	t.Headers = t.Headers[:j]

	for i, row := range t.Body {
		j = 0
		for k, v := range row {
			if _, ok := ignoredIndexes[k]; ok {
				continue
			}
			row[j] = v
			j++
		}
		t.Body[i] = row[:j]
	}
}

func (t *Table) annotatedTable() {
	hasProblem := false
	for _, hasProblem = range t.ProblemRows {
		if hasProblem {
			break
		}
	}
	if !hasProblem {
		return
	}

	if t.Headers != nil && len(t.Headers) > 0 {
		// Use a known sentinel value to insert a blank header value. Otherwise, blank headers will be stripped from
		// the markdown output.
		t.Headers = append([]string{BlankHeaderSentinel}, t.Headers...)
	}

	for i, row := range t.Body {
		// If there are issues, we always add a new column to the start of the row.
		//   - If there is an issue for a row, add an icon to the new column and highlight the row.
		//   - If there is not an issue for a row, add an empty string to the new column.
		if t.ProblemRows[i] {
			for j, cellContent := range row {
				row[j] = "***" + cellContent + "***"
			}
			row = append([]string{"❌"}, row...)
		} else {
			row = append([]string{""}, row...)
		}

		(t.Body)[i] = row
	}
}
