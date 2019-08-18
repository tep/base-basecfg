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
	"os"
	"testing"

	// "github.com/kr/pretty"
	"github.com/spf13/pflag"
)

type baseConfig struct {
	Name string `cfg:"name"`
	Port uint32 `cfg:"port"`

	*Config
}

func (bc *baseConfig) FlagSet(fs *pflag.FlagSet) {
	fs.StringVar(&bc.Name, "name", bc.Name, "The name")
	fs.Uint32Var(&bc.Port, "port", bc.Port, "The port")
}

func (bc *baseConfig) Validate(context.Context) error {
	return nil
}

func TestConfigBase(t *testing.T) {
	bc := &baseConfig{
		Name: "service1",
		Port: 9991,
	}

	bc.Config = New("testbase", Base(bc))
	bc.AddConfigPath("testdata")

	os.Setenv("XXX_TESTBASE_PORT", "2345")
	os.Args = append(os.Args, "--port", "3456")
	pflag.CommandLine.AddFlagSet(bc.Config.flags)
	pflag.Parse()

	// bc.Config.flags.Parse([]string{"--port", "3456"})

	if err := bc.Load(context.Background()); err != nil {
		t.Errorf("bc.Load() == %v; Wanted %v", err, nil)
	}

	t.Logf("port: %d", bc.Port)
}

func TestIgnore(t *testing.T) {
	bc := new(baseConfig)
	bc.Config = New("tb2", IgnoreConfigFileErrors)

	pflag.Parse()

	if err := bc.Load(context.Background()); err != nil {
		t.Errorf("bc.Load() == %v; Wanted %v", err, nil)
	}
}
