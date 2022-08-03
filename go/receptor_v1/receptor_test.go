package receptor_v1_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustero/api/go/pkg/printer"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type Foo struct {
	Id   string    `tr:"primary_key"`
	Name string    `tr:"DisplayName:Name Of Foo;DisplayOrder:2"`
	Date time.Time `tr:"DisplayName:Created on;DisplayOrder:1"`
}

func TestNewStruct(t *testing.T) {
	timestamp := int64(1658575538938)
	data := []interface{}{
		&Foo{
			Id:   "123",
			Name: "Foo",
			Date: time.UnixMilli(timestamp),
		},
	}

	newStruct, err := receptor_v1.NewStruct(data)
	assert.Nil(t, err)
	assert.NotNil(t, newStruct)
	assert.NotNil(t, newStruct.GetRows())
	assert.Equal(t, newStruct.GetRows()[0].ServiceId, "123")
	assert.NotNil(t, newStruct.GetColDisplayOrder())
	assert.NotNil(t, newStruct.GetColDisplayNames())
}

func TestNewEvidence(t *testing.T) {
	timestamp := int64(1658575538938)
	data := []interface{}{
		&Foo{
			Id:   "123",
			Name: "Foo",
			Date: time.UnixMilli(timestamp),
		},
	}

	newStruct, err := receptor_v1.NewStruct(data)
	assert.Nil(t, err)
	assert.NotNil(t, newStruct)
	assert.NotNil(t, newStruct.GetRows())
	assert.NotNil(t, newStruct.GetColDisplayOrder())
	assert.NotNil(t, newStruct.GetColDisplayNames())

	assert.Len(t, newStruct.GetColDisplayNames(), 3)
	assert.Equal(t, newStruct.GetColDisplayNames()["Date"], "Created on")
	assert.Equal(t, newStruct.GetColDisplayNames()["Name"], "Name Of Foo")
}

func TestTableFromMapFinding(t *testing.T) {
	data := make(map[string]*receptor_v1.Struct_Row_Value)
	data["hello"] = &receptor_v1.Struct_Row_Value{ValueType: &receptor_v1.Struct_Row_Value_StringValue{
		StringValue: "world",
	}}
	timestamp := int64(1658575538938)
	data["foo"] = &receptor_v1.Struct_Row_Value{
		ValueType: &receptor_v1.Struct_Row_Value_TimestampValue{
			TimestampValue: timestamppb.New(time.UnixMilli(timestamp)),
		},
	}
	mapFindings := []*receptor_v1.Struct_Row{
		{
			ServiceId: "123",
			Cols:      data,
		},
	}
	testStruct := &receptor_v1.Struct{
		Rows: mapFindings,
		ColDisplayNames: map[string]string{
			"hello": "Name",
			"foo":   "Date",
		},
		ColDisplayOrder: []string{"foo", "hello"},
	}

	table, err := testStruct.Tabulate(printer.SimpleStrconv)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(table.Body))
	assert.Equal(t, "2022-Jul-23", table.Body[0][0])
	assert.Equal(t, "world", table.Body[0][1])
}

type TestData struct {
	id    string    `tr:"primary_key"`
	hello string    `tr:"DisplayName:Name;DisplayOrder:2"`
	Foo   time.Time `tr:"DisplayName:Date;DisplayOrder:1"`
}

func TestTableFromAny(t *testing.T) {
	testDate := &TestData{
		id:    "123",
		hello: "world",
		Foo:   time.UnixMilli(int64(1658575538938)),
	}
	testStruct, _ := receptor_v1.NewStruct([]interface{}{testDate})

	table, err := testStruct.Tabulate(printer.SimpleStrconv)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(table.Body))
	assert.Equal(t, "2022-Jul-23", table.Body[0][0])
	assert.Equal(t, "world", table.Body[0][1])
}
