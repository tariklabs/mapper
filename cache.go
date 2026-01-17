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
	ConvertTo string
}

type structMeta struct {
	Type         reflect.Type
	FieldsByName map[string]*fieldMeta
	FieldsByTag  map[string]*fieldMeta
	HasComposite bool
}

type cacheKey struct {
	typ     reflect.Type
	tagName string
}

var structMetaCache sync.Map // map[cacheKey]*structMeta

func getStructMeta(t reflect.Type, tagName string) (*structMeta, error) {
	if t.Kind() != reflect.Struct {
		return nil, &MappingError{
			SrcType:   "",
			DstType:   "",
			FieldPath: "",
			Reason:    "type is not a struct",
		}
	}

	key := cacheKey{typ: t, tagName: tagName}

	if cached, ok := structMetaCache.Load(key); ok {
		return cached.(*structMeta), nil
	}

	m := buildStructMeta(t, tagName)

	actual, _ := structMetaCache.LoadOrStore(key, m)
	return actual.(*structMeta), nil
}

func buildStructMeta(t reflect.Type, tagName string) *structMeta {
	numFields := t.NumField()

	m := &structMeta{
		Type:         t,
		FieldsByName: make(map[string]*fieldMeta, numFields),
		FieldsByTag:  make(map[string]*fieldMeta, numFields),
		HasComposite: false,
	}

	for i := 0; i < numFields; i++ {
		sf := t.Field(i)

		if !sf.IsExported() {
			continue
		}

		switch sf.Type.Kind() {
		case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Struct:
			m.HasComposite = true
		}

		meta := &fieldMeta{
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

	return m
}
