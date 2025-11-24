package mapper

import "fmt"

// MappingError describes a failure when mapping between two types.
type MappingError struct {
	SrcType   string
	DstType   string
	FieldPath string
	Reason    string
}

// Error implements the error interface.
func (e *MappingError) Error() string {
	return fmt.Sprintf(
		"mapper: cannot map %s → %s at field %q: %s",
		e.SrcType, e.DstType, e.FieldPath, e.Reason,
	)
}
