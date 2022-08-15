// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package receptor_v1 provides the Go GRPC client bindings to communicate with the Trustero service.
package receptor_v1

import (
	"strconv"
	"strings"
)

const dateTimeLayout = "2006-Jan-02"

// Tabulate converts a receptor_v1.Struct to an ordered and displayable array header strings, and an array of
// rows of strings.  Each row's columns are ordered according to its headers in the headers array.
func (s *Struct) Tabulate() (headers []string, rows [][]string, err error) {
	headerKeys := s.getHeaderKeys()

	for _, row := range s.Rows {
		cols := make([]string, len(headerKeys))
		for i, key := range headerKeys {
			if len(key) == 0 {
				cols[i] = ""
				continue
			}
			cols[i] = toStringValue(row.Cols[key])
		}
		rows = append(rows, cols)
	}

	// Get displayable headers
	for _, h := range headerKeys {
		name := h
		if v, ok := s.ColDisplayNames[h]; ok {
			name = v
		}
		headers = append(headers, name)
	}

	return
}

func (s *Struct) getHeaderKeys() (displayHeaders []string) {
	if len(s.Rows) == 0 {
		return
	}

	// Build displayHeader array of column key names
	var allHeaders []string
	for k, _ := range s.Rows[0].Cols {
		allHeaders = append(allHeaders, k)
	}

	displayHeaders = s.ColDisplayOrder
	for _, h := range allHeaders {
		// If displayHeader doesn't contain a known column key, add the column key
		if !contains(s.ColDisplayOrder, h) {
			displayHeaders = append(displayHeaders, h)
		}
	}

	return
}

func toStringValue(value *Value) (str string) {
	str = ""
	v := value.GetValueType()
	if dv, ok := v.(*Value_DoubleValue); ok {
		str = strconv.FormatFloat(dv.DoubleValue, 'g', -1, 64)

	} else if fv, ok := v.(*Value_FloatValue); ok {
		str = strconv.FormatFloat(float64(fv.FloatValue), 'g', -1, 32)

	} else if i, ok := v.(*Value_Int32Value); ok {
		str = strconv.FormatInt(int64(i.Int32Value), 10)

	} else if ii, ok := v.(*Value_Int64Value); ok {
		str = strconv.FormatInt(ii.Int64Value, 10)

	} else if ui, ok := v.(*Value_Uint32Value); ok {
		str = strconv.FormatUint(uint64(ui.Uint32Value), 10)

	} else if uii, ok := v.(*Value_Uint64Value); ok {
		str = strconv.FormatUint(uii.Uint64Value, 10)

	} else if b, ok := v.(*Value_BoolValue); ok {
		if b.BoolValue {
			str = ":heavy_check_mark:"
		} else {
			str = "-"
		}

	} else if s, ok := v.(*Value_StringValue); ok {
		str = strings.TrimSpace(s.StringValue)

	} else if t, ok := v.(*Value_TimestampValue); ok {
		dt := t.TimestampValue
		if dt == nil || (dt.Nanos == 0 && dt.Seconds == 0) {
			str = "-"
		} else {
			str = dt.AsTime().Format(dateTimeLayout)
		}
	}
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
