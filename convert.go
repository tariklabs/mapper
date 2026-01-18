package mapper

import (
	"reflect"
	"strconv"
)

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
