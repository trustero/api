package cmd_test

import (
	"testing"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
)

// A control's custom field values travel as a StructList column named after the
// field, one struct per value with its Name and Value, which is what Trustero
// decodes on the other side.
func TestRowToStructRowSerializesControlCustomFields(t *testing.T) {
	row := receptor_sdk.Control{
		Id:   "c1",
		Name: "Access Control",
		CustomFields: []receptor_sdk.CustomFieldValue{
			{Name: "Implementation", Value: "Quarterly review"},
		},
	}
	reportStruct := &receptor_v1.Struct{ColDisplayNames: map[string]string{}, ColTags: map[string]string{}}
	entityIDField, fieldNames, err := cmd.ExtractMetaData(row, reportStruct)
	if err != nil {
		t.Fatal(err)
	}
	got := cmd.RowToStructRow(row, entityIDField, fieldNames)

	column, ok := got.Cols["CustomFields"]
	if !ok {
		t.Fatalf("no CustomFields column in %v", got.Cols)
	}
	list, ok := column.ValueType.(*receptor_v1.Value_StructListValue)
	if !ok {
		t.Fatalf("CustomFields is %T, want a StructList", column.ValueType)
	}
	if len(list.StructListValue.Values) != 1 {
		t.Fatalf("want one value, got %v", list.StructListValue.Values)
	}
	fields := list.StructListValue.Values[0].Fields
	if fields["Name"].GetStringValue() != "Implementation" || fields["Value"].GetStringValue() != "Quarterly review" {
		t.Fatalf("unexpected fields %v", fields)
	}
}
