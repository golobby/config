[![GoDoc](https://godoc.org/github.com/golobby/config?status.svg)](https://godoc.org/github.com/golobby/config)
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

### A simple example
The following example demonstrates how to set and get a simple key/value.

```go
c, err := config.New()
// Check error...

c.Set("name", "John Doe")

name, err := c.Get("name")
// Check error...
```

### Feeders
Feeders provide content of the configuration. Currently, these feeders exist out of the box:
* `Map`: Feeds a simple `map[string]interface{}`.
* `Json`: Feeds a JSON file.
* `JsonDirectory`: Feeds a directory of JSON files.

Of course, you are free to implement your feeders by implementing the `Feeder` interface.

You can pass your desired feeder through Options to the `New()` function this way:
```go
c, err := config.New(config.Options{
    Feeder: TheFeederGoesHere,
})
```

#### Feeding using Map feeder

You don't like config files!? It's OK you can pass a `Map` feeder to the Config initializer like this example:

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

#### Feeding using Json feeder

Storing configuration data in a JSON file could be a brilliant idea. The example below shows how to use Json feeder.

`config.json`:

```json
{
  "name": "MyAppUsingGoLobbyConfig",
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

Go code:

```go
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "path/to/config.json"},
})

v, err := c.Get("version") // 3.14

v, err := c.Get("numbers.2") // 3

v, err := c.Get("users.0.address.city") // Delfan
```

#### Feeding using JsonDirectory

If you have many configuration data and it doesn't fit in a single JSON file.
In this case, you can use multiple JSON files and feed them using JsonDirectory feeder like this example:

Sample project directory structure:

```
- main.go
- config
- - app.json
- - db.json
```

`app.json`:

```json
{
  "name": "MyApp",
  "version": 3.14
}
```

`db.json`:

```json
{
  "sqlite": { "path": "app.db" },
  "mysql": { "host": "localhost", "user": "root", "pass": "secret" }
}
```

Go code:

```go
c, err := config.New(config.Options{
    Feeder: feeder.JsonDirectory{Path: "config"},
})

v, err := c.Get("app.version") // 3.14
v, err := c.Get("db.mysql.host") // localhost
```

### OS variables and environment files

#### OS variables

Sometimes you need to use environment variables stored in OS alongside your configuration data.
You can refer to OS variables using a simple syntax, `${ VARIABLE }`, in the config values.
This example demonstrates how to use OS variables.

`db.json`:

```json
{"name": "MyApp", "port": "${ APP_PORT }"}
```

Go code:

```go
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "db.json"},
})

v, err := c.Get("name") // MyApp
v, err := c.Get("port") // equivalent to os.Getenv("APP_PORT")
```

If you need to have a default value in case of lacking OS variable you can use this syntax:

```
${ VARIABLE | DEFAULT }
```

Example of JSON file using OS variable:

```json
{"name": "MyApp", "port": "${ APP_PORT | 3306 }"}
```

#### environment files

You maybe want to use ".env" files. Good news! It's so easy to work with environment files.
You can pass an environment file path alongside your config feeder when you initialize a new instance of Config.

Sample project directory structure:

```
- main.go
- .env
- config.json
```

`config.json`:

```json
{"name": "MyApp", "key": "${ APP_KEY }", "port": "${ APP_PORT | 3306 }"}
```

`.env`:

```env
APP_KEY=secret
APP_PORT=
```

Go code:

```go
c, err := config.New(config.Options{
    Feeder: feeder.Json{Path: "config.json"},
    Env: ".env",
})

v, err := c.Get("name") // MyApp (from config.json)
v, err := c.Get("key") // secret (from .env)
v, err := c.Get("port") // 3306 (from config.json, the default value)
```

### Reload the config and env files

One of the benefits of using config management tools is the ability to change the configurations without redeployment.
The Config package takes advantage of OS signals to handle this need.
It listens to the "SIGHUP" signal and reloads the env and config files on receive.

You can send this signal to your application with following shell command:

```shell script
KILL -SIGHUP [YOUR-APP-PROCESS-ID]
```

To get your application process id you can use `ps` shell command.

### Altogether!

In this section, we illustrate a complete example that shows many of the package features.

Sample project directory structure:

```
- main.go
- .env
- config
- - app.json
- - db.json
```

`app.json`:

```json
{
  "name": "MyApp",
  "key": "${ APP_KEY }"
}
```

`db.json`:

```json
{
  "sqlite": {
    "path": "app.db"
  },
  "mysql": { 
    "host": "${ DB_HOST | localhost }",
    "user": "${ DB_USER | root }",
    "pass": "${ DB_PASS | secret }"
  }
}
```

`.env`:

```env
APP_KEY=theKey
DB_HOST=127.0.0.1
DB_USER=
DB_PASS=something
```

Go code:

```go
_ := os.Setenv("DB_HOST", "192.168.0.13")

c, err := config.New(config.Options{
    Feeder: feeder.JsonDirectory{Path: "config"},
    Env: ".env",
})

v, err := c.Get("app.name") // MyApp (from app.json)
v, err := c.Get("app.key")  // theKey (from .env)
v, err := c.Get("db.mysql.host")  // 192.168.0.13 (from OS variables)
v, err := c.Get("db.mysql.user")  // root (from app.json, the default value)
v, err := c.Get("db.mysql.pass")  // something (from .env)
```

You may ask what would happen if the value existed in the config file, the environment file, and the OS variables?
It's the order of the Config priorities:

1. OS Variables
1. Environment (.env) files
1. Default Value

So if the value was defined is OS variables, the Config would return it.
If it wasn't in OS variables, the Config would return the value stored in the environment file.
If it also wasn't in the environment file, it'd eventually return the value stored in the config file as default value.

## Contributors

* [@miladrahimi](https://github.com/miladrahimi)
* [@amirrezaask](https://github.com/amirrezaask)

## License

GoLobby Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
