package mapper

import (
	"fmt"
	"reflect"
)

// assignMap handles map assignment with proper deep copying and type conversion.
// It ensures that:
// - nil maps remain nil
// - empty maps remain empty (not nil)
// - a new underlying map is created (modifications to source don't affect destination)
// - key and value types are converted if compatible
func assignMap(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath string) error {
	return assignMapWithStructs(dst, src, srcStructType, dstStructType, fieldPath, "")
}

// assignMapWithStructs handles map assignment with support for nested structs.
// It extends assignMap to properly map struct values within maps.
func assignMapWithStructs(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string) error {
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

		if !needsProcessing && valuesAssignable {
			dstVal = srcVal
		} else if valuesAreStructs {
			dstVal = reflect.New(dstValType).Elem()
			valPath := fieldPath + "[" + fmt.Sprint(srcKey.Interface()) + "]"
			if err := assignStruct(dstVal, srcVal, srcStructType, dstStructType, valPath, tagName); err != nil {
				return err
			}
		} else if valuesAreNestedMaps {
			dstVal = reflect.New(dstValType).Elem()
			valPath := fieldPath + "[" + fmt.Sprint(srcKey.Interface()) + "]"
			if err := assignMapWithStructs(dstVal, srcVal, srcStructType, dstStructType, valPath, tagName); err != nil {
				return err
			}
		} else if valuesAreNestedSlices {
			dstVal = reflect.New(dstValType).Elem()
			valPath := fieldPath + "[" + fmt.Sprint(srcKey.Interface()) + "]"
			if err := assignSliceWithStructs(dstVal, srcVal, srcStructType, dstStructType, valPath, tagName); err != nil {
				return err
			}
		} else if valuesArePtrs {
			dstVal = reflect.New(dstValType).Elem()
			valPath := fieldPath + "[" + fmt.Sprint(srcKey.Interface()) + "]"
			if err := assignPointerElement(dstVal, srcVal, srcStructType, dstStructType, valPath, tagName); err != nil {
				return err
			}
		} else if valuesConvertible {
			dstVal = srcVal.Convert(dstValType)
		}

		newMap.SetMapIndex(dstKey, dstVal)
	}

	dst.Set(newMap)
	return nil
}
