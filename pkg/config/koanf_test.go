package config

import (
	"testing"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: koanf.UnmarshalWithConf won't work with private fields
type ConfigTest struct {
	nonKoanfField string
	// version       int   `koanf:"version"`
	SomeString    string `koanf:"some_string"`
	privateString string `koanf:"private_string"`
	FromFile      string `koanf:"from_file"`
	FromStruct    string `koanf:"from_struct"`
	// App     App  `koanf:"app"`
}

func TestKoanf_Defaults(t *testing.T) {
	del := "."

	tests := []struct {
		name               string
		defaults           ConfigTest
		defaultsResult     ConfigTest
		defaultsToFile     string
		loadFromFile       string
		loadFromFileResult ConfigTest
		loadFromFileToFile string
	}{
		{
			defaults: ConfigTest{
				nonKoanfField: "nonKoanfValue",
				SomeString:    "value from struct",
				privateString: "private value",
				FromStruct:    "value from struct",
			},
			defaultsResult: ConfigTest{
				nonKoanfField: "",
				SomeString:    "value from struct",
				privateString: "",
				FromFile:      "",
				FromStruct:    "value from struct",
			},
			defaultsToFile: `
				from_file: ""
				from_struct: value from struct
				some_string: value from struct
				`,
			loadFromFile: `
				nonKoanfField: "value from file"
				some_string: "value from file"
				privateString: "value from file"
				from_file: "value from file"
				top:
					nested: "value from file"
				`,
			loadFromFileResult: ConfigTest{
				nonKoanfField: "",
				SomeString:    "value from file",
				privateString: "",
				FromFile:      "value from file",
				FromStruct:    "value from struct",
			},
			loadFromFileToFile: `
				from_file: value from file
				from_struct: value from struct
				nonKoanfField: value from file
				privateString: value from file
				some_string: value from file
				top:
					nested: value from file
				`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var cfgResult ConfigTest
			k := koanf.New(del)
			// Load defaults
			err = k.Load(structs.Provider(tt.defaults, "koanf"), nil)
			require.NoError(t, err)

			err = k.UnmarshalWithConf("", &cfgResult, koanf.UnmarshalConf{})
			require.NoError(t, err)

			assert.Equal(t, tt.defaultsResult, cfgResult)

			// Check saved file
			cfgBytes, _ := k.Marshal(parser)
			assert.Equal(t, string(cfgBytes), dedentYAMLString(tt.defaultsToFile))

			tmpFile := createTmpConfigFile(
				t, []byte(dedentYAMLString(tt.loadFromFile)),
			)

			// Load config file over defaults
			err = k.Load(file.Provider(tmpFile), parser)
			require.NoError(t, err, "error loading config: %w", err)

			cfgResult = ConfigTest{}
			err = k.UnmarshalWithConf("", &cfgResult, koanf.UnmarshalConf{})
			require.NoError(t, err)
			// koanf.UnmarshalConf{
			// 	DecoderConfig: &mapstructure.DecoderConfig{
			// 		ErrorUnused: bool,
			// 		ErrorUnset: bool,
			// 		IgnoreUntaggedFields: bool,
			// 		MatchName: func(mapKey, fieldName string) bool,
			// 		DecodeNil: bool,
			// 	},
			// }

			assert.Equal(t, tt.loadFromFileResult, cfgResult)

			// Check saved file
			cfgBytes, _ = k.Marshal(parser)
			assert.Equal(t, dedentYAMLString(tt.loadFromFileToFile), string(cfgBytes))
		})
	}
}
