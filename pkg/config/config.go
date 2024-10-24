package config

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var (
	delimiter = "."
	parser    = yaml.Parser()
)

type App struct {
	Debug bool `koanf:"enable_debug"`

	TimeRangeSeconds  int `koanf:"time_range_seconds"`
	UpdateRateSeconds int `koanf:"update_rate_seconds"`

	BarSpacing int `koanf:"bar_spacing"`
	BarWidth   int `koanf:"bar_width"`

	PlotHeight int `koanf:"plot_height"`

	GraphSettings map[string]*GraphSettings `koanf:"graph_settings"`

	Position image.Point `koanf:"position"`

	Theme Theme `koanf:"theme"`
}

type GraphSettings struct {
	Enabled bool `koanf:"enabled"`
}

type Theme struct {
	Window ThemeWindow `koanf:"window"`
	Plot   ThemePlot   `koanf:"plot"`
}
type ThemeWindow struct {
	Active      ThemeWindowType `koanf:"active"`
	Passthrough ThemeWindowType `koanf:"passthrough"`
}

type ThemeWindowType struct {
	Background color.RGBA `koanf:"background"`
}

type ThemePlot struct {
	Border          color.RGBA `koanf:"border"`
	Midline         color.RGBA `koanf:"midline"`
	Bar             color.RGBA `koanf:"bar"`
	LabelText       color.RGBA `koanf:"label_text"`
	LabelBackground color.RGBA `koanf:"label_background"`
}

func NewApp() App {
	return App{
		TimeRangeSeconds:  120,
		UpdateRateSeconds: 1,

		BarSpacing: 0,
		BarWidth:   1,

		PlotHeight: 30,

		Position: image.Pt(10, 10),

		Theme: Theme{
			Window: ThemeWindow{
				Active: ThemeWindowType{
					Background: color.RGBA{20, 20, 0, 20},
				},
				Passthrough: ThemeWindowType{
					Background: color.RGBA{10, 0, 0, 20},
				},
			},
			Plot: ThemePlot{
				Border:          color.RGBA{150, 100, 100, 205},
				Midline:         color.RGBA{0, 200, 0, 100},
				Bar:             color.RGBA{250, 0, 0, 105},
				LabelText:       color.RGBA{255, 255, 255, 180},
				LabelBackground: color.RGBA{0, 0, 0, 0},
			},
		},
	}
}

type Config struct {
	// TODO: use mutex for config
	// m sync.Mutex
	k *koanf.Koanf

	path    string
	Version int `koanf:"version"`

	App App `koanf:"app"`
}

var defaultConfig = Config{
	Version: 1,
	App:     NewApp(),
}

// TODO: add path validation?
func NewConfig(path string) *Config {
	c := &Config{k: koanf.New(delimiter), path: path}

	err := c.k.Load(structs.Provider(&defaultConfig, "koanf"), nil)
	if err != nil {
		panic(err)
	}
	err = c.k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{})
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) Load() error {
	if c.path == "" {
		return errors.New("path to config file cannot be empty")
	}

	err := c.isLoadPossible()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// TODO: Use rawbytes provider to load file only once?
	err = c.k.Load(file.Provider(c.path), parser)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	err = c.k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{})
	if err != nil {
		return fmt.Errorf("cannot parse %s: %w", c.path, err)
	}

	return nil
}

// TODO: backup old config
func (c *Config) Save() error {
	if c.path == "" {
		return errors.New("path to config file cannot be empty")
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"cannot save config '%s': "+
				"file must either not exist, be empty, or "+
				"contain a valid config: %w",
			c.path, err,
		)
	}

	var err error

	err = c.isSavePossible()
	if err != nil {
		return wrapError(err)
	}

	err = c.k.Load(structs.Provider(&c, "koanf"), nil)
	if err != nil {
		panic(err)
	}

	cfgRaw, _ := c.k.Marshal(parser)

	err = os.WriteFile(c.path, cfgRaw, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write config to %s: %w", c.path, err)
	}

	return nil
}

func (c *Config) isVersionMatch() bool {
	return c.Version == defaultConfig.Version
}

func (c *Config) isLoadPossible() error {
	fInfo, err := os.Stat(c.path)
	if err != nil {
		return err
	}
	if fInfo.Size() == 0 {
		return nil
	}

	tmpK := koanf.New(delimiter)

	err = tmpK.Load(file.Provider(c.path), parser)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	tmpCfg := &Config{}
	err = tmpK.UnmarshalWithConf("", tmpCfg, koanf.UnmarshalConf{})
	if err != nil {
		return fmt.Errorf("cannot parse %s: %w", c.path, err)
	}
	// TODO: Implement better check that config is belong to our app
	if !tmpCfg.isVersionMatch() {
		return errors.New("file doesn't look like a valid config")
	}

	return nil
}

// NOTE: We only write to the file if it does not exist, is empty, or contains
// an old configuration, to avoid accidental writes to an unwanted file
func (c *Config) isSavePossible() error {
	fInfo, err := os.Stat(c.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if fInfo.Size() == 0 {
		return nil
	}

	return c.isLoadPossible()
}

func GetDefaultPath() (string, error) {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("can't determine config dir location: %w", err)
	}
	appConfDir := filepath.Join(userConfDir, "govermon")
	return filepath.Join(appConfDir, "govermon.yaml"), nil
}

func InitDefaultFile() (string, error) {
	cfgPath, err := GetDefaultPath()
	if err != nil {
		return "", err
	}
	appConfDir := filepath.Dir(cfgPath)

	err = os.MkdirAll(appConfDir, os.ModeDir)
	if err != nil {
		return "", fmt.Errorf("can't create config dir: %w", err)
	}

	_, err = os.Stat(cfgPath)
	if os.IsNotExist(err) {
		log.Printf(
			"Default config file does not exist. Creating: %s\n", cfgPath,
		)
		err = os.WriteFile(cfgPath, []byte{}, 0o600)
		if err != nil {
			return "", fmt.Errorf("can't create config file: %w", err)
		}
	}
	return cfgPath, nil
}
