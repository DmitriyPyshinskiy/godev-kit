package options

import (
	"testing"
)

type options struct {
	Field1 string
	Field2 int
}

func WithField1(value string) OptionFn[options] {
	return func(o *options) {
		o.Field1 = value
	}
}

func WithField2(value int) OptionFn[options] {
	return func(o *options) {
		o.Field2 = value
	}
}

func defaultOptions() options {
	return options{
		Field1: "test",
		Field2: 234,
	}
}

func TestParse(t *testing.T) {
	// Parse the options
	parsedOptions := Parse[options](WithField1("test"), WithField2(123))

	// Check if the options are parsed correctly
	if parsedOptions.Field1 != "test" || parsedOptions.Field2 != 123 {
		t.Errorf("Parse() = %v; want Field1:test, Field2:123", parsedOptions)
	}
}

func TestParseWithDefaults(t *testing.T) {
	// Parse the options with defaults
	parsedOptions := ParseWithDefaults[options](defaultOptions(), WithField2(456))

	// Check if the options are parsed correctly
	if parsedOptions.Field1 != "test" || parsedOptions.Field2 != 456 {
		t.Errorf("ParseWithDefaults() = %v; want Field1:test, Field2:456", parsedOptions)
	}
}
