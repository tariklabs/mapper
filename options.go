package mapper

type config struct {
	tagName          string
	ignoreZeroSource bool
	strictMode       bool
}

func defaultConfig() *config {
	return &config{
		tagName:          "map",
		ignoreZeroSource: false,
		strictMode:       false,
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
