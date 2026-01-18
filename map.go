package mapper

import (
	"reflect"
	"strconv"
)

// buildMapPath constructs a field path with map key notation.
// Only called when an error occurs to avoid allocation in the hot path.
func buildMapPath(basePath string, key reflect.Value) string {
	keyStr := formatMapKey(key)
	// Pre-allocate buffer: basePath + "[" + keyStr + "]"
	buf := make([]byte, 0, len(basePath)+len(keyStr)+2)
	buf = append(buf, basePath...)
	buf = append(buf, '[')
	buf = append(buf, keyStr...)
	buf = append(buf, ']')
	return string(buf)
}

// formatMapKey converts a map key to string representation efficiently.
func formatMapKey(key reflect.Value) string {
	switch key.Kind() {
	case reflect.String:
		return key.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(key.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(key.Uint(), 10)
	default:
		// Fallback for other types - this is rare
		return "<key>"
	}
}

// prependMapKeyPath prepends the map key path to a MappingError's FieldPath.
// This is called only when an error occurs, making path building lazy.
func prependMapKeyPath(err error, basePath string, key reflect.Value) error {
	if me, ok := err.(*MappingError); ok {
		keyPath := buildMapPath(basePath, key)
		if me.FieldPath != "" {
			me.FieldPath = keyPath + "." + me.FieldPath
		} else {
			me.FieldPath = keyPath
		}
	}
	return err
}

// assignMap handles map assignment with proper deep copying and type conversion.
// It ensures that:
// - nil maps remain nil
// - empty maps remain empty (not nil)
// - a new underlying map is created (modifications to source don't affect destination)
// - key and value types are converted if compatible
// - nested structs within maps are properly mapped using the provided tagName
func assignMap(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string, depth int) error {
	if depth <= 0 {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "maximum nesting depth exceeded (possible circular reference)",
		}
	}

	if src.IsNil() {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	sType := src.Type()
	dType := dst.Type()
	srcKeyType := sType.Key()
	dstKeyType := dType.Key()
	srcValType := sType.Elem()
	dstValType := dType.Elem()

	keysAssignable := srcKeyType.AssignableTo(dstKeyType)
	keysConvertible := srcKeyType.ConvertibleTo(dstKeyType)

	if !keysAssignable && !keysConvertible {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "map key types are incompatible: " + srcKeyType.String() + " -> " + dstKeyType.String(),
		}
	}

	srcValKind := srcValType.Kind()
	dstValKind := dstValType.Kind()

	valuesAreStructs := srcValKind == reflect.Struct && dstValKind == reflect.Struct
	valuesAreNestedMaps := srcValKind == reflect.Map && dstValKind == reflect.Map
	valuesAreNestedSlices := srcValKind == reflect.Slice && dstValKind == reflect.Slice
	valuesArePtrs := srcValKind == reflect.Ptr && dstValKind == reflect.Ptr
	valuesAssignable := srcValType.AssignableTo(dstValType)
	valuesConvertible := srcValType.ConvertibleTo(dstValType)

	if !valuesAssignable && !valuesConvertible && !valuesAreStructs && !valuesAreNestedMaps && !valuesAreNestedSlices && !valuesArePtrs {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "map value types are incompatible: " + srcValType.String() + " -> " + dstValType.String(),
		}
	}

	newMap := reflect.MakeMapWithSize(dType, src.Len())

	needsProcessing := valuesAreStructs || valuesAreNestedMaps || valuesAreNestedSlices || valuesArePtrs || (!valuesAssignable && valuesConvertible)

	iter := src.MapRange()
	for iter.Next() {
		srcKey := iter.Key()
		srcVal := iter.Value()

		var dstKey reflect.Value
		if keysAssignable {
			dstKey = srcKey
		} else {
			dstKey = srcKey.Convert(dstKeyType)
		}

		var dstVal reflect.Value
		var err error

		if !needsProcessing && valuesAssignable {
			dstVal = srcVal
		} else if valuesAreStructs {
			dstVal = reflect.New(dstValType).Elem()
			// Pass empty path; path is built only on error (lazy)
			err = assignStruct(dstVal, srcVal, srcStructType, dstStructType, "", tagName, depth-1)
			if err != nil {
				return prependMapKeyPath(err, fieldPath, srcKey)
			}
		} else if valuesAreNestedMaps {
			dstVal = reflect.New(dstValType).Elem()
			// Pass empty path; path is built only on error (lazy)
			err = assignMap(dstVal, srcVal, srcStructType, dstStructType, "", tagName, depth-1)
			if err != nil {
				return prependMapKeyPath(err, fieldPath, srcKey)
			}
		} else if valuesAreNestedSlices {
			dstVal = reflect.New(dstValType).Elem()
			// Pass empty path; path is built only on error (lazy)
			err = assignSlice(dstVal, srcVal, srcStructType, dstStructType, "", tagName, depth-1)
			if err != nil {
				return prependMapKeyPath(err, fieldPath, srcKey)
			}
		} else if valuesArePtrs {
			dstVal = reflect.New(dstValType).Elem()
			// Pass empty path; path is built only on error (lazy)
			err = assignPointerElement(dstVal, srcVal, srcStructType, dstStructType, "", tagName, depth-1)
			if err != nil {
				return prependMapKeyPath(err, fieldPath, srcKey)
			}
		} else if valuesConvertible {
			dstVal = srcVal.Convert(dstValType)
		}

		newMap.SetMapIndex(dstKey, dstVal)
	}

	dst.Set(newMap)
	return nil
}
