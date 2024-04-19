// Package options provides a flexible way for applying configuration in options-pattern style
package options

// OptionFn is a generic function type that applies an option to any options of type O.
type OptionFn[O any] func(*O)

// Applier is a generic interface representing any type that can apply options to any options of type O.
type Applier[O any] interface {
	Apply(options *O)
}

// Apply implements the Applier interface by OptionFn
func (f OptionFn[O]) Apply(o *O) {
	f(o)
}

// Parse is a generic options parser
func Parse[O any](opts ...Applier[O]) *O {
	var o O
	for _, opt := range opts {
		opt.Apply(&o)
	}
	return &o
}

// ParseWithDefaults is a generic options parser with defaults
func ParseWithDefaults[O any](defaults O, opts ...Applier[O]) *O {
	o := defaults
	for _, opt := range opts {
		opt.Apply(&o)
	}
	return &o
}
