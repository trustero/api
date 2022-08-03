package proto_utils

import (
	"errors"
	"fmt"
	"github.com/trustero/api/go/pkg/printer"
	"github.com/trustero/api/go/pkg/tabular"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"time"
)

func Tabulate(message proto.Message, strconv printer.Strconv) (attributes []string, err error) {
	return tabulate(message.ProtoReflect(), strconv)
}

func tabulate(m protoreflect.Message, strconv printer.Strconv) (attributes []string, err error) {
	if strconv == nil {
		err = errors.New("string converter not set")
		return
	}
	for i := 0; i < m.Descriptor().Fields().Len(); i++ {
		fDesc, fVal := getFieldReflectionByIndex(m, i)
		var strVal string
		if strVal, err = printField(fDesc, fVal, strconv); err != nil {
			return
		}
		attributes = append(attributes, strVal)
	}

	return
}

func isFieldListType(message protoreflect.Message, idx int) bool {
	fd, _ := getFieldReflectionByIndex(message, idx)
	return fd.Kind() == protoreflect.StringKind && fd.Cardinality() == protoreflect.Repeated
}

func TabulateOneOf(message proto.Message, oneOfFieldName string, strconv printer.Strconv) (result *tabular.Table, err error) {
	oneOfValue := GetOneOfValue(message, oneOfFieldName)
	// Assuming the first field of the wrapper class contains is a list of finding elements
	findingElements := GetFieldAsList(oneOfValue, 0)

	result = &tabular.Table{}
	for i := 0; i < findingElements.Len(); i++ {
		row := findingElements.Get(i).Message()
		// Assume that, as specified by the evidence.Observations contract, the first field
		// of each element contains the id or list of ids for instances in that can filter the element out
		var attributes []string
		if attributes, err = tabulateRow(row, strconv); err != nil {
			return
		}
		result.Body = append(result.Body, attributes)
	}
	return
}

func tabulateRow(m protoreflect.Message, strconv printer.Strconv) (attributes []string, err error) {
	if isFieldListType(m, 1) {
		attributes, err = tabulateField(m, 1, strconv)
		return
	}

	if attributes, err = tabulate(m, strconv); err != nil {
		return
	}
	// skip the first attribute, which is the id
	attributes = attributes[1:]
	return
}
func TabulateField(message proto.Message, idx int, strconv printer.Strconv) (attributes []string, err error) {
	return tabulateField(message.ProtoReflect(), idx, strconv)
}

func tabulateField(m protoreflect.Message, idx int, strconv printer.Strconv) (attributes []string, err error) {
	if m.Descriptor().Oneofs().Len() > 0 {
		onefOfDescriptor := m.Descriptor().Oneofs().Get(idx)
		fv := m.WhichOneof(onefOfDescriptor)
		if onefOfDescriptor != nil {
			vd := m.Get(fv)
			var strVal string
			if strVal, err = printField(fv, vd, strconv); err == nil {
				attributes = append(attributes, strVal)
			}
			return
		}

	}
	fd, fv := getFieldReflectionByIndex(m, idx)

	if fd.Kind() == protoreflect.StringKind && fd.Cardinality() == protoreflect.Repeated {
		for i := 0; i < fv.List().Len(); i++ {
			attributes = append(attributes, fv.List().Get(i).String())
		}
		return
	}
	var strVal string
	if strVal, err = printField(fd, fv, strconv); err != nil {
		return
	}
	attributes = append(attributes, strVal)
	return
}
func getFieldReflectionByIndex(m protoreflect.Message, idx int) (protoreflect.FieldDescriptor, protoreflect.Value) {
	fieldDescriptor := m.Descriptor().Fields().Get(idx)
	fieldValue := m.Get(fieldDescriptor)
	return fieldDescriptor, fieldValue
}

// getFieldValueAsString
// Supported types
//const (
//	BoolKind     Kind = 8
//	Int32Kind    Kind = 5
//	Sint32Kind   Kind = 17
//	Uint32Kind   Kind = 13
//	Int64Kind    Kind = 3
//	Sint64Kind   Kind = 18
//	Uint64Kind   Kind = 4
//	Sfixed32Kind Kind = 15
//	Fixed32Kind  Kind = 7
//	FloatKind    Kind = 2
//	Sfixed64Kind Kind = 16
//	Fixed64Kind  Kind = 6
//	DoubleKind   Kind = 1
//	StringKind   Kind = 9
//)
func printField(fd protoreflect.FieldDescriptor, fv protoreflect.Value, customStrconv printer.Strconv) (value string, err error) {
	if customStrconv == nil {
		err = errors.New("string converter not set")
		return
	}
	switch fd.Kind() {
	case protoreflect.BoolKind:
		value = customStrconv.FormatBool(fv.Bool())
		return
	case protoreflect.EnumKind,
		protoreflect.BytesKind,
		protoreflect.MessageKind,
		protoreflect.GroupKind:
		if isTimestamp(fd) {
			value = customStrconv.FormatTimestamp(parseTimestamp(fv))
			return
		}
		if isDuration(fd) {
			value = customStrconv.FormatDuration(parseDuration(fv))
			return
		}

		err = errors.New("error parsing unknown value type")
		return
	case protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind:
		value = strconv.FormatInt(fv.Int(), 10)
		return
	case protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind:
		value = strconv.FormatUint(fv.Uint(), 10)
		return
	case protoreflect.FloatKind,
		protoreflect.DoubleKind:
		value = fmt.Sprintf("%.3f", fv.Float())
		return
	case protoreflect.StringKind:
		if fd.Cardinality() == protoreflect.Repeated {
			var items []string
			for i := 0; i < fv.List().Len(); i++ {
				items = append(items, fv.List().Get(i).String())
			}
			value = strings.Join(items, ", ")
			return
		}
		value = fv.String()
		return
	default:
		err = errors.New("error parsing unknown value type")
	}
	return
}

func parseTimestamp(fv protoreflect.Value) *timestamppb.Timestamp {
	message := fv.Message()
	descriptor := message.Descriptor()
	seconds := message.Get(descriptor.Fields().Get(0)).Int()
	nanos := int32(message.Get(descriptor.Fields().Get(1)).Int())
	return &timestamppb.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
}

func parseDuration(fv protoreflect.Value) *durationpb.Duration {
	message := fv.Message()
	descriptor := message.Descriptor()
	seconds := message.Get(descriptor.Fields().Get(0)).Int()
	nanos := int32(message.Get(descriptor.Fields().Get(1)).Int())
	return &durationpb.Duration{
		Seconds: seconds,
		Nanos:   nanos,
	}
}

func isTimestamp(fd protoreflect.FieldDescriptor) bool {
	return timestamppb.Now().ProtoReflect().Descriptor().FullName() == fd.Message().FullName()
}
func isDuration(fd protoreflect.FieldDescriptor) bool {
	return durationpb.New(time.Nanosecond).ProtoReflect().Descriptor().FullName() == fd.Message().FullName()
}

func GetFieldAsList(value protoreflect.Value, fieldIndex int) protoreflect.List {
	oneOfFieldMessageDescriptor := value.Message().Descriptor()
	oneOfFirstFieldDescriptor := oneOfFieldMessageDescriptor.Fields().Get(fieldIndex)
	findingElements := value.Message().Get(oneOfFirstFieldDescriptor).List()
	return findingElements
}

// GetOneOfValue
// Digs into a proto message where a one-of type is known to exist, and returs the value of the one-of instance in the message.
func GetOneOfValue(protoMessage protoreflect.ProtoMessage, oneofFieldName string) protoreflect.Value {
	findingReflectDescriptor := protoMessage.ProtoReflect().Descriptor()
	findTypeOneOfFieldDescriptor := findingReflectDescriptor.Oneofs().ByName(protoreflect.Name(oneofFieldName)).Fields()
	var rowsValue protoreflect.Value
	// Loop over the OneOf of evidence.Observations and find the *one* field that is set
	for i := 0; i < findingReflectDescriptor.Fields().Len(); i++ {
		if protoMessage.ProtoReflect().Has(findTypeOneOfFieldDescriptor.Get(i)) {
			rowsValue = protoMessage.ProtoReflect().Get(findTypeOneOfFieldDescriptor.Get(i))
			break
		}
	}
	return rowsValue
}

func GetId(message protoreflect.Message) (instanceIds []string, err error) {
	// Assume that, as specified by the evidence.Observations contract, the first field
	// of each element contains the id or list of ids for instances in that can filter the element out
	instanceIdFieldDescriptor := message.Descriptor().Fields().Get(0)
	instanceIdFieldValue := message.Get(instanceIdFieldDescriptor)

	if instanceIdFieldDescriptor.Kind() != protoreflect.StringKind {
		err = errors.New("expected string or repeated string in field 1 of message")
		return
	}
	if instanceIdFieldDescriptor.Cardinality() == protoreflect.Repeated {
		for k := 0; k < instanceIdFieldValue.List().Len(); k++ {
			instanceIds = append(instanceIds, instanceIdFieldValue.List().Get(k).String())
		}
		return
	}
	instanceIds = append(instanceIds, instanceIdFieldValue.String())
	return
}
