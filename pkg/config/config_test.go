package config

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"n4/gui-test/internal/utils"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var testDefaultApp = App{
	TimeRangeSeconds:  120,
	UpdateRateSeconds: 1,

	BarSpacing: 0,
	BarWidth:   1,

	PlotHeight: 30,

	GraphSettings: make(map[string]*GraphSettings),

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

func TestNewConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "Emtpy path",
			args: args{""},
			want: &Config{path: "", Version: 1, App: testDefaultApp},
		},
		{
			name: "Non emtpy path",
			args: args{"test"},
			want: &Config{path: "test", Version: 1, App: testDefaultApp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig(tt.args.path)
			require.NotNil(t, got.k)
			require.Equal(t, tt.want.path, got.path)
			require.Equal(t, tt.want.Version, got.Version)
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}

// TODO: fix test
func TestConfig_Load(t *testing.T) {
	t.Run("Empty path", func(t *testing.T) {
		err := (&Config{path: ""}).Load()
		require.ErrorContains(t, err, "path to config file cannot be empty")
	})

	tests := []struct {
		name       string
		createFile bool
		cfgFile    string
		wantErr    string
		want       *Config
	}{
		{
			name:       "File does not exist",
			createFile: false,
			cfgFile:    ``,
			wantErr:    "cannot find the file specified",
		},
		{
			name:       "Empty file",
			createFile: true,
			cfgFile:    ``,
			want: &Config{
				Version: 1,
				App:     testDefaultApp,
			},
		},
		{
			name:       "Invalid version",
			createFile: true,
			cfgFile:    `version: 0`,
			wantErr:    "file doesn't look like a valid config",
		},
		{
			name:       "Invalid YAML",
			createFile: true,
			cfgFile:    `invalid yaml`,
			wantErr:    "error loading config",
		},
		{
			name:       "Invalid config file struct",
			createFile: true,
			cfgFile:    `no_version: there`,
			wantErr:    "file doesn't look like a valid config",
		},
		{
			name:       "Version only",
			createFile: true,
			cfgFile:    `version: 1`,
			want: &Config{
				Version: 1,
				App:     testDefaultApp,
			},
		},
		{
			name:       "Partial config",
			createFile: true,
			cfgFile: `
				version: 1
				app:
					time_range_seconds: 10
				`,
			want: &Config{
				Version: 1,
				App: App{
					TimeRangeSeconds: 10,
					Position:         image.Point{10, 10},
				},
			},
		},
		{
			name:       "Ignore private fields from file",
			createFile: true,
			cfgFile: `
				path: some/path
				version: 1
				app:
					time_range_seconds: 10
				`,
			want: &Config{
				Version: 1,
				App: App{
					TimeRangeSeconds: 10,
					Position:         image.Point{10, 10},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.createFile {
				path = createTmpConfigFile(
					t, []byte(dedentYAMLString(tt.cfgFile)),
				)
			} else {
				path = getTmpConfigPath(t)
			}

			cfg := NewConfig(path)

			err := cfg.Load()
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			require.Equal(t, path, cfg.path)
			require.EqualExportedValues(t, tt.want, cfg)
		})
	}
}

// TODO: fix test
func TestConfig_Save(t *testing.T) {
	t.Run("Empty path", func(t *testing.T) {
		err := (&Config{path: ""}).Save()
		require.ErrorContains(t, err, "path to config file cannot be empty")
	})

	defaultFile := `
		version: 1
		app:
			enable_debug: false
			time_range_seconds: 120
			update_rate_seconds: 1
			bar_spacing: 0
			bar_width: 1
			plot_height: 30
			graph_settings: {}
			position:
				X: 10
				Y: 10
			theme:
				window:
					active:
						background: {"R": 20, "G": 20, "B": 0, "A": 20}
					passthrough:
						background: {"R": 10, "G": 0, "B": 0, "A": 20}
				plot:
					border: {"R": 150, "G": 100, "B": 100, "A": 205}
					midline: {"R": 0, "G": 200, "B": 0, "A": 100}
					bar: {"R": 250, "G": 0, "B": 0, "A": 105}
					label_text: {"R": 255, "G": 255, "B": 255, "A": 180}
					label_background: {"R": 0, "G": 0, "B": 0, "A": 0}
		`

	tests := []struct {
		name       string
		cfg        *Config
		createFile bool
		origFile   string
		wantFile   string
		wantErr    string
	}{
		// TODO: Add test for backup
		{
			name:       "File does not exist",
			createFile: false,
			wantFile:   defaultFile,
		},
		{
			name:       "Empty file",
			createFile: true,
			wantFile:   defaultFile,
		},
		{
			name:       "Invalid version",
			createFile: true,
			origFile:   `version: 0`,
			wantErr:    "file doesn't look like a valid config",
		},
		{
			name:       "Invalid YAML",
			createFile: true,
			origFile:   `invalid yaml`,
			wantErr:    "error loading config",
		},
		{
			name:       "Invalid config file struct",
			createFile: true,
			origFile:   `no_version: there`,
			wantErr:    "file doesn't look like a valid config",
		},
		{
			name: "Struct changed",
			cfg: &Config{
				Version: 42,
				App: App{
					TimeRangeSeconds: 60,
				},
			},
			wantFile: `
				version: 42
				app:
					time_range_seconds: 60
					position:
						X: 0
						Y: 0
				`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.createFile {
				path = createTmpConfigFile(
					t, []byte(dedentYAMLString(tt.origFile)),
				)
			} else {
				path = getTmpConfigPath(t)
			}

			cfg := NewConfig(path)
			if tt.cfg != nil {
				cfg.Version = tt.cfg.Version
				cfg.App = tt.cfg.App
			}

			err := cfg.Save()
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			wantMap := make(map[string]interface{})
			wantRaw := []byte(dedentYAMLString(tt.wantFile))
			err = yaml.Unmarshal(wantRaw, &wantMap)
			require.NoError(t, err)

			gotMap := make(map[string]interface{})
			gotRaw, err := os.ReadFile(path)
			require.NoError(t, err)
			err = yaml.Unmarshal(gotRaw, &gotMap)
			require.NoError(t, err)

			require.Equal(t, wantMap, gotMap)
		})
	}
}

func TestConfig_isVersionMatch(t *testing.T) {
	tests := []struct {
		version int
		want    bool
	}{
		{version: 0, want: false},
		{version: 1, want: true},
		{version: 2, want: false},
		{version: 10, want: false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := (&Config{Version: tt.version}).isVersionMatch()
			if got != tt.want {
				t.Errorf("Config.isVersionMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dedentYAMLString(str string) (result string) {
	return utils.TabToSpaces(utils.Dedent(str), 4)
}

func getTmpConfigPath(t *testing.T) (path string) {
	return filepath.Join(t.TempDir(), "koanf_mock")
}

func createTmpConfigFile(t *testing.T, content []byte) (path string) {
	path = getTmpConfigPath(t)
	err := os.WriteFile(path, content, 0o600)
	require.NoError(t, err, "error creating temp config file: %w", err)
	return path
}
