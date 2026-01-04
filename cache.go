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

// Two-level cache: first by reflect.Type (most common), then by tagName
// This optimizes the common case where tagName is constant across an application
var (
	typeCache     = make(map[reflect.Type]*tagCache)
	typeCacheLock sync.RWMutex
)

type tagCache struct {
	sync.RWMutex
	entries map[string]*structMeta
}

func getStructMeta(t reflect.Type, tagName string) (*structMeta, error) {
	if t.Kind() != reflect.Struct {
		return nil, &MappingError{
			SrcType:   "",
			DstType:   "",
			FieldPath: "",
			Reason:    "type is not a struct",
		}
	}

	// Fast path: read lock to check cache
	typeCacheLock.RLock()
	tc := typeCache[t]
	typeCacheLock.RUnlock()

	if tc != nil {
		tc.RLock()
		if m := tc.entries[tagName]; m != nil {
			tc.RUnlock()
			return m, nil
		}
		tc.RUnlock()
	}

	// Slow path: compute metadata
	m := buildStructMeta(t, tagName)

	// Store in cache
	typeCacheLock.Lock()
	tc = typeCache[t]
	if tc == nil {
		tc = &tagCache{entries: make(map[string]*structMeta, 2)}
		typeCache[t] = tc
	}
	typeCacheLock.Unlock()

	tc.Lock()
	// Double-check in case another goroutine computed it
	if existing := tc.entries[tagName]; existing != nil {
		tc.Unlock()
		return existing, nil
	}
	tc.entries[tagName] = m
	tc.Unlock()

	return m, nil
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

		// Check for composite types that need deep copy
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
