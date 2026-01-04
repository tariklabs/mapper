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
// - nested structs within slices are properly mapped using the provided tagName
func assignSlice(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string, depth int) error {
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
	srcElemType := sType.Elem()
	dstElemType := dType.Elem()

	length := src.Len()
	newSlice := reflect.MakeSlice(dType, length, length)

	// Fast path: identical simple element types can use reflect.Copy
	if srcElemType == dstElemType {
		elemKind := srcElemType.Kind()
		if elemKind != reflect.Struct && elemKind != reflect.Slice && elemKind != reflect.Map && elemKind != reflect.Ptr {
			reflect.Copy(newSlice, src)
			dst.Set(newSlice)
			return nil
		}
	}

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

	// Slow path: need per-element processing
	for i := 0; i < length; i++ {
		srcElem := src.Index(i)
		dstElem := newSlice.Index(i)

		var err error
		if elementsAreStructs {
			err = assignStruct(dstElem, srcElem, srcStructType, dstStructType, fieldPath+"["+strconv.Itoa(i)+"]", tagName, depth-1)
		} else if elementsAreSlices {
			err = assignSlice(dstElem, srcElem, srcStructType, dstStructType, fieldPath+"["+strconv.Itoa(i)+"]", tagName, depth-1)
		} else if elementsAreMaps {
			err = assignMap(dstElem, srcElem, srcStructType, dstStructType, fieldPath+"["+strconv.Itoa(i)+"]", tagName, depth-1)
		} else if elementsArePtrs {
			err = assignPointerElement(dstElem, srcElem, srcStructType, dstStructType, fieldPath+"["+strconv.Itoa(i)+"]", tagName, depth-1)
		} else if elementsAssignable {
			dstElem.Set(srcElem)
		} else if elementsConvertible {
			dstElem.Set(srcElem.Convert(dstElemType))
		}

		if err != nil {
			return err
		}
	}

	dst.Set(newSlice)
	return nil
}

// assignPointerElement handles pointer elements within slices and maps.
func assignPointerElement(dst, src reflect.Value, srcStructType, dstStructType reflect.Type, fieldPath, tagName string, depth int) error {
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

	srcElem := src.Elem()
	dstElemType := dst.Type().Elem()

	newPtr := reflect.New(dstElemType)

	srcElemKind := srcElem.Kind()
	dstElemKind := dstElemType.Kind()

	if srcElemKind == reflect.Struct && dstElemKind == reflect.Struct {
		if err := assignStruct(newPtr.Elem(), srcElem, srcStructType, dstStructType, fieldPath, tagName, depth-1); err != nil {
			return err
		}
	} else if srcElemKind == reflect.Slice && dstElemKind == reflect.Slice {
		if err := assignSlice(newPtr.Elem(), srcElem, srcStructType, dstStructType, fieldPath, tagName, depth-1); err != nil {
			return err
		}
	} else if srcElemKind == reflect.Map && dstElemKind == reflect.Map {
		if err := assignMap(newPtr.Elem(), srcElem, srcStructType, dstStructType, fieldPath, tagName, depth-1); err != nil {
			return err
		}
	} else if srcElemKind == reflect.Ptr && dstElemKind == reflect.Ptr {
		if err := assignPointerElement(newPtr.Elem(), srcElem, srcStructType, dstStructType, fieldPath, tagName, depth-1); err != nil {
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
