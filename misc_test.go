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

	"github.com/spf13/pflag"
)

func useTestRegistry() func() {
	var testRegistry featureRegistry
	orig := registry
	registry = testRegistry

	return func() { registry = orig }
}

type testFeature struct {
	ThingOne string `cfg:"thing-one"`
	Other    int64  `cfg:"other"`
	Stuff    string `cfg:"stuff,nodefault"`
	fs       *pflag.FlagSet
}

func (tf *testFeature) FlagSet(fs *pflag.FlagSet) {
	fs.StringVar(&tf.ThingOne, "thing-one", "", "Thing number one")

	tf.fs = fs
}

func (tf *testFeature) Validate(context.Context) error { return nil }

func mkTestFeature() *testFeature {
	return &testFeature{
		ThingOne: "foo",
		Other:    12,
		Stuff:    "bar",
	}
}
