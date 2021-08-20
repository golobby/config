[![GoDoc](https://godoc.org/github.com/golobby/config/v3?status.svg)](https://godoc.org/github.com/golobby/config/v3)
[![CI](https://github.com/golobby/config/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/config/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/config/workflows/CodeQL/badge.svg)](https://github.com/golobby/config/actions?query=workflow%3ACodeQL)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/config)](https://goreportcard.com/report/github.com/golobby/config)
[![Coverage Status](https://coveralls.io/repos/github/golobby/config/badge.svg)](https://coveralls.io/github/golobby/config?branch=master)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome) 

# Config
GoLobby Config is lightweight yet powerful configuration management for Go projects.
It takes advantage of dot env files and OS variables alongside config files to be your ultimate requirement.

## Documentation
### Required Go Version
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/golobby/config/v3
```

### Quick Start
The following example demonstrates how to use the package using a JSON configuration file.

```go
// My configuration struct
type MyConfig struct {
    App struct {
        Name string
        Port int
    }
    Debug      bool
    Production bool
    Pi         float64
}

// An instance of my configuration struct
myConfig := MyConfig{}

// Create a feeder that provides the configuration data from a JSON file
jsonFeeder := feeder.Json{Path: "config.json"}

// Create a Config instance, pass the feeder and feed the configuration struct
err := config.New(jsonFeeder).Feed(&myConfig)

// Use myConfig...
```

### Feeders
Feeders provide the configuration data.
The GoLobby Config package supports the following feeders out of the box.

* `Json`: It Feeds using a JSON file.
* `Yaml`: It Feeds using a YAML file.
* `DotEnv`: It Feeds using a dot env (.env) file.
* `Env`: It Feeds using OS environment variables.

You can also create your custom feeders by implementing the `Feeder` interface or use third-party feeders.

#### Json Feeder
Storing configuration data in a JSON file could be a brilliant idea.
The `Json` feeder uses Go built-in `json` package to load JSON files.

The example below shows how to use the `Json` feeder.

The JSON file: https://github.com/golobby/config/blob/v3/assets/sample1.json

```go
type MyConfig struct {
    App struct {
        Name string
        Port int
    }
    Debug      bool
    Production bool
    Pi         float64
}

myConfig := MyConfig{}

jsonFeeder := feeder.Json{Path: "sample1.json"}

err := config.New(jsonFeeder).Feed(&myConfig)

// Use myConfig...
```

#### Yaml Feeder
YAML files are also easy to use.
They could be another candidate for your configuration file.
The `Yaml` feeder uses the [YAML package](https://github.com/go-yaml/yaml) (v3) to load YAML files.

The YAML file: https://github.com/golobby/config/blob/v3/assets/sample1.yaml

```go
type MyConfig struct {
    App struct {
        Name string
        Port int
    }
    Debug      bool
    Production bool
    Pi         float64
}

myConfig := MyConfig{}

yamlFeeder := feeder.Yaml{Path: "sample1.yaml"}

err := config.New(yamlFeeder).Feed(&myConfig)

// Use myConfig...
```

#### DotEnv Feeder
Dot env (.env) files are popular configuration files.
They are usually declared per environment (production, local, test, etc.) differently.
The `DotEnv` feeder uses the [GoLobby DotEnv](https://github.com/golobby/dotenv) package to load `.env` files.

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
}

myConfig := MyConfig{}

dotEnvFeeder := feeder.DotEnv{Path: ".env"}

err := config.New(dotEnvFeeder).Feed(&myConfig)

// Use myConfig...
```

You must add a `env` tag for each field that determines the related dot env key.
If there isn't any value for a field in the related file, it ignores the struct field.

You can read more about this feeder in the [GoLobby DotEnv](https://github.com/golobby/dotenv) package repository.

#### Env Feeder
You may keep it simple stupid with no configuration files at all!

The `Env` feeder works fine in simple cases and cloud environments.
It feeds your structs by OS environment variables.
This feeder is built on top of the [GoLobby Env](https://github.com/golobby/env) package.

```go
_ = os.Setenv("APP_NAME", "Shop")
_ = os.Setenv("APP_PORT", "8585")
_ = os.Setenv("DEBUG", "true")
_ = os.Setenv("PRODUCTION", "false")
_ = os.Setenv("PI", "3.14")

type MyConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool    `env:"DEBUG"`
    Production bool    `env:"PRODUCTION"`
    Pi         float64 `env:"PI"`
}

myConfig := MyConfig{}
envFeeder := feeder.DotEnv{}

err := config.New(envFeeder).Feed(&myConfig)

// Use myConfig...
```

You must add a `env` tag for each field that determines the related OS environment variable name.
If there isn't any value for a field in OS environment variables, it ignores the struct field.

You can read more about this feeder in the [GoLobby Env](https://github.com/golobby/env) package repository.

### Multiple Feeders
One of the key features in the GoLobby Config package is feeding using multiple feeders.
Lately added feeders overrides early added ones.

* JSON file: https://github.com/golobby/config/blob/v3/assets/sample1.json
* DotEnv file: https://github.com/golobby/config/blob/v3/assets/.env.sample2
* Env (OS) variables: Defined in the Go code!

```go
_ = os.Setenv("PRODUCTION", "true")
_ = os.Setenv("APP_PORT", "6969")

type MyConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool    `env:"DEBUG"`
    Production bool    `env:"PRODUCTION"`
    Pi         float64 `env:"PI"`
}

myConfig := MyConfig{}

feeder1 := feeder.Json{Path: "sample1.json"}
feeder2 := feeder.DotEnv{Path: ".env.sample2"}
feeder3 := feeder.Env{}

err := config.New(feeder1, feeder2, feeder3).Feed(&myConfig)

fmt.Println(c.App.Name)   // Blog  [from DotEnv]
fmt.Println(c.App.Port)   // 6969  [from Env]
fmt.Println(c.Debug)      // false [from DotEnv]
fmt.Println(c.Production) // true  [from Env]
fmt.Println(c.Pi)         // 3.14  [from Json]
```

What happened?

* The `Json` feeder as the first feeder sets all the struct fields from the JSON file.
* The `DotEnv` feeder as the second feeder overrides existing fields.
  The `APP_NAME` and `DEBUG` fields exist in the `.env.sample2` file.
* The `Env` feeder as the last feeder overrides existing variables in the OS environment.
  The `APP_PORT` and `PRODUCTION` fields are defined.

### Refresh
The `Refresh()` method re-feeds the structs using the provided feeders.
It makes each feeder reload configuration data and feed the given structs again.

```go
c := config.New(feeder1, feeder2, feeder3)
err := c.Feed(&myConfig)

err = c.Refresh()

// myConfig fields are updated!
```

### Listener
One of the GoLobby Config features is the ability to update the configuration structs without redeployment.
It takes advantage of OS signals to handle this requirement.
Config instances listen to the "SIGHUP" operating system signal and refresh structs (call the `Refresh()` method).

To enable the listener for a Config instance, you should call the `WithListener()` method.
It gets a fallback function and calls it when the `Refresh()` method fails and returns an error.

```go
c := config.New(feeder).WithListener(func(err error) {
    fmt.Println(err)
})

err := c.Feed(&myConfig)
```

You can send the `SIGHUP` signal to your running application with the following shell command.

```shell script
KILL -SIGHUP [YOUR-APP-PROCESS-ID]
```

To get your application process ID, you can use the `ps` shell command.

## See Also
* [GoLobby/DotEnv](https://github.com/golobby/dotenv):
  A lightweight package for loading dot env (.env) files into structs for Go projects
* [GoLobby/Env](https://github.com/golobby/env):
  A lightweight package for loading OS environment variables into structs for Go projects

## License

GoLobby Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
