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

package basecfg // import "toolman.org/base/basecfg"

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"toolman.org/base/log/v2"
	"toolman.org/base/toolman/v2"
)

type Config struct {
	file  string
	path  []string
	defs  []*featureDefn
	fmap  map[Label]Feature
	oomap map[string]Label
	flags *pflag.FlagSet
	opts  *cfgOptions
	*viper.Viper
}

func New(name string, options ...Option) *Config {
	fsName := filepath.Base(os.Args[0])
	fs := pflag.NewFlagSet(fsName, pflag.ExitOnError)
	fs.SortFlags = false

	opts := &cfgOptions{envPrefix: name}

	for _, o := range options {
		o.setopt(opts)
	}

	v := viper.New()

	v.SetConfigName(name)
	v.SetEnvPrefix(opts.envPrefix)

	c := &Config{
		defs:  registry.reify(),
		fmap:  make(map[Label]Feature),
		oomap: make(map[string]Label),
		flags: fs,
		opts:  opts,
		Viper: v,
	}

	fs.StringVar(&c.file, "config-file", "",
		"If specified, use only this specific config file (i.e. don't search config path)")

	fs.StringSliceVar(&c.path, "config-path", defaultCfgPath(),
		"Comma separated list of directories to search for config (may be specified more than once)")

	// If a base feature is provided, we prepend it to our list of features.
	var bf *featureDefn
	if opts.base != nil {
		bf = &featureDefn{
			Feature: opts.base,
		}

		bf.extractDefaults()

		c.defs = append([]*featureDefn{bf}, c.defs...)
	}

	for _, fd := range c.defs {
		var pfx string
		fsn := fsName
		if bf == nil || fd != bf {
			// This is a non-base feature
			c.fmap[fd.label] = fd.Feature
			fsn = fmt.Sprintf("%s:%s", fsn, fd.label)
			pfx = string(fd.label) + "."
		}

		ffs := pflag.NewFlagSet(fsn, pflag.ExitOnError)

		fd.FlagSet(ffs)

		ffs.VisitAll(func(f *pflag.Flag) {
			// Prepend the feature's label to the flag name
			// (if it's not there already)
			if pfx != "" && !strings.HasPrefix(f.Name, pfx) {
				f.Name = pfx + f.Name
			}

			// Set the flag's default value
			f.DefValue = stringify(fd.defaults[f.Name])
		})

		fs.AddFlagSet(ffs)
	}

	return c
}

func (c *Config) Flags() *toolman.InitOption {
	return toolman.FlagSet(c.flags)
}

func (c *Config) Load() error {
	if c.file != "" {
		log.Infof("Ignoring config path in lieu of: %s", c.file)
		c.SetConfigFile(c.file)
	} else {
		for _, d := range c.path {
			log.Infof("Updating config search path: %q", d)
			c.AddConfigPath(filepath.Clean(d))
		}
	}

	c.AddConfigPath(".")
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	c.AutomaticEnv()

	if err := c.readConfig(); err != nil {
		return err
	}

	c.flags.Visit(func(f *pflag.Flag) {
		// Only bind changed, non-hidden flags
		if f.Changed && !f.Hidden {
			c.BindPFlag(f.Name, f)
		}
	})

	// TODO: Rewrite this comment
	// `all` and `oom` are used to determine whether a particular "oneof" has
	// already been configured. `all` is a map[string]interface{} containing all
	// currently available key/val pairs while `oom` is the "one of map" used for
	// marking previously encountered "oneof" names.
	all := c.AllSettings()

	for _, fd := range c.defs {
		// If a feature's label is found in the `all` map (meaning we have config
		// values for that Feature) *and* that feature has a non-empty "oneof"
		// name, then this is the configured Feature for that "oneof" set.
		// There can be only one of these per "oneof" name.
		if _, has := all[string(fd.label)]; has && fd.oneof != "" {
			if oo, ok := c.oomap[fd.oneof]; ok {
				delete(c.oomap, fd.oneof)
				return multipleOneOfError(fd.oneof, oo, fd.label)
			}
			log.Infof("Using %q as %q", fd.label, fd.oneof)
			c.oomap[fd.oneof] = fd.label
		}

		for k, v := range fd.defaults {
			c.SetDefault(k, v)
		}

		if err := c.unmarshal(fd); err != nil {
			return err
		}

		// We skip the call to Validate for oneof Features that are not currently
		// configured.
		if fd.oneof == "" || c.oomap[fd.oneof] != "" {
			if err := fd.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) OneOf(name string) Feature {
	return c.Feature(c.oomap[name])
}

// TODO(tep): Add a struct tag option to indicate required params.
//
//            1) A mapstructure.DecodeHookFunc will be required to parse the
//               struct tag looking for the `required` option and stash this
//               info into a `map[string]bool`.
//
//            2) After config is loaded (and before calling `Validate()`),
//               each `required` value should be checked with `isZero()` to
//               determine whether it's been set.
//

func isZero(in interface{}) bool {
	v := reflect.ValueOf(in)

	if k := v.Kind(); k == reflect.Interface || k == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		return true
	}

	t := v.Type()

	if t.Comparable() {
		return v.Interface() == reflect.Zero(t).Interface()
	}

	if v.Kind() == reflect.Slice {
		return v.Len() == 0
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
}

func (c *Config) readConfig() error {
	err := c.readConfigData()
	if err == nil {
		log.Infof("Config loaded from: %q", c.ConfigFileUsed())
		return nil
	}

	switch c.opts.ocfErr {
	case RequireConfigFile:
		return err
	case WarnOnConfigFileErrors:
		log.Warningf("Error loading configuration file: %v", err)
		fallthrough
	default:
		return nil
	}
}

func (c *Config) readConfigData() error {
	if cr := c.opts.cfgReader; cr != nil {
		c.SetConfigType(cr.typ)
		return c.ReadConfig(cr.readr)
	}
	return c.ReadInConfig()
}

func tagname(c *mapstructure.DecoderConfig) {
	c.TagName = "cfg"
}

func (c *Config) unmarshal(fd *featureDefn) error {
	// A base Feature will have no create func and uses viper's Unmarshal method.
	if fd.create == nil {
		return c.Unmarshal(fd.Feature, tagname)
	}

	// All others use UnmarshalKey,
	// however...

	// What follows is a kludge to work around viper bug #188 that wreaks
	// havok between ENV overrides and `UnmarshalKey` (amongst other things).
	//
	// The gist is: we call `c.Get` on the viper-key constructed from each of
	// the map keys returned by "Get"ing the raw map for this feature - if
	// this key's `Get` value differs from the raw map's value, we force an
	// override to the former.
	//
	// IOW: `c.Get("foo")["bar"]` might be wrong (it ignores "${FOO_BAR}")
	// but `c.Get("foo.bar")` is correct (i.e. it will be "${FOO_BAR}" if its
	// set) -- so, we compare the two and if they're different we force the
	// value for "foo.bar" to be what's returned by `c.Get("foo.bar")`.
	//
	// Once we do this, `c.UnmarshalKey(...)` will do the right thing.
	//
	ls := string(fd.label)

	if fm, ok := c.Get(ls).(map[string]interface{}); ok {
		for k := range fm {
			fk := fd.label.Key(k)
			if kv := c.Get(fk); kv != fm[fk] {
				c.Set(fk, kv)
			}
		}
	}

	return c.UnmarshalKey(ls, fd.Feature, tagname)
}

func (c *Config) Feature(l Label) Feature {
	return c.fmap[l]
}

func (c *Config) Features() []Label {
	var i int
	list := make([]Label, len(c.fmap))

	for fl := range c.fmap {
		list[i] = fl
		i++
	}

	return list
}

// BinDir is a wrapper around filepath.Join returning a file path constructed
// by joining the directory location of the current program's binary with the
// provided subdirectory components.
func BinDir(subdirs ...string) string {
	return filepath.Clean(filepath.Join(append([]string{filepath.Dir(os.Args[0])}, subdirs...)...))
}

func defaultCfgPath() []string {
	d := BinDir("conf")
	return []string{filepath.Dir(d), d}
}

func stringify(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case *string:
		return *v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
