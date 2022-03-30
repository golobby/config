[![GoDoc](https://godoc.org/github.com/golobby/config/v3?status.svg)](https://godoc.org/github.com/golobby/config/v3)
[![CI](https://github.com/golobby/config/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/config/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/config/workflows/CodeQL/badge.svg)](https://github.com/golobby/config/actions?query=workflow%3ACodeQL)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/config)](https://goreportcard.com/report/github.com/golobby/config)
[![Coverage Status](https://coveralls.io/repos/github/golobby/config/badge.svg)](https://coveralls.io/github/golobby/config?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

# Config
GoLobby Config is a lightweight yet powerful configuration manager for Go projects.
It takes advantage of Dot-env (.env) files and OS environment variables alongside config files (JSON, YAML, and TOML) to meet all of your requirements.

## Documentation
### Required Go Version
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/golobby/config/v3
```

### Quick Start
The following example demonstrates how to use a JSON configuration file.

```go
// The configuration struct
type MyConfig struct {
    App struct {
        Name string
        Port int
    }
    Debug      bool
    Production bool
    Pi         float64
}

// Create an instance of the configuration struct
myConfig := MyConfig{}

// Create a feeder that provides the configuration data from a JSON file
jsonFeeder := feeder.Json{Path: "config.json"}

// Create a Config instance and feed `myConfig` using `jsonFeeder`
c := config.New()
c.AddFeeder(jsonFeeder)
c.AddStruct(&myConfig)
err := c.Feed()

// Or use method chaining:
// err := config.New().AddFeeder(jsonFeeder).AddStruct(&myConfig).Feed()

// Use `myConfig`...
```

### Feeders
Feeders provide the configuration data.
The GoLobby Config package supports the following feeders out of the box.

* `Json`: It feeds using a JSON file.
* `Yaml`: It feeds using a YAML file.
* `Toml`: It feeds using a TOML file.
* `DotEnv`: It feeds using a dot env (.env) file.
* `Env`: It feeds using OS environment variables.

You can also create your custom feeders by implementing the `Feeder` interface or use third-party feeders.

#### Json Feeder
The `Json` feeder uses Go built-in `json` package to load JSON files.
The snippet below shows how to use the `Json` feeder.

```go
jsonFeeder := feeder.Json{Path: "sample1.json"}
c := config.New().AddFeeder(jsonFeeder)
```

#### Yaml Feeder
The `Yaml` feeder uses the [YAML package](https://github.com/go-yaml/yaml) (v3) to load YAML files.
The snippet below shows how to use the `Yaml` feeder.

```go
yamlFeeder := feeder.Yaml{Path: "sample1.yaml"}
c := config.New().AddFeeder(yamlFeeder)
```

#### Toml Feeder
The `Toml` feeder uses the [BurntSushi TOML package](https://github.com/BurntSushi/toml) to load TOML files.
The snippet below shows how to use the `Toml` feeder.

```go
tomlFeeder := feeder.Toml{Path: "sample1.toml"}
c := config.New().AddFeeder(tomlFeeder)
```

#### DotEnv Feeder
The `DotEnv` feeder uses the [GoLobby DotEnv](https://github.com/golobby/dotenv) package to load `.env` files.
The example below shows how to use the `DotEnv` feeder.

The `.env` file: https://github.com/golobby/config/blob/v3/assets/.env.sample1

```go
type MyConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool    `env:"DEBUG"`
    Production bool    `env:"PRODUCTION"`
    Pi         float64 `env:"PI"`
    IDs        []int   `env:"IDS"`
}

myConfig := MyConfig{}
dotEnvFeeder := feeder.DotEnv{Path: ".env"}
err := config.New().AddFeeder(dotEnvFeeder).AddStruct(&myConfig).Feed()
```

You must add a `env` tag for each field that determines the related dot env variable.
If there isn't any value for a field in the related file, it ignores the struct field.
You can read more about this feeder in the [GoLobby DotEnv](https://github.com/golobby/dotenv) package repository.

#### Env Feeder
The `Env` feeder is built on top of the [GoLobby Env](https://github.com/golobby/env) package.
The example below shows how to use the `Env` feeder.

```go
_ = os.Setenv("APP_NAME", "Shop")
_ = os.Setenv("APP_PORT", "8585")
_ = os.Setenv("DEBUG", "true")
_ = os.Setenv("PRODUCTION", "false")
_ = os.Setenv("PI", "3.14")
_ = os.Setenv("IPS", "192.168.0.1", "192.168.0.2")
_ = os.Setenv("IDS", "10, 11, 12, 13")

type MyConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool     `env:"DEBUG"`
    Production bool     `env:"PRODUCTION"`
    Pi         float64  `env:"PI"`
    IPs        []string `env:"IPS"`
    IDs        []int16  `env:"IDS"`
}

myConfig := MyConfig{}
envFeeder := feeder.DotEnv{}
err := config.New().AddFeeder(envFeeder).AddStruct(&myConfig).Feed()
```

You must add a `env` tag for each field that determines the related OS environment variable name.
If there isn't any value for a field in OS environment variables, it ignores the struct field.
You can read more about this feeder in the [GoLobby Env](https://github.com/golobby/env) package repository.

### Multiple Feeders
One of the key features in the GoLobby Config package is feeding using multiple feeders.
Lately added feeders overrides early added ones.

The example below demonstrates how to use a JSON file as the main configuration feeder and override the configurations with dot env and os variables.

* JSON file: https://github.com/golobby/config/blob/v3/assets/sample1.json
* DotEnv file: https://github.com/golobby/config/blob/v3/assets/.env.sample2
* Env (OS) variables: Defined in the Go code!

```go
_ = os.Setenv("PRODUCTION", "true")
_ = os.Setenv("APP_PORT", "6969")
_ = os.Setenv("IDs", "6, 9")

type MyConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool    `env:"DEBUG"`
    Production bool    `env:"PRODUCTION"`
    Pi         float64 `env:"PI"`
    IDs        []int32 `env:"IDS"`
}

myConfig := MyConfig{}

feeder1 := feeder.Json{Path: "sample1.json"}
feeder2 := feeder.DotEnv{Path: ".env.sample2"}
feeder3 := feeder.Env{}

err := config.New()
        .AddFeeder(feeder1)
        .AddFeeder(feeder2)
        .AddFeeder(feeder3)
        .AddStruct(&myConfig)
        .Feed()

fmt.Println(c.App.Name)   // Blog  [from DotEnv]
fmt.Println(c.App.Port)   // 6969  [from Env]
fmt.Println(c.Debug)      // false [from DotEnv]
fmt.Println(c.Production) // true  [from Env]
fmt.Println(c.Pi)         // 3.14  [from Json]
fmt.Println(c.IDs)        // 6, 9  [from Env]
```

What happened?

* The `Json` feeder as the first feeder sets all the struct fields from the JSON file.
* The `DotEnv` feeder as the second feeder overrides existing fields.
  The `APP_NAME` and `DEBUG` fields exist in the `.env.sample2` file.
* The `Env` feeder as the last feeder overrides existing fields, as well.
  The `APP_PORT` and `PRODUCTION` fields are defined in the OS environment.

### Re-feed
You can re-feed the structs every time you need to.
Just call the `Feed()` method again.

```go
c := config.New().AddFeeder(feeder).AddStruct(&myConfig)
err := c.Feed()

// Is it time to re-feed?
err = c.Feed()

// Use `myConfig` with updated data!
```

### Listener
One of the GoLobby Config features is the ability to update the configuration structs without redeployment.
It takes advantage of OS signals to handle this requirement.
Config instances listen to the "SIGHUP" operating system signal and refresh structs (call the `Feed()` method).

To enable the listener for a Config instance, you should call the `SetupListener()` method.
It gets a fallback function and calls it when the `Feed()` method fails and returns an error.

```go
c := config.New().AddFeeder(feeder).AddStruct(&myConfig)
c.SetupListener(func(err error) {
    fmt.Println(err)
})

err := c.Feed()
```

You can send the `SIGHUP` signal to your running application with the following shell command.

```shell script
KILL -SIGHUP [YOUR-APP-PROCESS-ID]
```

You can get your application process ID using the `ps` command.

## See Also
* [GoLobby/DotEnv](https://github.com/golobby/dotenv):
  A lightweight package for loading dot env (.env) files into structs for Go projects
* [GoLobby/Env](https://github.com/golobby/env):
  A lightweight package for loading OS environment variables into structs for Go projects
* [GoLobby/Container](https://github.com/golobby/container):
  A lightweight yet powerful IoC dependency injection container for Go projects

## License
GoLobby Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
