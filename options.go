package mapper

// DefaultMaxDepth is the default maximum nesting depth for struct mapping.
// This limit prevents stack overflow from deeply nested or circular references.
// The default value of 64 is sufficient for most real-world use cases.
const DefaultMaxDepth = 64

type config struct {
	tagName          string
	ignoreZeroSource bool
	strictMode       bool
	maxDepth         int
}

// defaultConfig returns default configuration values.
// Returns a value (not pointer) to enable stack allocation in the caller.
func defaultConfig() config {
	return config{
		tagName:          "map",
		ignoreZeroSource: false,
		strictMode:       false,
		maxDepth:         DefaultMaxDepth,
	}
}

// Option configures the behavior of [MapWithOptions].
// Options are applied in the order they are passed.
type Option func(*config)

// WithTagName sets the struct tag name used to read field aliases from source
// struct fields. The default tag name is "map".
//
// This option is useful when you want to reuse existing struct tags (like "json"
// or "db") for mapping instead of adding separate "map" tags.
//
// Example:
//
//	type APIResponse struct {
//	    UserName string `json:"name"`
//	    UserAge  int    `json:"age"`
//	}
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	err := mapper.MapWithOptions(&user, response, mapper.WithTagName("json"))
//
// With this option, the mapper will look for "json" tags instead of "map" tags
// when determining field aliases.
func WithTagName(tag string) Option {
	return func(c *config) {
		c.tagName = tag
	}
}

// WithIgnoreZeroSource configures the mapper to skip assignments when the
// source field has the zero value for its type. This enables patch semantics
// where only explicitly set fields are copied.
//
// Zero values that are skipped:
//   - Empty string ("")
//   - Zero numbers (0, 0.0)
//   - false for booleans
//   - nil for pointers, slices, and maps
//   - Empty structs
//
// Example - Patch operation:
//
//	type User struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	existing := User{Name: "Rafa", Email: "rafa@old.com", Age: 25}
//	patch := User{Name: "Alicia", Email: "", Age: 0}  // Only Name should update
//
//	err := mapper.MapWithOptions(&existing, patch, mapper.WithIgnoreZeroSource())
//	// existing = {Name: "Alicia", Email: "rafa@old.com", Age: 25}
//
// Without this option, the empty string and zero would overwrite the existing values.
func WithIgnoreZeroSource() Option {
	return func(c *config) {
		c.ignoreZeroSource = true
	}
}

// WithStrictMode configures the mapper to return an error when a destination
// field has no matching source field (by name or tag).
//
// This is useful for ensuring that all destination fields are populated,
// catching mistakes like typos in field names or missing fields in the source.
//
// Example:
//
//	type Source struct {
//	    Name string
//	}
//
//	type Destination struct {
//	    Name  string
//	    Email string  // No matching source field
//	}
//
//	err := mapper.MapWithOptions(&dst, src, mapper.WithStrictMode())
//	// Returns error: no matching source field found for "Email"
//
// Without strict mode, unmatched destination fields are silently left unchanged.
func WithStrictMode() Option {
	return func(c *config) {
		c.strictMode = true
	}
}

// WithMaxDepth sets the maximum nesting depth for struct mapping. This prevents
// stack overflow from deeply nested structures or circular references.
//
// The default depth is [DefaultMaxDepth] (64), which is sufficient for most
// use cases. Each level of nesting (nested struct, slice element, map value)
// decrements the depth counter.
//
// Values less than or equal to 0 are ignored, and the default depth is used.
//
// Example:
//
//	// For very deeply nested structures
//	err := mapper.MapWithOptions(&dst, src, mapper.WithMaxDepth(200))
//
//	// For strict depth limiting
//	err := mapper.MapWithOptions(&dst, src, mapper.WithMaxDepth(10))
//	// Returns error if nesting exceeds 10 levels
//
// When the depth limit is exceeded, the mapper returns a [*MappingError] with
// the reason "maximum nesting depth exceeded".
func WithMaxDepth(depth int) Option {
	return func(c *config) {
		if depth > 0 {
			c.maxDepth = depth
		}
	}
}
