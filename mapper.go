package mapper

// Map copies fields from src to dst using default configuration.
//
// The dst argument must be a non-nil pointer to a struct. The src argument
// must be a struct or a non-nil pointer to a struct. Fields are matched by
// name (case-sensitive) or by the "map" tag value on source fields.
//
// Map performs deep copying for slices, maps, and nested structs. Pointer
// fields are handled flexibly: values can map to pointers and vice versa.
//
// Only exported fields are mapped. Unexported fields are silently ignored.
//
// Example:
//
//	type Source struct {
//	    Name  string
//	    Email string `map:"ContactEmail"`
//	}
//
//	type Destination struct {
//	    Name         string
//	    ContactEmail string
//	}
//
//	src := Source{Name: "Rafa", Email: "rafa@example.com"}
//	var dst Destination
//
//	if err := mapper.Map(&dst, src); err != nil {
//	    log.Fatal(err)
//	}
//	// dst.Name = "Rafa", dst.ContactEmail = "rafa@example.com"
//
// Returns a [*MappingError] if mapping fails. Common failure causes include
// type incompatibility, nil pointers, and string conversion errors.
func Map(dst any, src any) error {
	return MapWithOptions(dst, src)
}

// MapWithOptions copies fields from src to dst with custom configuration.
//
// Options are applied in order using the functional options pattern.
// See [WithTagName], [WithIgnoreZeroSource], [WithStrictMode], and
// [WithMaxDepth] for available options.
//
// Example with multiple options:
//
//	err := mapper.MapWithOptions(&dst, src,
//	    mapper.WithTagName("json"),       // Use json tags instead of map tags
//	    mapper.WithIgnoreZeroSource(),    // Skip zero-value fields
//	    mapper.WithStrictMode(),          // Error on unmapped destination fields
//	)
//
// Example for patch operations:
//
//	// Only update fields that have non-zero values in the patch
//	existing := User{Name: "Rafa", Age: 30}
//	patch := UpdateRequest{Name: "Alicia"}  // Age is zero, will be skipped
//
//	err := mapper.MapWithOptions(&existing, patch, mapper.WithIgnoreZeroSource())
//	// existing.Name = "Alicia", existing.Age = 30 (unchanged)
//
// Returns a [*MappingError] if mapping fails.
func MapWithOptions(dst any, src any, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return runMapping(dst, src, &cfg)
}
