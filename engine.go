package mapper

import (
	"reflect"
	"strconv"
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

// assignSlice handles slice assignment with proper deep copying and type conversion.
// It ensures that:
// - nil slices remain nil
// - empty slices remain empty (not nil)
// - a new underlying array is created (modifications to source don't affect destination)
// - element types are converted if compatible
func assignSlice(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath string) error {
	// Handle nil slice: result should be nil.
	if src.IsNil() {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	sType := src.Type()
	dType := dst.Type()
	srcElemType := sType.Elem()
	dstElemType := dType.Elem()

	length := src.Len()

	// Create a new slice with the same length and capacity.
	newSlice := reflect.MakeSlice(dType, length, length)

	// Check if elements are directly assignable or need conversion.
	elementsAssignable := srcElemType.AssignableTo(dstElemType)
	elementsConvertible := srcElemType.ConvertibleTo(dstElemType)

	if !elementsAssignable && !elementsConvertible {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "slice element types are incompatible: " + srcElemType.String() + " -> " + dstElemType.String(),
		}
	}

	// Copy each element.
	for i := 0; i < length; i++ {
		srcElem := src.Index(i)
		dstElem := newSlice.Index(i)

		if elementsAssignable {
			// For same types, we still need to handle nested slices/structs properly.
			if srcElemType.Kind() == reflect.Slice {
				if err := assignSlice(dstElem, srcElem, srcStructType, dstStructType, fieldPath+"["+strconv.Itoa(i)+"]"); err != nil {
					return err
				}
			} else {
				dstElem.Set(srcElem)
			}
		} else if elementsConvertible {
			dstElem.Set(srcElem.Convert(dstElemType))
		}
	}

	dst.Set(newSlice)
	return nil
}

func typeOf(v any) string {
	if v == nil {
		return "<nil>"
	}
	return reflect.TypeOf(v).String()
}

// convertString converts a string to the specified type.
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool.
func convertString(str string, targetType string, srcStructType, dstStructType reflect.Type, fieldPath string) (reflect.Value, error) {
	switch targetType {
	case "int":
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(int(val)), nil

	case "int8":
		val, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(int8(val)), nil

	case "int16":
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(int16(val)), nil

	case "int32":
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(int32(val)), nil

	case "int64":
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(val), nil

	case "uint":
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(uint(val)), nil

	case "uint8":
		val, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(uint8(val)), nil

	case "uint16":
		val, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(uint16(val)), nil

	case "uint32":
		val, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(uint32(val)), nil

	case "uint64":
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(val), nil

	case "float32":
		val, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(float32(val)), nil

	case "float64":
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(val), nil

	case "bool":
		val, err := strconv.ParseBool(str)
		if err != nil {
			return reflect.Value{}, conversionError(str, targetType, err, srcStructType, dstStructType, fieldPath)
		}
		return reflect.ValueOf(val), nil

	default:
		return reflect.Value{}, &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "unsupported mapconv target type: " + targetType,
		}
	}
}

// conversionError creates a MappingError for string conversion failures.
func conversionError(str, targetType string, parseErr error, srcStructType, dstStructType reflect.Type, fieldPath string) error {
	return &MappingError{
		SrcType:   srcStructType.String(),
		DstType:   dstStructType.String(),
		FieldPath: fieldPath,
		Reason:    "cannot convert \"" + str + "\" to " + targetType + ": " + parseErr.Error(),
	}
}
