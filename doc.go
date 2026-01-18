// Package mapper provides a zero-boilerplate, reflection-based struct mapper for Go.
// It automatically maps fields between structs using name matching, tag aliasing,
// and type conversion, with support for nested structs, slices, maps, and pointers.
//
// # Quick Start
//
// Basic mapping between structs with matching field names:
//
//	type Source struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	type Destination struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	src := Source{Name: "Alice", Email: "alice@example.com", Age: 30}
//	var dst Destination
//
//	err := mapper.Map(&dst, src)
//	// dst = {Name: "Alice", Email: "alice@example.com", Age: 30}
//
// # Tag-Based Field Aliasing
//
// Use the "map" tag to map fields with different names:
//
//	type APIResponse struct {
//	    UserName    string `map:"Name"`
//	    UserEmail   string `map:"Email"`
//	    YearsOld    int    `map:"Age"`
//	}
//
//	type User struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	src := APIResponse{UserName: "Bob", UserEmail: "bob@example.com", YearsOld: 25}
//	var dst User
//
//	err := mapper.Map(&dst, src)
//	// dst = {Name: "Bob", Email: "bob@example.com", Age: 25}
//
// # String-to-Type Conversion
//
// Use the "mapconv" tag to convert string fields to numeric or boolean types:
//
//	type FormInput struct {
//	    Age      string `mapconv:"int"`
//	    Score    string `mapconv:"float64"`
//	    Active   string `mapconv:"bool"`
//	}
//
//	type User struct {
//	    Age    int
//	    Score  float64
//	    Active bool
//	}
//
//	src := FormInput{Age: "42", Score: "95.5", Active: "true"}
//	var dst User
//
//	err := mapper.Map(&dst, src)
//	// dst = {Age: 42, Score: 95.5, Active: true}
//
// Supported conversion types: int, int8, int16, int32, int64, uint, uint8,
// uint16, uint32, uint64, float32, float64, bool.
//
// Tags can be combined for aliasing with conversion:
//
//	type Input struct {
//	    UserAge string `map:"Age" mapconv:"int"`
//	}
//
// # Nested Struct Mapping
//
// Nested structs are mapped recursively:
//
//	type SrcAddress struct {
//	    Street string
//	    City   string
//	}
//
//	type SrcPerson struct {
//	    Name    string
//	    Address SrcAddress
//	}
//
//	type DstAddress struct {
//	    Street string
//	    City   string
//	}
//
//	type DstPerson struct {
//	    Name    string
//	    Address DstAddress
//	}
//
//	src := SrcPerson{
//	    Name: "Alice",
//	    Address: SrcAddress{Street: "123 Main St", City: "Seattle"},
//	}
//	var dst DstPerson
//
//	err := mapper.Map(&dst, src)
//	// dst.Address = {Street: "123 Main St", City: "Seattle"}
//
// # Slice and Map Deep Copying
//
// Slices and maps are deep-copied to avoid shared references:
//
//	type Source struct {
//	    Tags   []string
//	    Config map[string]string
//	}
//
//	type Destination struct {
//	    Tags   []string
//	    Config map[string]string
//	}
//
//	src := Source{
//	    Tags:   []string{"go", "mapper"},
//	    Config: map[string]string{"env": "prod"},
//	}
//	var dst Destination
//
//	err := mapper.Map(&dst, src)
//	// Modifying src.Tags after mapping does not affect dst.Tags
//
// Element type conversion is supported for slices and maps:
//
//	type Source struct {
//	    Values []int32
//	}
//
//	type Destination struct {
//	    Values []int64  // Different element type
//	}
//
// # Pointer Handling
//
// Flexible conversion between pointer and value types:
//
//	type Source struct {
//	    Name string   // value
//	}
//
//	type Destination struct {
//	    Name *string  // pointer
//	}
//
//	src := Source{Name: "Alice"}
//	var dst Destination
//
//	err := mapper.Map(&dst, src)
//	// *dst.Name == "Alice"
//
// Nil pointers are handled gracefully and do not overwrite destination values.
//
// # Options
//
// Use [MapWithOptions] for customized behavior:
//
//	// Use a different tag name
//	err := mapper.MapWithOptions(&dst, src, mapper.WithTagName("json"))
//
//	// Skip zero-value fields (patch semantics)
//	err := mapper.MapWithOptions(&dst, src, mapper.WithIgnoreZeroSource())
//
//	// Error on missing source fields
//	err := mapper.MapWithOptions(&dst, src, mapper.WithStrictMode())
//
//	// Increase max depth for deeply nested structs
//	err := mapper.MapWithOptions(&dst, src, mapper.WithMaxDepth(100))
//
//	// Combine multiple options
//	err := mapper.MapWithOptions(&dst, src,
//	    mapper.WithTagName("json"),
//	    mapper.WithIgnoreZeroSource(),
//	    mapper.WithStrictMode(),
//	)
//
// # Patch Semantics
//
// Use [WithIgnoreZeroSource] for partial updates where only non-zero values
// should be applied:
//
//	existing := User{Name: "Alice", Email: "alice@old.com", Age: 25}
//	patch := PatchRequest{Name: "Alicia", Email: "", Age: 0}  // Only update name
//
//	err := mapper.MapWithOptions(&existing, patch, mapper.WithIgnoreZeroSource())
//	// existing = {Name: "Alicia", Email: "alice@old.com", Age: 25}
//
// # Error Handling
//
// Errors are returned as [*MappingError] with detailed context:
//
//	err := mapper.Map(&dst, src)
//	if err != nil {
//	    var mappingErr *mapper.MappingError
//	    if errors.As(err, &mappingErr) {
//	        fmt.Printf("Error at field %q: %s\n", mappingErr.FieldPath, mappingErr.Reason)
//	    }
//	}
//
// # Performance
//
// The mapper uses struct metadata caching to minimize reflection overhead.
// First-time mapping of a struct type incurs reflection cost, but subsequent
// mappings use cached metadata for faster execution.
//
// For performance-critical code paths where nanoseconds matter, consider
// manual field assignment. For typical application code (API handlers, DTOs),
// the mapper provides a good balance of convenience and performance.
//
// # Thread Safety
//
// All functions are safe for concurrent use. The internal metadata cache uses
// sync.Map for thread-safe access.
//
// # Limitations
//
//   - Only exported (public) fields are mapped
//   - Interface types are not supported as field types
//   - No custom converter functions (only built-in type conversions)
//   - Circular references are protected by depth limit, not runtime detection
package mapper
