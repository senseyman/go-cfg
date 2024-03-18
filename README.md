# Go Config

The `go-cfg` package is used for mapping ENV variables to your config struct with default options.


# Install
```shell
    go get github.com/senseyman/go-cfg
```

# Usage

## Organize your config structure
### Example of config structure
```go
type (
	Config struct {
		LogLevel  string `mapstructure:"LOG_LEVEL"    def:"DEBUG"`
		DebugMode bool   `mapstructure:"DEBUG_MODE"   def:"true"`
		AppName   string `mapstructure:"APP_NAME"`
		DB        MySQL  `mapstructure:"DB"`
	}

	MySQL struct {
		Host     string `mapstructure:"HOST"       def:"127.0.0.1"`
		Port     int    `mapstructure:"PORT"       def:"1234"`
		Username string `mapstructure:"USERNAME"   def:"admin"`
	}
)
```
Use `mapstructure` tag to set ENV name. You can create nested structures (like DB in the example).

Use `def` tag for setting default value.

### Example of using
```go
package main

import (
	cfgReader "github.com/senseyman/go-cfg"
)

func main() {
	var cfg Config
	err := cfgReader.Read(&cfg)
	// other code
}
```

## Dependencies

`go-cfg` uses [viper](https://github.com/spf13/viper)