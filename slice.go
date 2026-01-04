package mapper

import (
	"reflect"
	"strconv"
)

// assignSlice handles slice assignment with proper deep copying and type conversion.
// It ensures that:
// - nil slices remain nil
// - empty slices remain empty (not nil)
// - a new underlying array is created (modifications to source don't affect destination)
// - element types are converted if compatible
func assignSlice(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath string) error {
	return assignSliceWithStructs(dst, src, srcStructType, dstStructType, fieldPath, "")
}

// assignSliceWithStructs handles slice assignment with support for nested structs.
// It extends assignSlice to properly map struct elements within slices.
func assignSliceWithStructs(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string) error {
	if src.IsNil() {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	sType := src.Type()
	dType := dst.Type()
	srcElemType := sType.Elem()
	dstElemType := dType.Elem()

	length := src.Len()

	newSlice := reflect.MakeSlice(dType, length, length)

	srcElemKind := srcElemType.Kind()
	dstElemKind := dstElemType.Kind()

	elementsAreStructs := srcElemKind == reflect.Struct && dstElemKind == reflect.Struct
	elementsAreSlices := srcElemKind == reflect.Slice && dstElemKind == reflect.Slice
	elementsAreMaps := srcElemKind == reflect.Map && dstElemKind == reflect.Map
	elementsArePtrs := srcElemKind == reflect.Ptr && dstElemKind == reflect.Ptr
	elementsAssignable := srcElemType.AssignableTo(dstElemType)
	elementsConvertible := srcElemType.ConvertibleTo(dstElemType)

	if !elementsAssignable && !elementsConvertible && !elementsAreStructs && !elementsAreSlices && !elementsAreMaps && !elementsArePtrs {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "slice element types are incompatible: " + srcElemType.String() + " -> " + dstElemType.String(),
		}
	}

	for i := 0; i < length; i++ {
		srcElem := src.Index(i)
		dstElem := newSlice.Index(i)
		elemPath := fieldPath + "[" + strconv.Itoa(i) + "]"

		if elementsAreStructs {
			if err := assignStruct(dstElem, srcElem, srcStructType, dstStructType, elemPath, tagName); err != nil {
				return err
			}
		} else if elementsAreSlices {
			if err := assignSliceWithStructs(dstElem, srcElem, srcStructType, dstStructType, elemPath, tagName); err != nil {
				return err
			}
		} else if elementsAreMaps {
			if err := assignMapWithStructs(dstElem, srcElem, srcStructType, dstStructType, elemPath, tagName); err != nil {
				return err
			}
		} else if elementsArePtrs {
			if err := assignPointerElement(dstElem, srcElem, srcStructType, dstStructType, elemPath, tagName); err != nil {
				return err
			}
		} else if elementsAssignable {
			dstElem.Set(srcElem)
		} else if elementsConvertible {
			dstElem.Set(srcElem.Convert(dstElemType))
		}
	}

	dst.Set(newSlice)
	return nil
}

// assignPointerElement handles pointer elements within slices and maps.
func assignPointerElement(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string) error {
	if src.IsNil() {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	srcElem := src.Elem()
	dstElemType := dst.Type().Elem()

	newPtr := reflect.New(dstElemType)

	srcElemKind := srcElem.Kind()
	dstElemKind := dstElemType.Kind()

	if srcElemKind == reflect.Struct && dstElemKind == reflect.Struct {
		if err := assignStruct(newPtr.Elem(), srcElem, srcStructType, dstStructType, fieldPath, tagName); err != nil {
			return err
		}
	} else if srcElem.Type().AssignableTo(dstElemType) {
		newPtr.Elem().Set(srcElem)
	} else if srcElem.Type().ConvertibleTo(dstElemType) {
		newPtr.Elem().Set(srcElem.Convert(dstElemType))
	} else {
		return &MappingError{
			SrcType:   srcStructType.String(),
			DstType:   dstStructType.String(),
			FieldPath: fieldPath,
			Reason:    "incompatible pointer element types: " + srcElem.Type().String() + " -> " + dstElemType.String(),
		}
	}

	dst.Set(newPtr)
	return nil
}
