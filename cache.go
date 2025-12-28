package mapper

import (
	"reflect"
	"sync"
)

type fieldMeta struct {
	Name      string
	Index     []int
	Type      reflect.Type
	Tag       string
	ConvertTo string // Target type for string conversion (e.g., "int", "float64", "bool")
}

type structMeta struct {
	Type         reflect.Type
	FieldsByName map[string]fieldMeta
	FieldsByTag  map[string]fieldMeta
}

type metaCacheKey struct {
	Type    reflect.Type
	TagName string
}

var metaCache sync.Map // map[metaCacheKey]*structMeta

func getStructMeta(t reflect.Type, tagName string) (*structMeta, error) {
	if t.Kind() != reflect.Struct {
		return nil, &MappingError{
			SrcType:   "",
			DstType:   "",
			FieldPath: "",
			Reason:    "type is not a struct",
		}
	}

	key := metaCacheKey{Type: t, TagName: tagName}
	if v, ok := metaCache.Load(key); ok {
		return v.(*structMeta), nil
	}

	m := &structMeta{
		Type:         t,
		FieldsByName: make(map[string]fieldMeta),
		FieldsByTag:  make(map[string]fieldMeta),
	}

	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		sf := t.Field(i)

		// Skip unexported fields.
		if !sf.IsExported() {
			continue
		}

		meta := fieldMeta{
			Name:  sf.Name,
			Index: sf.Index,
			Type:  sf.Type,
		}

		if convTag := sf.Tag.Get("mapconv"); convTag != "" {
			meta.ConvertTo = convTag
		}

		m.FieldsByName[sf.Name] = meta

		if tagName != "" {
			if tag := sf.Tag.Get(tagName); tag != "" {
				meta.Tag = tag
				m.FieldsByTag[tag] = meta
			}
		}
	}

	metaCache.Store(key, m)
	return m, nil
}
