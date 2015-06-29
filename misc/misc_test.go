package misc

import (
	"reflect"
	"testing"
)

type ToFloatTest struct {
	in  string
	out float64
}

var ToFloatTests = []ToFloatTest{
	{"1.00", 1.00},
	{"100", 100.00},
	{"123.123457843274", 123.123457843274},
}

func TestToFloat(t *testing.T) {
	for i, test := range ToFloatTests {
		output := ToFloat(test.in)
		if output != test.out {
			t.Errorf("#%d: Input %s; want %f, got %f", i, test.in, test.out, output)
		}
	}
}

func TestUrandom(t *testing.T) {
	input := 123
	output := Urandom(input)

	typ := reflect.TypeOf(output)
	if typ.Kind() != reflect.Slice {
		t.Errorf("Not a slice: %s\n", typ.Kind())
	}

}
