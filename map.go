package mapper

import (
	"reflect"
)

// assignMap handles map assignment with proper deep copying and type conversion.
// It ensures that:
// - nil maps remain nil
// - empty maps remain empty (not nil)
// - a new underlying map is created (modifications to source don't affect destination)
// - key and value types are converted if compatible
func assignMap(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath string) error {
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

	valuesAssignable := srcValType.AssignableTo(dstValType)
	valuesConvertible := srcValType.ConvertibleTo(dstValType)
	valuesAreNestedMaps := srcValType.Kind() == reflect.Map && dstValType.Kind() == reflect.Map
	valuesAreNestedSlices := srcValType.Kind() == reflect.Slice && dstValType.Kind() == reflect.Slice

	if !valuesAssignable && !valuesConvertible && !valuesAreNestedMaps && !valuesAreNestedSlices {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "map value types are incompatible: " + srcValType.String() + " -> " + dstValType.String(),
		}
	}

	newMap := reflect.MakeMapWithSize(dType, src.Len())

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
		if valuesAreNestedMaps {
			dstVal = reflect.New(dstValType).Elem()
			if err := assignMap(dstVal, srcVal, srcStructType, dstStructType, fieldPath+"[key]"); err != nil {
				return err
			}
		} else if valuesAreNestedSlices {
			dstVal = reflect.New(dstValType).Elem()
			if err := assignSlice(dstVal, srcVal, srcStructType, dstStructType, fieldPath+"[key]"); err != nil {
				return err
			}
		} else if srcValType.Kind() == reflect.Ptr {
			if srcVal.IsNil() {
				dstVal = reflect.Zero(dstValType)
			} else {
				newPtr := reflect.New(dstValType.Elem())
				newPtr.Elem().Set(srcVal.Elem())
				dstVal = newPtr
			}
		} else if valuesAssignable {
			dstVal = srcVal
		} else {
			dstVal = srcVal.Convert(dstValType)
		}

		newMap.SetMapIndex(dstKey, dstVal)
	}

	dst.Set(newMap)
	return nil
}
