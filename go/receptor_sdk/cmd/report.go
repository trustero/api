package cmd

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func report(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

	// Discover evidence
	var discovered []*receptor_sdk.Evidence
	if discovered, err = receptorImpl.Report(credentials); err != nil {
		return
	}

	// Report discovered evidence to Trustero
	var finding receptor_v1.Finding

	finding.ReceptorType = receptorImpl.GetReceptorType()
	finding.ServiceProviderAccount = serviceProviderAccount

	// Convert and append discovered evidences to reported evidences
	for _, evidence := range discovered {
		reportStruct := receptor_v1.Struct{
			Rows:            []*receptor_v1.Struct_Row{},
			ColDisplayNames: map[string]string{},
			ColDisplayOrder: []string{},
		}

		reportEvidence := receptor_v1.Evidence{
			Caption:      evidence.Caption,
			Description:  evidence.Description,
			ServiceName:  evidence.ServiceName,
			Sources:      []*receptor_v1.Evidence_Source{},
			EvidenceType: &receptor_v1.Evidence_Struct{Struct: &reportStruct},
		}

		// Convert sources
		for _, source := range evidence.Sources {
			reportEvidence.Sources = append(reportEvidence.Sources, &receptor_v1.Evidence_Source{
				RawApiRequest:  source.ProviderAPIRequest,
				RawApiResponse: source.ProviderAPIResponse,
			})
		}

		// Convert rows
		var serviceIdFieldName string
		var rowFieldNames []string
		for idx, row := range evidence.Rows {
			if idx == 0 {
				if serviceIdFieldName, rowFieldNames, err = extractMetaData(row, &reportStruct); err != nil {
					return // fail to extract metadata, likely an invalid row type
				}
			}
			reportStruct.Rows = append(reportStruct.Rows, rowToStructRow(row, serviceIdFieldName, rowFieldNames))
		}

		// Append to Finding
		finding.Evidences = append(finding.Evidences, &reportEvidence)
	}

	// Report evidence findings to Trustero
	_, err = rc.Report(context.Background(), &finding)

	return
}

func extractMetaData(row interface{}, reportStruct *receptor_v1.Struct) (serviceIdFieldName string, rowFieldNames []string, err error) {
	rowFieldNames = []string{}
	rowType := reflect.TypeOf(row)
	if err = assertStruct(rowType); err != nil {
		return
	}

	fieldOrder := map[int]string{}
	fieldOrderKeys := []int{}
	for i := 0; i < rowType.NumField(); i++ {
		field := rowType.Field(i)
		tags := expandFieldTag(field)
		rowFieldNames = append(rowFieldNames, field.Name)

		// Is it the id field?
		if _, ok := tags["id"]; ok {
			serviceIdFieldName = field.Name
		}

		// Get the field order
		if val, ok := tags["order"]; ok {
			if i, err := strconv.Atoi(val); err == nil {
				fieldOrder[i] = field.Name
				fieldOrderKeys = append(fieldOrderKeys, i)
			}
		}

		// Get display name
		if val, ok := tags["display"]; ok {
			reportStruct.ColDisplayNames[field.Name] = val
		} else {
			reportStruct.ColDisplayNames[field.Name] = field.Name
		}
	}

	// order the display columns
	sort.Ints(fieldOrderKeys)
	for _, key := range fieldOrderKeys {
		reportStruct.ColDisplayOrder = append(reportStruct.ColDisplayOrder, fieldOrder[key])
	}
	return
}

func rowToStructRow(row interface{}, serviceIdFieldName string, rowFieldNames []string) (reportRow *receptor_v1.Struct_Row) {
	reportRow = &receptor_v1.Struct_Row{
		ServiceId: getField(row, serviceIdFieldName),
		Cols:      map[string]*receptor_v1.Struct_Row_Value{},
	}

	rowValue := reflect.Indirect(reflect.ValueOf(row))
	for _, fieldName := range rowFieldNames {
		v := rowValue.FieldByName(fieldName)

		if dateTime, ok := v.Interface().(time.Time); ok {
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_TimestampValue{
					TimestampValue: timestamppb.New(dateTime),
				},
			}
			continue
		}

		if dateTime, ok := v.Interface().(*time.Time); ok && dateTime != nil {
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_TimestampValue{
					TimestampValue: timestamppb.New(*dateTime),
				},
			}
			continue
		}

		switch v.Kind() {
		case reflect.Bool:
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_BoolValue{
					BoolValue: v.Bool(),
				},
			}
			break
		case reflect.Int:
		case reflect.Int8:
		case reflect.Int16:
		case reflect.Int32:
		case reflect.Int64:
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_Int64Value{
					Int64Value: v.Int(),
				},
			}
			break
		case reflect.Uint:
		case reflect.Uint8:
		case reflect.Uint16:
		case reflect.Uint32:
		case reflect.Uint64:
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_Uint64Value{
					Uint64Value: v.Uint(),
				},
			}
			break
		case reflect.Float32:
		case reflect.Float64:
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_DoubleValue{
					DoubleValue: v.Float(),
				},
			}
			break
		case reflect.String:
			reportRow.Cols[fieldName] = &receptor_v1.Struct_Row_Value{
				ValueType: &receptor_v1.Struct_Row_Value_StringValue{
					StringValue: v.String(),
				},
			}
			break
		case reflect.Complex64:
		case reflect.Complex128:
		case reflect.Array:
		case reflect.Chan:
		case reflect.Func:
		case reflect.Interface:
		case reflect.Map:
		case reflect.Pointer:
		case reflect.Slice:
		case reflect.Struct:
		case reflect.UnsafePointer:
		case reflect.Uintptr:
		default:
			log.Warn().Msg("unsupported evidence row field (" + fieldName + ") type")
			break
		}
	}

	return
}

func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func assertStruct(rowType reflect.Type) (err error) {
	if rowType.Kind() != reflect.Struct {
		err = errors.New("evidence row must be a struct. " + rowType.String())
	}
	return
}
