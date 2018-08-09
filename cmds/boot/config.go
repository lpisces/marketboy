package boot

import (
	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
)

var (
	Conf *Config
)

type (
	Config struct {
		Debug      bool
		ConfigFile string
		*WSConfig
		*RestConfig
		*AuthConfig
		*Subscribe
		*Trading
	}

	WSConfig struct {
		Scheme string
		Host   string
		Path   string
	}

	RestConfig struct {
		Scheme string
		Host   string
		Prefix string
		Port   string
	}

	AuthConfig struct {
		Key    string
		Secret string
	}

	Subscribe struct {
		Topic []string
	}

	Trading struct {
		UnitQty    float64
		MaxHoldQty float64
		Symbol     []string
		Spread     float64
		PriceUint  float64
		Range      int64
	}
)

func init() {
	Conf = Default()
}

// DefaultConfig get default config
func Default() *Config {

	return &Config{
		false,
		"config.ini",
		&WSConfig{
			"ws",
			"localhost",
			"/",
		},
		&RestConfig{
			"https",
			"localhost",
			"/",
			"",
		},
		&AuthConfig{
			"",
			"",
		},
		&Subscribe{
			[]string{},
		},
		&Trading{
			100,
			1000,
			[]string{"XBTUSD"},
			1,
			0.5,
			5,
		},
	}

}

// LoadFromIni load config from ini override default config
func (config *Config) LoadFromIni() (err error) {
	return ini.MapTo(config, config.ConfigFile)
}

// Load load config from command line param
func (config *Config) Load(c *cli.Context) (err error) {

	if c.String("config") != "" {
		Conf.ConfigFile = c.String("config")
		if err = Conf.LoadFromIni(); err != nil {
			return
		}
	}

	if c.Bool("debug") {
		Conf.Debug = true
	}

	return
}
