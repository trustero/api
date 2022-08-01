package receptor_v1

import (
	"fmt"
	"github.com/trustero/api/go/pkg/printer"
	. "github.com/trustero/api/go/pkg/proto-utils"
	. "github.com/trustero/api/go/pkg/reflect-utils"
	"github.com/trustero/api/go/pkg/tabular"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"time"
)

func NewStruct(values []interface{}) (strct *Struct, err error) {
	if len(values) == 0 {
		return
	}

	var trTag *TrTag
	trTag, err = NewTrTag(values[0])

	strct = &Struct{
		Rows:            make([]*Struct_Row, len(values)),
		ColDisplayNames: trTag.DisplayName,
		ColDisplayOrder: trTag.DisplayOrder,
	}
	for i, value := range values {
		strct.Rows[i] = newStruct_Row(value, trTag)
	}
	return
}

func newStruct_Row(value interface{}, trTag *TrTag) (row *Struct_Row) {
	reflectedValue, reflectedField := GetValueAndType(value)
	row = &Struct_Row{
		ServiceId: trTag.Id,
		Cols:      make(map[string]*Struct_Row_Value),
	}

	for i := 0; i < reflectedValue.NumField(); i++ {
		field := reflectedField.Field(i)
		fieldValue := reflectedValue.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_StringValue{
				StringValue: fieldValue.String(),
			}}
			continue
		case reflect.Bool:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_BoolValue{
				BoolValue: fieldValue.Bool(),
			}}
			continue
		case reflect.Int32:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_Int32Value{
				Int32Value: int32(fieldValue.Int()),
			}}
			continue
		case reflect.Int, reflect.Int64:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_Int64Value{
				Int64Value: fieldValue.Int(),
			}}
			continue
		case reflect.Uint32:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_Uint32Value{
				Uint32Value: uint32(fieldValue.Uint()),
			}}
			continue
		case reflect.Uint, reflect.Uint64:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_Uint64Value{
				Uint64Value: fieldValue.Uint(),
			}}
			continue
		case reflect.Float32, reflect.Float64:
			row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_FloatValue{
				FloatValue: float32(fieldValue.Float()),
			}}
			continue
		default:
			if dateTime, ok := fieldValue.Interface().(time.Time); ok {
				row.Cols[field.Name] = &Struct_Row_Value{ValueType: &Struct_Row_Value_TimestampValue{
					TimestampValue: timestamppb.New(dateTime),
				}}
				continue
			}
		}
		panic(fmt.Sprintf("unsupported type of field %s", field.Name))
	}
	return
}

// Tabulate returns the collection of elements wrapped by Struct as a 2D array, where each row is a Struct_Row and each column is an attribute.
// The order and display names of the attributes are given by the struct's `display_order` and `display_name` fields respectively.
func (x *Struct) Tabulate(strconv printer.Strconv) (result *tabular.Table, err error) {
	result = &tabular.Table{}

	for _, row := range x.GetRows() {
		attributes := make([]string, len(x.GetColDisplayOrder()))
		for i, key := range x.GetColDisplayOrder() {
			if len(key) == 0 {
				attributes[i] = ""
				continue
			}
			var value []string
			if value, err = TabulateField(row.GetCols()[key], 0, strconv); err != nil {
				continue
			}
			attributes[i] = strings.Join(value, "")
			result.Headers = append(result.Headers, x.GetColDisplayNames()[key])
		}
		result.Body = append(result.Body, attributes)
	}
	return
}
