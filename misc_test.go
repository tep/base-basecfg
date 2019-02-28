package basecfg

import "github.com/spf13/pflag"

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

func (tf *testFeature) Validate() error { return nil }

func mkTestFeature() *testFeature {
	return &testFeature{
		ThingOne: "foo",
		Other:    12,
		Stuff:    "bar",
	}
}
