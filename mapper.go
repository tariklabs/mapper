package mapper

func Map(dst any, src any) error {
	return MapWithOptions(dst, src)
}

func MapWithOptions(dst any, src any, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return runMapping(dst, src, cfg)
}
