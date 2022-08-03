package proto_utils

import (
	"fmt"
	reflect_utils "github.com/trustero/api/go/pkg/reflect-utils"
	"reflect"
	"strconv"
	"strings"
)

type fieldTag struct {
	IsId         bool
	DisplayName  string
	DisplayOrder int
}

type TrTag struct {
	Id           string
	DisplayName  map[string]string
	DisplayOrder []string
}

const (
	tagName        = "tr"
	tfDisplayName  = "DisplayName"
	tfDisplayOrder = "DisplayOrder"
	tfId           = "primary_key"
)

func NewTrTag(anyVal interface{}) (tag *TrTag, err error) {
	reflectedValue, reflectedField := reflect_utils.GetValueAndType(anyVal)
	idIdx := -1
	tag = &TrTag{
		DisplayName:  make(map[string]string),
		DisplayOrder: make([]string, reflectedField.NumField()),
	}
	for i := 0; i < reflectedValue.NumField(); i++ {
		field := reflectedField.Field(i)
		var ft *fieldTag
		if ft, err = getFieldTag(field); err != nil || ft == nil {
			continue
		}
		tag.DisplayName[field.Name] = ft.DisplayName
		tag.DisplayOrder[ft.DisplayOrder] = field.Name
		if ft.IsId {
			if idIdx >= 0 {
				err = fmt.Errorf("duplicate Id field")
				return
			}
			idIdx = i
		}
	}

	return
}

func getFieldTag(field reflect.StructField) (ft *fieldTag, err error) {
	tagValue := field.Tag.Get(tagName)
	if tagValue == "" {
		return
	}
	ft = &fieldTag{}
	for _, tagSection := range strings.Split(tagValue, ";") {
		tagSection = strings.TrimSpace(tagSection)
		if tagSection == "" {
			continue
		}
		tagKeyValue := strings.Split(tagSection, ":")
		if len(tagKeyValue) != 2 {
			if tagSection != tfId {
				err = fmt.Errorf("invalid tag: %s", tagSection)
				return
			}
			if field.Type.Kind() != reflect.String {
				err = fmt.Errorf("Id field must be string")
				return
			}
			ft.IsId = true
			continue
		}
		sectionKey := strings.TrimSpace(tagKeyValue[0])
		sectionValue := strings.TrimSpace(tagKeyValue[1])
		switch sectionKey {
		case tfDisplayName:
			ft.DisplayName = sectionValue
		case tfDisplayOrder:
			if ft.DisplayOrder, err = strconv.Atoi(sectionValue); err != nil || ft.DisplayOrder < 0 {
				err = fmt.Errorf("invalid order for field %s", field.Name)
				return
			}
			ft.DisplayOrder-- // cardinality starts at 1, but arrays start at 0
		case tfId:
			if field.Type.Kind() != reflect.String {
				err = fmt.Errorf("Id field must be string")
				return
			}
			ft.IsId = true
		default:
			err = fmt.Errorf("invalid tag: %s", tagSection)
			return
		}
	}
	return
}
