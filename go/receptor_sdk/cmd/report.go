// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/multipartkit"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func report(rc receptor_v1.ReceptorClient, credentials interface{}, config interface{}) (err error) {

	// Report discovered evidence to Trustero
	var finding receptor_v1.Finding

	// Discover service entities
	if finding.Entities, err = receptorImpl.Discover(credentials, config); err != nil {
		return
	}
	finding.ReceptorType = GetParsedReceptorType()
	finding.ServiceProviderAccount = serviceProviderAccount

	// report in single batch
	var evidences []*receptor_sdk.Evidence
	if evidences, err = receptorImpl.Report(credentials, config); err == nil && len(evidences) > 0 {
		_ = reportEvidence(rc, &finding, evidences)
	}

	// report in multiple batches
	evidenceChannel := make(chan []*receptor_sdk.Evidence)
	go func() {
		defer close(evidenceChannel)
		receptorImpl.ReportBatch(credentials, evidenceChannel)
	}()

	for evidences := range evidenceChannel {
		// Receive evidence and report them one batch at a time
		err = reportEvidence(rc, &finding, evidences)
		if err != nil {
			log.Err(err).Msg("failed to report evidence")
			//return err
			continue
		}
		finding.Evidences = []*receptor_v1.Evidence{} // Empty every time for new evidence
	}

	return
}

func reportEvidence(rc receptor_v1.ReceptorClient, finding *receptor_v1.Finding, evidences []*receptor_sdk.Evidence) (err error) {
	for _, evidence := range evidences {
		reportStruct := receptor_v1.Struct{
			Rows:            []*receptor_v1.Row{},
			ColDisplayNames: map[string]string{},
			ColDisplayOrder: []string{},
			ColTags:         map[string]string{},
		}

		reportEvidence := receptor_v1.Evidence{
			Caption:          evidence.Caption,
			Description:      evidence.Description,
			ServiceName:      evidence.ServiceName,
			EntityType:       evidence.EntityType,
			Sources:          evidence.Sources,
			ServiceAccountId: evidence.ServiceAccountId,
		}

		if evidence.Document != nil { // evidence is a blob or path to blob
			// create a new finding from current finding and add evidence
			reportEvidence.EvidenceType = &receptor_v1.Evidence_Doc{
				Doc: &receptor_v1.Document{
					Body:           evidence.Document.Body,
					Mime:           evidence.Document.Mime,
					StreamFilePath: evidence.Document.StreamFilePath,
				},
			}
			reportFinding := receptor_v1.Finding{
				ReceptorType:           finding.ReceptorType,
				ServiceProviderAccount: finding.ServiceProviderAccount,
				Entities:               finding.Entities,
				Evidences:              []*receptor_v1.Evidence{&reportEvidence},
			}
			contentType, streamFile, err := multipartEvidence(&reportFinding)
			os.Remove(evidence.Document.StreamFilePath)
			if err != nil {
				log.Err(err).Msg("failed to create multipart evidence")
				err = nil
				continue
			}

			// make a multipart file and then stream it

			stream, err := rc.StreamReport(context.Background())
			defer stream.CloseSend()
			if err != nil {
				log.Err(err).Msg("failed to stream report")
				continue
			}

			//send boundary of the multipart first
			if err = stream.Send(&receptor_v1.ReportChunk{Content: []byte(contentType), IsBoundary: true}); err != nil {
				log.Err(err).Msg("failed to send data chunk")
				break
			}

			//read from the file path and stream in chunks
			file, err := os.Open(streamFile)
			defer func() {
				file.Close()
				os.Remove(streamFile)
			}()

			if err != nil {
				log.Err(err).Msg("failed to open file")
				continue
			}
			buf := make([]byte, 1024)
			for {
				n, err := file.Read(buf)
				if err != nil {
					break
				}
				if err = stream.Send(&receptor_v1.ReportChunk{Content: buf[:n]}); err != nil {
					log.Err(err).Msg("failed to send data chunk")
					break
				}
			}
		} else { // evidence is structured
			reportEvidence.EvidenceType = &receptor_v1.Evidence_Struct{Struct: &reportStruct}

			// Convert rows
			var entityIdFieldName string
			var rowFieldNames []string
			for idx, row := range evidence.Rows {
				if idx == 0 {
					if entityIdFieldName, rowFieldNames, err = ExtractMetaData(row, &reportStruct); err != nil {
						return // fail to extract metadata, likely an invalid row type
					}
				}
				reportStruct.Rows = append(reportStruct.Rows, RowToStructRow(row, entityIdFieldName, rowFieldNames))
			}

			// Append to Finding
			finding.Evidences = append(finding.Evidences, &reportEvidence)
		}

		// Report evidence findings to Trustero
		_, err = rc.Report(context.Background(), finding)
	}
	return

}

// ExtractMetaData Extracts tag information from struct
func ExtractMetaData(row interface{}, reportStruct *receptor_v1.Struct) (entityIdFieldName string, rowFieldNames []string, err error) {
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
		if _, ok := tags[idField]; ok {
			entityIdFieldName = field.Name
		}

		// Get the field order
		if val, ok := tags[orderField]; ok {
			if i, err := strconv.Atoi(val); err == nil {
				fieldOrder[i] = field.Name
				fieldOrderKeys = append(fieldOrderKeys, i)
			}
		}

		// Get display name
		if val, ok := tags[displayField]; ok {
			reportStruct.ColDisplayNames[field.Name] = val
		} else {
			reportStruct.ColDisplayNames[field.Name] = field.Name
		}

		// Get the check
		if val, ok := tags[controlTestField]; ok {
			reportStruct.ColTags[val] = field.Name
		}
	}

	// order the display columns
	sort.Ints(fieldOrderKeys)
	for _, key := range fieldOrderKeys {
		reportStruct.ColDisplayOrder = append(reportStruct.ColDisplayOrder, fieldOrder[key])
	}
	return
}

// RowToStructRow Builds structured row of evidence
func RowToStructRow(row interface{}, entityIdFieldName string, rowFieldNames []string) (reportRow *receptor_v1.Row) {
	reportRow = &receptor_v1.Row{
		EntityInstanceId: getField(row, entityIdFieldName),
		Cols:             map[string]*receptor_v1.Value{},
	}

	rowValue := reflect.Indirect(reflect.ValueOf(row))
	for _, fieldName := range rowFieldNames {
		v := rowValue.FieldByName(fieldName)

		if dateTime, ok := v.Interface().(time.Time); ok {
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_TimestampValue{
					TimestampValue: timestamppb.New(dateTime),
				},
			}
			continue
		}

		dateTime, ok := v.Interface().(*time.Time)
		if ok {
			reportRow.Cols[fieldName] = &receptor_v1.Value{}
			if dateTime != nil {
				reportRow.Cols[fieldName].ValueType = &receptor_v1.Value_TimestampValue{
					TimestampValue: timestamppb.New(*dateTime),
				}
			}
			continue
		}

		switch v.Kind() {
		case reflect.Bool:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_BoolValue{
					BoolValue: v.Bool(),
				},
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Int32Value{
					Int32Value: int32(v.Int()),
				},
			}
			break
		case reflect.Int64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Int64Value{
					Int64Value: v.Int(),
				},
			}
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Uint32Value{
					Uint32Value: uint32(v.Uint()),
				},
			}
			break
		case reflect.Uint64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Uint64Value{
					Uint64Value: v.Uint(),
				},
			}
			break
		case reflect.Float32, reflect.Float64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_DoubleValue{
					DoubleValue: v.Float(),
				},
			}
			break
		case reflect.String:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_StringValue{
					StringValue: v.String(),
				},
			}
			break
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

func multipartEvidence(finding *receptor_v1.Finding) (contentType string, evidencePath string, err error) {
	if len(finding.Evidences) == 0 {
		err = errors.New("no evidence found")
		log.Error().Msg("no evidence found")
		return
	}

	evidence := finding.Evidences[0]
	if evidence.GetDoc() == nil {
		err = errors.New("evidence doc is nil")
		log.Error().Msg("evidence doc is nil")
	} else {
		// evidence should be protobuf of evidence + blob in a mulitpart/mixed
		// the mime of the part should be the mime from the evidence.doc.Mime
		dstFile, err := os.CreateTemp("", "multipart-evidence_*.tmp")
		if err != nil {
			log.Error().Msgf("failed to create multipart file: %v", err)
			return "", "", err
		}

		mime := evidence.GetDoc().GetMime()
		body := evidence.GetDoc().GetBody()
		streamFilePath := evidence.GetDoc().GetStreamFilePath()

		bufferSize := multipartkit.DefaultBufferSize

		// Initialize the multipart builder
		builder, err := multipartkit.NewMultipartBuilder(dstFile, bufferSize)
		defer func() {
			err = builder.Finalize()
			if err != nil {
				log.Error().Msgf("failed to finalize multipart builder: %v", err)
			}
		}()

		if err != nil {
			log.Error().Msgf("failed to create multipart builder: %v", err)
			return "", "", err
		}

		boundary := builder.GetBoundary()
		contentType = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

		// 1. Part1 : protobuf of Finding without evidence
		err = builder.AddProtobuf("receptor_v1.Finding", finding)
		if err != nil {
			log.Error().Msgf("failed to add protobuf message: %v", err)
		}

		//2. Part2 : evidence blob
		if streamFilePath != "" {
			err = builder.AddFile(evidence.Caption, streamFilePath, mime)
		} else {
			err = builder.AddBytes(evidence.Caption, evidence.Caption, mime, body)
		}
		if err != nil {
			log.Error().Msgf("failed to add blob part: %v", err)
		}

		//3. Part3 : Sources
		err = builder.AddProtobuf("receptor_v1.Sources", evidence.Sources)
		// evidence.Sources
		return contentType, dstFile.Name(), nil
	}
	return
}
