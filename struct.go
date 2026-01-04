package mapper

import (
	"reflect"
)

// assignStruct handles nested struct assignment by recursively mapping fields.
// It ensures that:
// - struct fields are mapped by name or tag
// - nested structs are recursively processed
// - a new struct is created (deep copy behavior)
func assignStruct(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string, depth int) error {
	if depth <= 0 {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "maximum nesting depth exceeded (possible circular reference)",
		}
	}

	srcType := src.Type()
	dstType := dst.Type()

	srcMeta, err := getStructMeta(srcType, tagName)
	if err != nil {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "failed to get source struct metadata: " + err.Error(),
		}
	}

	if srcType == dstType && !srcMeta.HasComposite {
		dst.Set(src)
		return nil
	}

	dstMeta, err := getStructMeta(dstType, tagName)
	if err != nil {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "failed to get destination struct metadata: " + err.Error(),
		}
	}

	for dstName, dstFieldMeta := range dstMeta.FieldsByName {
		srcFieldMeta, ok := srcMeta.FieldsByName[dstName]

		if !ok {
			srcFieldMeta, ok = srcMeta.FieldsByTag[dstName]
		}

		if !ok {
			continue
		}

		srcField := src.FieldByIndex(srcFieldMeta.Index)
		dstField := dst.FieldByIndex(dstFieldMeta.Index)

		nestedPath := fieldPath + "." + dstName

		if err := assignNestedValue(dstField, srcField, srcStructType, dstStructType, nestedPath, tagName, srcFieldMeta.ConvertTo, depth); err != nil {
			return err
		}
	}

	return nil
}

// assignNestedValue handles value assignment within nested contexts (structs, slices, maps).
// It supports nested structs, slices, maps, pointers, and type conversions.
func assignNestedValue(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName, convertTo string, depth int) error {
	if depth <= 0 {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "maximum nesting depth exceeded (possible circular reference)",
		}
	}

	if !dst.CanSet() {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "destination field cannot be set",
		}
	}

	sType := src.Type()
	dType := dst.Type()

	if convertTo != "" && sType.Kind() == reflect.String {
		converted, err := convertString(src.String(), convertTo, srcStructType, dstStructType, fieldPath)
		if err != nil {
			return err
		}
		dst.Set(converted.Convert(dType))
		return nil
	}

	srcKind := sType.Kind()
	dstKind := dType.Kind()

	if srcKind == reflect.Struct && dstKind == reflect.Struct {
		return assignStruct(dst, src, srcStructType, dstStructType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Slice && dstKind == reflect.Slice {
		return assignSlice(dst, src, srcStructType, dstStructType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Map && dstKind == reflect.Map {
		return assignMap(dst, src, srcStructType, dstStructType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Ptr && dstKind == reflect.Ptr {
		if src.IsNil() {
			dst.Set(reflect.Zero(dType))
			return nil
		}
		newPtr := reflect.New(dType.Elem())
		if err := assignNestedValue(newPtr.Elem(), src.Elem(), srcStructType, dstStructType, fieldPath, tagName, convertTo, depth-1); err != nil {
			return err
		}
		dst.Set(newPtr)
		return nil
	}

	if srcKind == reflect.Ptr && dstKind != reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		return assignNestedValue(dst, src.Elem(), srcStructType, dstStructType, fieldPath, tagName, convertTo, depth-1)
	}

	if srcKind != reflect.Ptr && dstKind == reflect.Ptr {
		newPtr := reflect.New(dType.Elem())
		if err := assignNestedValue(newPtr.Elem(), src, srcStructType, dstStructType, fieldPath, tagName, convertTo, depth-1); err != nil {
			return err
		}
		dst.Set(newPtr)
		return nil
	}

	if sType.AssignableTo(dType) {
		dst.Set(src)
		return nil
	}

	if sType.ConvertibleTo(dType) {
		dst.Set(src.Convert(dType))
		return nil
	}

	return &MappingError{
		SrcType:   srcStructType.String(),
		DstType:   dstStructType.String(),
		FieldPath: fieldPath,
		Reason:    "incompatible field types: " + sType.String() + " -> " + dType.String(),
	}
}
