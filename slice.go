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
