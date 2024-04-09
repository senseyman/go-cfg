package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	testConfigStruct struct {
		LogLevel     string       `mapstructure:"LOG_LEVEL"    def:"DEBUG"`
		DebugMode    bool         `mapstructure:"DEBUG_MODE"   def:"true"`
		ID           string       `mapstructure:"ID"`
		DBTestStruct testDBStruct `mapstructure:"DB"`
	}

	testDBStruct struct {
		Host string `mapstructure:"HOST"   def:"127.0.0.1"`
		Port int    `mapstructure:"PORT"   def:"1234"`
	}
)

func TestRead(t *testing.T) {
	testCases := []struct {
		name        string
		hasErr      bool
		input       testConfigStruct
		expected    testConfigStruct
		setDefaults func()
	}{
		{
			name:   "success/defaults/all",
			hasErr: false,
			expected: testConfigStruct{
				LogLevel:  "DEBUG",
				DebugMode: true,
				ID:        "",
				DBTestStruct: testDBStruct{
					Host: "127.0.0.1",
					Port: 1234,
				},
			},
			setDefaults: func() {},
		},
		{
			name:   "success/defaults/part",
			hasErr: false,
			expected: testConfigStruct{
				LogLevel:  "INFO",
				DebugMode: true,
				ID:        "",
				DBTestStruct: testDBStruct{
					Host: "127.0.0.1",
					Port: 4321,
				},
			},
			setDefaults: func() {
				os.Setenv("LOG_LEVEL", "INFO")
				os.Setenv("DB_PORT", "4321")
			},
		},
		{
			name:   "success/all_set",
			hasErr: false,
			expected: testConfigStruct{
				LogLevel:  "INFO",
				DebugMode: false,
				ID:        "000999",
				DBTestStruct: testDBStruct{
					Host: "192.168.0.1",
					Port: 5555,
				},
			},
			setDefaults: func() {
				os.Setenv("LOG_LEVEL", "INFO")
				os.Setenv("DEBUG_MODE", "false")
				os.Setenv("ID", "000999")
				os.Setenv("DB_HOST", "192.168.0.1")
				os.Setenv("DB_PORT", "5555")
			},
		},
		{
			name:   "err/wrong_type",
			hasErr: true,
			expected: testConfigStruct{
				LogLevel:  "DEBUG",
				DebugMode: false,
				ID:        "",
				DBTestStruct: testDBStruct{
					Host: "127.0.0.1",
					Port: 1234,
				},
			},
			setDefaults: func() {
				os.Setenv("DEBUG_MODE", "abc")
			},
		},
	}

	for idx := range testCases {
		tc := testCases[idx]

		t.Run(tc.name, func(t *testing.T) {
			defer t.Cleanup(resetEnvs)

			tc.setDefaults()

			err := Read(&tc.input)

			if tc.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.EqualValues(t, tc.expected, tc.input)
		})
	}
}

func resetEnvs() {
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("DEBUG_MODE")
	os.Unsetenv("ID")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}
