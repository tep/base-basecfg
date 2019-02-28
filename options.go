package basecfg

import "io"

type Option interface {
	setopt(*cfgOptions)
}

type cfgOptions struct {
	base      Feature
	ocfErr    onConfFileErr
	envPrefix string
	cfgReader *cfgReader
}

//--------------------------------------

func Base(f Feature) Option {
	return &baseOpt{f}
}

type baseOpt struct {
	f Feature
}

func (b *baseOpt) setopt(c *cfgOptions) {
	c.base = b.f
}

//--------------------------------------

const (
	RequireConfigFile onConfFileErr = iota
	WarnOnConfigFileErrors
	IgnoreConfigFileErrors
)

type onConfFileErr int

func (o onConfFileErr) setopt(c *cfgOptions) {
	c.ocfErr = o
}

//--------------------------------------

func EnvPrefix(s string) Option {
	return envPrefix(s)
}

type envPrefix string

func (e envPrefix) setopt(c *cfgOptions) {
	c.envPrefix = string(e)
}

//--------------------------------------

func FromReader(typ string, readr io.Reader) Option {
	return &cfgReader{typ, readr}
}

type cfgReader struct {
	typ   string
	readr io.Reader
}

func (r *cfgReader) setopt(c *cfgOptions) {
	c.cfgReader = r
}
