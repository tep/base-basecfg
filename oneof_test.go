package basecfg

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

type ooTestcase struct {
	text string
	want Feature
	err  error
}

func TestOneOf(t *testing.T) {
	fa := `"feata": { "option": "value1" }`
	fb := `"featb": { "option": "value1" }`
	fc := `"featc": { "other": "thing"}`
	bf := fa + "," + fb

	tests := map[string]*ooTestcase{
		"a-want-a":   &ooTestcase{fa, &oneofFeatureA{Option: "value1"}, nil},
		"b-want-b":   &ooTestcase{fb, &oneofFeatureB{Option: "value1"}, nil},
		"c-want-nil": &ooTestcase{fc, nil, nil},
		"ab-err":     &ooTestcase{bf, nil, multipleOneOfError("otf", "feata", "featb")},
	}

	for name, tc := range tests {
		t.Run(name, tc.test)
	}
}

func (tc *ooTestcase) test(t *testing.T) {
	reset := useTestRegistry()
	defer reset()

	RegisterOneOf("otf", "feata", func() Feature { return new(oneofFeatureA) })
	RegisterOneOf("otf", "featb", func() Feature { return new(oneofFeatureB) })

	buf := bytes.NewBufferString("{" + tc.text + "}")

	c := New("oneoftest", FromReader("json", buf))

	if err := c.Load(); !reflect.DeepEqual(err, tc.err) {
		t.Fatalf("c.Load() == (%v); Wanted (%v)", err, tc.err)
	}

	if got := c.OneOf("otf"); !reflect.DeepEqual(got, tc.want) {
		t.Errorf("c.OneOf(%q) == %#v; Wanted %#v", "otf", got, tc.want)
	}
}

type oneofFeature struct{}

func (f *oneofFeature) FlagSet(*pflag.FlagSet) {}
func (f *oneofFeature) Validate() error        { return nil }

type oneofFeatureA struct {
	Option string `cfg:"option"`
	oneofFeature
}

type oneofFeatureB struct {
	Option string `cfg:"option"`
	oneofFeature
}

type oneofFeatureC struct {
	Other string `cfg:"other"`
	oneofFeature
}
