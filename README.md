[![Build Status](https://travis-ci.org/golobby/config.svg?branch=master)](https://travis-ci.org/golobby/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/config)](https://goreportcard.com/report/github.com/golobby/config)
[![Coverage Status](https://coveralls.io/repos/github/golobby/config/badge.png?branch=master)](https://coveralls.io/github/golobby/config?branch=master)

# Config
GoLobby Config is a lightweight yet powerful config package for Go projects. 
It takes advantage of env files and OS variables alongside config files to be your ultimate requirement.

## Documentation

### Supported Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project

```bash
go get github.com/golobby/config
```

### Simple example of usage
The following example demonstrates how to set and get a simple key/value.

```go
c, err := config.New()
// Check error...

c.Set("name", "John Doe")

name, err := c.Get("name")
// Check error...
```

### Feeders
Feeders provide content of the configuration. Currently, these feeders are supported:
* `Map`: Feeds a simple `map[string]interface{}`.
* `Json`: Feeds a JSON file.
* `JsonDirectory`: Feeds a directory of JSON files.

You can pass your desired feeder through Options to the `New()` function this way:
```go
c, err := config.New(config.Options{
    Feeder: ...,
})
```

#### Feeding using Map

```go
c, err := config.New(config.Options{
    Feeder: feeder.Map{
        "name":     "Hey You",
        "band":     "Pink Floyd",
        "year":     1979,
        "duration": 4.6,
    },
})
if err != nil {
    panic(err)
}

name, err := c.Get("name")
// OR
name, err := c.GetString("name")

year, err := c.Get("year")
// OR
year, err := c.GetInt("year")

year, err := c.Get("duration")
// OR
duration, err := c.GetFloat("duration")

```

#### Feeding using Json

A sample of JSON file:

```json
{
  "name": "MyAppUsingConfig",
  "version": 3.14,
  "numbers": [
    1,
    2,
    3
  ],
  "users": [
    {
      "name": "Milad Rahimi",
      "year": 1993,
      "address": {
        "country": "Iran",
        "state": "Lorestan",
        "city": "Delfan"
      }
    },
    {
      "name": "Amirreza Askarpour",
      "year": 1998,
      "address": {
        "country": "Iran",
        "state": "Khouzestan",
        "city": "Ahvaz"
      }
    }
  ]
}
```

A sample of usage:

```go
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "path/to/config.json"},
})

v, err := c.Get("version") // 3.14

v, err := c.Get("numbers.2") // 3

v, err := c.Get("users.0.address.city") // Delfan
```

#### Feeding using JsonDirectory

Example of directory structure:

```
- main.go
- config
- - app.json
- - db.json
```

`app.json`:
```json
{"name": "MyApp", "version": 3.14}
```

`db.json`:
```json
{"sqlite": {"path": "app.db"}, "mysql":{"host": "localhost", "user": "root", "pass": "secret"}}
```

```go
c, err := config.New(config.Options{
    Feeder: feeder.JsonDirectory{Path: "config"},
})

v, err := c.Get("app.version") // 3.14

v, err := c.Get("db.mysql.host") // localhost
```


### Using OS variables and .env files
#### Using OS variables

```json
{"name": "MyApp", "port": "${ APP_PORT }"}
```

```go
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "path/to/config.json"},
})

v, err := c.Get("port") // equivalent to os.Getenv("APP_PORT")
```

#### Using .env files

```json
{"name": "MyApp", "port": "${ APP_PORT }"}
```

```env
APP_PORT=8585
```

```json
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "path/to/config.json"},
    EnvFile: "path/to/.env",
})

v, err := c.Get("port") // 8585

```

#### Default Values

```json
{"name": "MyApp", "port": "${ APP_PORT | 80 }"}
```

```env
APP_PORT=
```

```json
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "path/to/config.json"},
    EnvFile: "path/to/.env",
})

v, err := c.Get("port") // 80
```

#### Priority of value sources

1. OS Variables
1. Environment (.env) files
1. Default Value

## Contributors

* [@miladrahimi](https://github.com/miladrahimi)

## License

GoLobby Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
