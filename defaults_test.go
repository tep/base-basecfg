package basecfg

import (
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

	if err := c.Load(); err != nil {
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

	if err := c.Load(); err != nil {
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
