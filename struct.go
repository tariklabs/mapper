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

		// Pass base path and field name separately; path is only built on error
		if err := assignNestedValue(dstField, srcField, srcStructType, dstStructType, fieldPath, dstName, tagName, srcFieldMeta.ConvertTo, depth); err != nil {
			return err
		}
	}

	return nil
}

// buildPath constructs the full field path from base path and field name.
// Only called when an error occurs to avoid allocation in the hot path.
func buildPath(basePath, fieldName string) string {
	if basePath == "" {
		return fieldName
	}
	return basePath + "." + fieldName
}

// assignNestedValue handles value assignment within nested contexts (structs, slices, maps).
// It supports nested structs, slices, maps, pointers, and type conversions.
// basePath and fieldName are kept separate to avoid string concatenation in the hot path;
// the full path is only built when an error occurs.
func assignNestedValue(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, basePath, fieldName, tagName, convertTo string, depth int) error {
	if depth <= 0 {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: buildPath(basePath, fieldName),
			Reason:    "maximum nesting depth exceeded (possible circular reference)",
		}
	}

	if !dst.CanSet() {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: buildPath(basePath, fieldName),
			Reason:    "destination field cannot be set",
		}
	}

	sType := src.Type()
	dType := dst.Type()

	if convertTo != "" && sType.Kind() == reflect.String {
		fullPath := buildPath(basePath, fieldName)
		converted, err := convertString(src.String(), convertTo, srcStructType, dstStructType, fullPath)
		if err != nil {
			return err
		}
		dst.Set(converted.Convert(dType))
		return nil
	}

	srcKind := sType.Kind()
	dstKind := dType.Kind()

	// Fast path: directly assignable or convertible types (most common for primitive fields)
	// Check these first to avoid path building for the majority of field assignments
	if srcKind != reflect.Struct && srcKind != reflect.Slice && srcKind != reflect.Map && srcKind != reflect.Ptr &&
		dstKind != reflect.Struct && dstKind != reflect.Slice && dstKind != reflect.Map && dstKind != reflect.Ptr {
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
			FieldPath: buildPath(basePath, fieldName),
			Reason:    "incompatible field types: " + sType.String() + " -> " + dType.String(),
		}
	}

	// Slow path: complex types that need recursion - build path once here
	fullPath := buildPath(basePath, fieldName)

	if srcKind == reflect.Struct && dstKind == reflect.Struct {
		return assignStruct(dst, src, srcStructType, dstStructType, fullPath, tagName, depth-1)
	}

	if srcKind == reflect.Slice && dstKind == reflect.Slice {
		return assignSlice(dst, src, srcStructType, dstStructType, fullPath, tagName, depth-1)
	}

	if srcKind == reflect.Map && dstKind == reflect.Map {
		return assignMap(dst, src, srcStructType, dstStructType, fullPath, tagName, depth-1)
	}

	if srcKind == reflect.Ptr && dstKind == reflect.Ptr {
		if src.IsNil() {
			dst.Set(reflect.Zero(dType))
			return nil
		}
		newPtr := reflect.New(dType.Elem())
		if err := assignNestedValue(newPtr.Elem(), src.Elem(), srcStructType, dstStructType, fullPath, "", tagName, convertTo, depth-1); err != nil {
			return err
		}
		dst.Set(newPtr)
		return nil
	}

	if srcKind == reflect.Ptr && dstKind != reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		return assignNestedValue(dst, src.Elem(), srcStructType, dstStructType, fullPath, "", tagName, convertTo, depth-1)
	}

	if srcKind != reflect.Ptr && dstKind == reflect.Ptr {
		newPtr := reflect.New(dType.Elem())
		if err := assignNestedValue(newPtr.Elem(), src, srcStructType, dstStructType, fullPath, "", tagName, convertTo, depth-1); err != nil {
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
		FieldPath: fullPath,
		Reason:    "incompatible field types: " + sType.String() + " -> " + dType.String(),
	}
}
