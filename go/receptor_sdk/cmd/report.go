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

const mulitpartPrefix = "multipart/tr-mixed"

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

	go receptorImpl.ReportBatch(credentials, evidenceChannel)

	for evidences := range evidenceChannel {
		// Receive evidence and report them one batch at a time
		err := reportEvidence(rc, &finding, evidences)
		if err != nil {
			log.Err(err).Msg("failed to report evidence")
			// Continue on to next batch even after an error
			continue
		}

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
			Caption:            evidence.Caption,
			Description:        evidence.Description,
			ServiceName:        evidence.ServiceName,
			EntityType:         evidence.EntityType,
			Sources:            evidence.Sources,
			ServiceAccountId:   evidence.ServiceAccountId,
			Controls:           evidence.Controls,
			IsManual:           evidence.IsManual,
			RelevantDate:       &evidence.RelevantDate,
			EvidenceObjectType: evidence.EvidenceObjectType,
			Policies:           evidence.Policies,
			RecordIds:          evidence.RecordIds,
			EvidenceKey:        evidence.EvidenceKey,
			Exceptions:         evidence.Exceptions,
			EvidenceLink:       evidence.EvidenceLink,
		}

		if evidence.Document != nil && len(*evidence.Document) > 0 {
			paths := []FilePathsInfo{}
			evidenceDocuments := receptor_v1.Documents{}
			evidenceDocuments.Docs = []*receptor_v1.Document{}
			reportFinding := receptor_v1.Finding{
				ReceptorType:           finding.ReceptorType,
				ServiceProviderAccount: finding.ServiceProviderAccount,
			}

			for _, doc := range *evidence.Document {
				newDoc := receptor_v1.Document{
					Body:     doc.Body,
					Mime:     doc.Mime,
					FileName: doc.FileName,
					Metadata: doc.Metadata,
				}
				if doc.LastModified != nil {
					newDoc.LastModified = doc.LastModified
				}
				evidenceDocuments.Docs = append(evidenceDocuments.Docs, &newDoc)
				// add to paths only if it is NOT bytes
				if len(doc.Body) == 0 {
					paths = append(paths, FilePathsInfo{
						Path:     doc.StreamFilePath,
						Metadata: doc.Metadata,
						FileName:  doc.FileName,
						PartName: doc.FileName,
					})
				}
			}
			if len(evidenceDocuments.Docs) == 1 { // single document
				reportEvidence.EvidenceType = &receptor_v1.Evidence_Doc{
					Doc: evidenceDocuments.Docs[0],
				}
			} else if len(evidenceDocuments.Docs) > 1 {
				reportEvidence.EvidenceType = &receptor_v1.Evidence_Docs{
					Docs: &evidenceDocuments,
				}
			}
			//extract sources and add to multipart and remove from finding
			sources := []*receptor_v1.Source{}
			for _, source := range evidence.Sources {
				sources = append(sources, &receptor_v1.Source{
					RawApiRequest:  source.RawApiRequest,
					RawApiResponse: source.RawApiResponse,
				})
			}
			reportEvidence.Sources = []*receptor_v1.Source{}

			reportFinding.Evidences = append(reportFinding.Evidences, &reportEvidence)

			contentType, streamFile, err := multipartEvidence(&reportFinding, paths, sources)

			// have the streamFile from receptor - remove the temp evidence files
			for _, doc := range *evidence.Document {
				os.Remove(doc.StreamFilePath)
			}

			if err != nil {
				log.Err(err).Msg("failed to create multipart evidence")
				err = nil
				continue
			}

			// make a multipart file and then stream it
			stream, err := rc.StreamReport(context.Background())
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
			_, err = stream.CloseAndRecv()
			if err != nil {
				log.Err(err).Msg("failed to close and receive stream")
				continue
			}
		} else { // evidence is structured
			reportEvidence.EvidenceType = &receptor_v1.Evidence_Struct{Struct: &reportStruct}

			// Convert rows
			var entityIdFieldName string
			var rowFieldNames []string
			for idx, row := range evidence.Rows {
				if idx == 0 {
					if entityIdFieldName, rowFieldNames, err = ExtractMetaData(row, &reportStruct); err != nil {
						return // failed to extract metadata, likely an invalid row type
					}
				}
				reportStruct.Rows = append(reportStruct.Rows, RowToStructRow(row, entityIdFieldName, rowFieldNames))
			}

			// Append to Finding
			finding.Evidences = append(finding.Evidences, &reportEvidence)
		}

	}
	// report all structured evidence at once
	_, err = rc.Report(context.Background(), finding)
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
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Int32Value{
					Int32Value: int32(v.Int()),
				},
			}
		case reflect.Int64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Int64Value{
					Int64Value: v.Int(),
				},
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Uint32Value{
					Uint32Value: uint32(v.Uint()),
				},
			}
		case reflect.Uint64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_Uint64Value{
					Uint64Value: v.Uint(),
				},
			}
		case reflect.Float32, reflect.Float64:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_DoubleValue{
					DoubleValue: v.Float(),
				},
			}
		case reflect.String:
			reportRow.Cols[fieldName] = &receptor_v1.Value{
				ValueType: &receptor_v1.Value_StringValue{
					StringValue: v.String(),
				},
			}
		case reflect.Slice:
			elemType := v.Type().Elem()
			if elemType.Kind() == reflect.String {
				var stringArray receptor_v1.StringList
				for i := 0; i < v.Len(); i++ {
					stringArray.Values = append(stringArray.Values, v.Index(i).String())
				}
				reportRow.Cols[fieldName] = &receptor_v1.Value{
					ValueType: &receptor_v1.Value_StringListValue{
						StringListValue: &stringArray,
					},
				}
			} else if elemType.Kind() == reflect.Struct {
				var structList receptor_v1.StructList
				for i := 0; i < v.Len(); i++ {
					structElem := v.Index(i)
					structMap := make(map[string]*receptor_v1.Value)
					for j := 0; j < structElem.NumField(); j++ {
						structField := structElem.Type().Field(j)
						fieldVal := structElem.Field(j).Interface()
						structMap[structField.Name] = mapToValue(fieldVal)
					}
					structList.Values = append(structList.Values, &receptor_v1.StructStruct{Fields: structMap})
				}
				reportRow.Cols[fieldName] = &receptor_v1.Value{
					ValueType: &receptor_v1.Value_StructListValue{
						StructListValue: &structList,
					},
				}
			} else {
				log.Warn().Msg("unsupported evidence row field (" + fieldName + ") type")
			}
		default:
			log.Warn().Msg("unsupported evidence row field (" + fieldName + ") type")
		}
	}

	return
}

func mapToValue(val interface{}) *receptor_v1.Value {
	switch v := val.(type) {
	case bool:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_BoolValue{BoolValue: v},
		}
	case int, int8, int16, int32:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_Int32Value{Int32Value: int32(reflect.ValueOf(v).Int())},
		}
	case int64:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_Int64Value{Int64Value: v},
		}
	case uint, uint8, uint16, uint32:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_Uint32Value{Uint32Value: uint32(reflect.ValueOf(v).Uint())},
		}
	case uint64:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_Uint64Value{Uint64Value: v},
		}
	case float32, float64:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_DoubleValue{DoubleValue: reflect.ValueOf(v).Float()},
		}
	case string:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_StringValue{StringValue: v},
		}
	case time.Time:
		return &receptor_v1.Value{
			ValueType: &receptor_v1.Value_TimestampValue{TimestampValue: timestamppb.New(v)},
		}
	case *time.Time:
		if v != nil {
			return &receptor_v1.Value{
				ValueType: &receptor_v1.Value_TimestampValue{TimestampValue: timestamppb.New(*v)},
			}
		}
		return &receptor_v1.Value{} // Return empty value for nil *time.Time
	default:
		// Handle unsupported types as needed
		return &receptor_v1.Value{} // Return a default value if needed
	}
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

type FilePathsInfo struct {
<<<<<<< HEAD
	FileName string
	PartName string
=======
>>>>>>> main
	Path     string
	Metadata map[string]string
}

func multipartEvidence(finding *receptor_v1.Finding, streamFilePathsInfo []FilePathsInfo, sources []*receptor_v1.Source) (contentType string, evidencePath string, err error) {
	if len(finding.Evidences) == 0 {
		err = errors.New("no evidence found")
		log.Error().Msg("no evidence found")
		return
	}

	evidence := finding.Evidences[0]

	if evidence.EvidenceType == nil {
		err = errors.New("evidence doc(s) is nil")
		log.Error().Msg("evidence doc(s) is nil")
	} else {

		// evidence should be protobuf of evidence + blob in a multipart/mixed
		// the mime of the part should be the mime from the evidence.doc.Mime
		dstFile, err := os.CreateTemp("", "multipart-evidence_*.tmp")
		if err != nil {
			log.Err(err).Msg("failed to create multipart file")
			return "", "", err
		}
<<<<<<< HEAD
		// TOD: Adjust mime for multiple document evidences
		mime := evidence.GetDoc().GetMime() 
=======
		mime := evidence.GetDoc().GetMime()
>>>>>>> main
		bufferSize := multipartkit.DefaultBufferSize

		// Initialize the multipart builder
		builder, err := multipartkit.NewMultipartBuilder(dstFile, bufferSize)
		defer func() {
			err = builder.Finalize()
			if err != nil {
				log.Err(err).Msg("failed to finalize multipart builder")
			}
		}()

		if err != nil {
			log.Error().Msgf("failed to create multipart builder: %v", err)
			return "", "", err
		}

		boundary := builder.GetBoundary()
		contentType = fmt.Sprintf("%s; %s; boundary=%s", mulitpartPrefix, mime, boundary)

		// 1. Part1 : protobuf of Finding without evidence

		err = builder.AddProtobuf("receptor_v1.Finding", finding) // need to remove evidences from this finding ...
		if err != nil {
			log.Error().Msgf("failed to add protobuf message: %v", err)
		}

		// 2. Part2 : evidence blob
		docs := []*receptor_v1.Document{}

		switch evidenceDocType := evidence.EvidenceType.(type) {
		case *receptor_v1.Evidence_Doc:
			docs = append(docs, evidenceDocType.Doc)
			mime = evidence.GetDoc().GetMime()
		case *receptor_v1.Evidence_Docs:
			docs = evidence.GetDocs().Docs
			mime = "application/tr-archive"
		}
		contentType = fmt.Sprintf("%s; %s; boundary=%s", mulitpartPrefix, mime, boundary)
		for _, doc := range docs {
			if len(doc.Body) > 0 {
				name := evidence.Caption
				if doc.FileName != "" {
					name = doc.FileName
				}
				err = builder.AddBytes(name, name, mime, doc.GetBody(), doc.GetMetadata())
				if err != nil {
					log.Err(err).Msgf("failed to add blob part: %s", evidence.Caption)
				}
			}

		}
		// 3. Part3 : evidence paths
		for _, streamFilePathInfo := range streamFilePathsInfo {
			if streamFilePathInfo.Path != "" {
				err = builder.AddFile(streamFilePathInfo.PartName, streamFilePathInfo.FileName, streamFilePathInfo.Path, mime, streamFilePathInfo.Metadata)
				if err != nil {
					log.Err(err).Msgf("failed to add stream file: %s", streamFilePathInfo.Path)
				}
			}
		}

		// 4. Part4 : Sources
		err = builder.AddProtobuf("receptor_v1.Sources", &receptor_v1.Sources{
			Sources: sources,
		})
		if err != nil {
			log.Error().Msgf("failed to add sources part: %v", err)
		}
		// clear evidences from the temp finding
		finding.Evidences = []*receptor_v1.Evidence{}

		return contentType, dstFile.Name(), nil
	}
	return
}
