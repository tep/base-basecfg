package basecfg

import (
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

func (bc *baseConfig) Validate() error {
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

	if err := bc.Load(); err != nil {
		t.Errorf("bc.Load() == %v; Wanted %v", err, nil)
	}

	t.Logf("port: %d", bc.Port)
}

func TestIgnore(t *testing.T) {
	bc := new(baseConfig)
	bc.Config = New("tb2", IgnoreConfigFileErrors)

	pflag.Parse()

	if err := bc.Load(); err != nil {
		t.Errorf("bc.Load() == %v; Wanted %v", err, nil)
	}
}
