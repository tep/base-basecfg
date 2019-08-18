// Copyright 2019 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package basecfg

import (
	"bytes"
	"context"
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

	if err := c.Load(context.Background()); !reflect.DeepEqual(err, tc.err) {
		t.Fatalf("c.Load() == (%v); Wanted (%v)", err, tc.err)
	}

	if got := c.OneOf("otf"); !reflect.DeepEqual(got, tc.want) {
		t.Errorf("c.OneOf(%q) == %#v; Wanted %#v", "otf", got, tc.want)
	}
}

type oneofFeature struct{}

func (f *oneofFeature) FlagSet(*pflag.FlagSet)         {}
func (f *oneofFeature) Validate(context.Context) error { return nil }

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
