package mapper

import (
	"reflect"
)

func runMapping(dst any, src any, cfg *config) error {
	if dst == nil || src == nil {
		return &MappingError{
			SrcType:   typeOf(src),
			DstType:   typeOf(dst),
			FieldPath: "",
			Reason:    "nil src or dst",
		}
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return &MappingError{
			SrcType:   typeOf(src),
			DstType:   typeOf(dst),
			FieldPath: "",
			Reason:    "dst must be a non-nil pointer to struct",
		}
	}

	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		return &MappingError{
			SrcType:   typeOf(src),
			DstType:   typeOf(dst),
			FieldPath: "",
			Reason:    "dst must point to a struct",
		}
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return &MappingError{
				SrcType:   typeOf(src),
				DstType:   typeOf(dst),
				FieldPath: "",
				Reason:    "src is a nil pointer",
			}
		}
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return &MappingError{
			SrcType:   typeOf(src),
			DstType:   typeOf(dst),
			FieldPath: "",
			Reason:    "src must be a struct or pointer to struct",
		}
	}

	srcType := srcVal.Type()
	dstType := dstElem.Type()

	srcMeta, err := getStructMeta(srcType, cfg.tagName)
	if err != nil {
		return err
	}
	dstMeta, err := getStructMeta(dstType, cfg.tagName)
	if err != nil {
		return err
	}

	for dstName, dstFieldMeta := range dstMeta.FieldsByName {
		srcFieldMeta, ok := srcMeta.FieldsByName[dstName]

		if !ok {
			srcFieldMeta, ok = srcMeta.FieldsByTag[dstName]
		}

		if !ok {
			if cfg.strictMode {
				return &MappingError{
					SrcType:   srcType.String(),
					DstType:   dstType.String(),
					FieldPath: dstName,
					Reason:    "no matching source field found",
				}
			}
			continue
		}

		srcField := srcVal.FieldByIndex(srcFieldMeta.Index)
		dstField := dstElem.FieldByIndex(dstFieldMeta.Index)

		if cfg.ignoreZeroSource && srcField.IsZero() {
			continue
		}

		if err := assignValue(dstField, srcField, srcType, dstType, dstName, srcFieldMeta.ConvertTo); err != nil {
			return err
		}
	}

	return nil
}

// assignValue tries to assign src to dst, handling basic cases and pointer/value combinations.
func assignValue(dst, src reflect.Value, srcType, dstType reflect.Type, fieldPath string, convertTo string) error {
	if !dst.CanSet() {
		return &MappingError{
			SrcType:   srcType.String(),
			DstType:   dstType.String(),
			FieldPath: fieldPath,
			Reason:    "destination field cannot be set",
		}
	}

	sType := src.Type()
	dType := dst.Type()

	// String conversion using mapconv tag.
	if convertTo != "" && sType.Kind() == reflect.String {
		converted, err := convertString(src.String(), convertTo, srcType, dstType, fieldPath)
		if err != nil {
			return err
		}
		dst.Set(converted.Convert(dType))
		return nil
	}

	// Exact or assignable type.
	if sType.AssignableTo(dType) {
		// For slices, create a deep copy to avoid sharing underlying array.
		if sType.Kind() == reflect.Slice {
			if err := assignSlice(dst, src, srcType, dstType, fieldPath); err != nil {
				return err
			}
			return nil
		}
		dst.Set(src)
		return nil
	}

	// Convertible types.
	if sType.ConvertibleTo(dType) {
		dst.Set(src.Convert(dType))
		return nil
	}

	// Slice handling for different but compatible element types.
	if sType.Kind() == reflect.Slice && dType.Kind() == reflect.Slice {
		return assignSlice(dst, src, srcType, dstType, fieldPath)
	}

	// Pointer -> value.
	if sType.Kind() == reflect.Ptr && dType.Kind() != reflect.Ptr {
		if src.IsNil() {
			return nil // nothing to assign
		}
		return assignValue(dst, src.Elem(), srcType, dstType, fieldPath, convertTo)
	}

	// Value -> pointer.
	if sType.Kind() != reflect.Ptr && dType.Kind() == reflect.Ptr {
		// Allocate new value for pointer.
		newVal := reflect.New(dType.Elem())
		if err := assignValue(newVal.Elem(), src, srcType, dstType, fieldPath, convertTo); err != nil {
			return err
		}
		dst.Set(newVal)
		return nil
	}

	return &MappingError{
		SrcType:   srcType.String(),
		DstType:   dstType.String(),
		FieldPath: fieldPath,
		Reason:    "incompatible field types: " + sType.String() + " -> " + dType.String(),
	}
}

func typeOf(v any) string {
	if v == nil {
		return "<nil>"
	}
	return reflect.TypeOf(v).String()
}
