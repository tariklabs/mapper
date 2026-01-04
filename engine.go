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

		if err := assignValue(dstField, srcField, srcType, dstType, dstName, srcFieldMeta.ConvertTo, cfg.tagName, cfg.maxDepth); err != nil {
			return err
		}
	}

	return nil
}

// assignValue tries to assign src to dst, handling basic cases and pointer/value combinations.
func assignValue(dst, src reflect.Value, srcType, dstType reflect.Type, fieldPath string, convertTo string, tagName string, depth int) error {
	if depth <= 0 {
		return &MappingError{
			SrcType:   srcType.String(),
			DstType:   dstType.String(),
			FieldPath: fieldPath,
			Reason:    "maximum nesting depth exceeded (possible circular reference)",
		}
	}

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

	if convertTo != "" && sType.Kind() == reflect.String {
		converted, err := convertString(src.String(), convertTo, srcType, dstType, fieldPath)
		if err != nil {
			return err
		}
		dst.Set(converted.Convert(dType))
		return nil
	}

	srcKind := sType.Kind()
	dstKind := dType.Kind()

	if srcKind == reflect.Struct && dstKind == reflect.Struct {
		return assignStruct(dst, src, srcType, dstType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Slice && dstKind == reflect.Slice {
		return assignSlice(dst, src, srcType, dstType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Map && dstKind == reflect.Map {
		return assignMap(dst, src, srcType, dstType, fieldPath, tagName, depth-1)
	}

	if srcKind == reflect.Ptr && dstKind == reflect.Ptr {
		if src.IsNil() {
			dst.Set(reflect.Zero(dType))
			return nil
		}

		srcElemType := sType.Elem()
		dstElemType := dType.Elem()

		if srcElemType == dstElemType {
			elemKind := srcElemType.Kind()
			if elemKind != reflect.Struct && elemKind != reflect.Slice && elemKind != reflect.Map && elemKind != reflect.Ptr {
				newPtr := reflect.New(dstElemType)
				newPtr.Elem().Set(src.Elem())
				dst.Set(newPtr)
				return nil
			}
		}

		newPtr := reflect.New(dstElemType)
		if err := assignValue(newPtr.Elem(), src.Elem(), srcType, dstType, fieldPath, convertTo, tagName, depth-1); err != nil {
			return err
		}
		dst.Set(newPtr)
		return nil
	}

	if srcKind == reflect.Ptr && dstKind != reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		return assignValue(dst, src.Elem(), srcType, dstType, fieldPath, convertTo, tagName, depth-1)
	}

	if srcKind != reflect.Ptr && dstKind == reflect.Ptr {
		newVal := reflect.New(dType.Elem())
		if err := assignValue(newVal.Elem(), src, srcType, dstType, fieldPath, convertTo, tagName, depth-1); err != nil {
			return err
		}
		dst.Set(newVal)
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
