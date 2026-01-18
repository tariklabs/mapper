package mapper

import "fmt"

// MappingError describes a failure that occurred during struct mapping.
// It provides detailed context about what went wrong and where.
//
// The error message format is:
//
//	mapper: cannot map {SrcType} → {DstType} at field "{FieldPath}": {Reason}
//
// FieldPath uses dot notation for nested fields (e.g., "Address.City") and
// bracket notation for slice indices (e.g., "Items[0].Name") and map keys
// (e.g., "Config[database].Host").
//
// Common Reasons:
//   - "nil src or dst" - nil was passed for source or destination
//   - "dst must be a non-nil pointer to struct" - destination is not a valid pointer
//   - "src must be a struct or pointer to struct" - source is not a struct type
//   - "src is a nil pointer" - source pointer is nil
//   - "no matching source field found" - strict mode enabled, field has no match
//   - "incompatible field types: X -> Y" - types cannot be converted
//   - "maximum nesting depth exceeded" - depth limit reached
//   - "cannot convert \"X\" to Y" - string conversion failed
//   - "unsupported mapconv target type: X" - invalid mapconv tag value
//   - "destination field cannot be set" - field is unexported
//
// Example - Error handling:
//
//	err := mapper.Map(&dst, src)
//	if err != nil {
//	    var mappingErr *mapper.MappingError
//	    if errors.As(err, &mappingErr) {
//	        log.Printf("Mapping failed at %s: %s", mappingErr.FieldPath, mappingErr.Reason)
//	    }
//	}
type MappingError struct {
	// SrcType is the name of the source struct type (e.g., "main.UserDTO").
	SrcType string

	// DstType is the name of the destination struct type (e.g., "main.User").
	DstType string

	// FieldPath is the path to the field where the error occurred.
	// Uses dot notation for nested fields ("Address.City") and brackets
	// for indices ("Items[0]") and map keys ("Config[key]").
	FieldPath string

	// Reason describes why the mapping failed.
	Reason string
}

// Error implements the error interface and returns a formatted error message.
func (e *MappingError) Error() string {
	return fmt.Sprintf(
		"mapper: cannot map %s → %s at field %q: %s",
		e.SrcType, e.DstType, e.FieldPath, e.Reason,
	)
}
