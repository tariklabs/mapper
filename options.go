package mapper

// DefaultMaxDepth is the default maximum nesting depth for struct mapping.
// This prevents stack overflow from circular references.
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

type Option func(*config)

// WithTagName sets the struct tag name used to read mapping hints
// from source struct fields. Default is "map".
func WithTagName(tag string) Option {
	return func(c *config) {
		c.tagName = tag
	}
}

// WithIgnoreZeroSource makes the mapper skip assignments where the
// source field has the zero value for its type.
func WithIgnoreZeroSource() Option {
	return func(c *config) {
		c.ignoreZeroSource = true
	}
}

// WithStrictMode makes the mapper return an error when a destination
// field has no matching source field (by name or tag).
func WithStrictMode() Option {
	return func(c *config) {
		c.strictMode = true
	}
}

// WithMaxDepth sets the maximum nesting depth for struct mapping.
// This prevents stack overflow from circular references.
// Default is 64. Values less than or equal to 0 are ignored and
// the default depth limit will be used.
func WithMaxDepth(depth int) Option {
	return func(c *config) {
		if depth > 0 {
			c.maxDepth = depth
		}
	}
}
