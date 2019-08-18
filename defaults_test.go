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
	"context"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/spf13/pflag"
)

func TestDefaults(t *testing.T) {
	reset := useTestRegistry()
	defer reset()

	fd := &featureDefn{Feature: mkTestFeature()}

	want := map[string]interface{}{
		"thing-one": "foo",
		"other":     int64(12),
	}

	fd.extractDefaults()

	if got := fd.defaults; !reflect.DeepEqual(got, want) {
		t.Errorf("%#vextractDefaults() == %s; wanted %s", fd, pretty.Sprint(got), pretty.Sprint(want))
	}
}

func TestFlagDefaults(t *testing.T) {
	reset := useTestRegistry()
	defer reset()

	tf := mkTestFeature()

	c := New("test", Base(tf), IgnoreConfigFileErrors)

	if err := c.Load(context.Background()); err != nil {
		t.Fatal(err)
	}

	got := make(map[string]string)
	tf.fs.VisitAll(func(f *pflag.Flag) {
		got[f.Name] = f.DefValue
	})

	want := map[string]string{
		"thing-one": "foo",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect flag default values: Got(%# v); Wanted(%# v)", pretty.Formatter(got), pretty.Formatter(want))
	}
}

func TestFeatureFlagDefaults(t *testing.T) {
	reset := useTestRegistry()
	defer reset()

	Register("feat", func() Feature { return mkTestFeature() })

	c := New("test", IgnoreConfigFileErrors)

	if err := c.Load(context.Background()); err != nil {
		t.Fatal(err)
	}

	tf := c.Feature("feat").(*testFeature)

	got := make(map[string]string)
	tf.fs.VisitAll(func(f *pflag.Flag) {
		got[f.Name] = f.DefValue
	})

	want := map[string]string{
		"feat.thing-one": "foo",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect flag default values: Got(%# v); Wanted(%# v)", pretty.Formatter(got), pretty.Formatter(want))
	}
}
