[![GoDoc](https://godoc.org/github.com/golobby/config?status.svg)](https://godoc.org/github.com/golobby/config)
[![Build Status](https://travis-ci.org/golobby/config.svg?branch=master)](https://travis-ci.org/golobby/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/config)](https://goreportcard.com/report/github.com/golobby/config)
[![Coverage Status](https://coveralls.io/repos/github/golobby/config/badge.svg?branch=master)](https://coveralls.io/github/golobby/config?branch=master)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome)

# Config

GoLobby Config is a lightweight yet powerful config package for Go projects. It takes advantage of env files and OS
variables alongside config files to be your ultimate requirement.

## Documentation

### Required Go Version

It requires Go `v1.11` or newer versions.

### Installation

To install this package run the following command in the root of your project

```bash
go get github.com/golobby/config
```

### Getting Started

The following example demonstrates how to set and get a simple configuration key/value.

```go
c, err := config.New()
// Check error...

c.Set("name", "Pink Floyd")

name, err := c.Get("name")
```

### Feeders

Feeders provide the configuration contents.
The Config package supports the following feeders out of the box.

* `Map` Feeds using a simple `map[string]interface{}`.
* `Json` Feeds using a JSON file.
* `JsonDirectory` Feeds using a directory of JSON files.
* `Yaml` Feeds using a YAML file.
* `YamlDirectory` Feeds using a directory of YAML files.
* `Env` Feeds using an environment file.
* `OS` Feeds using a list of OS variables.

Of course, you can implement your custom feeders by implementing the `Feeder` interface.

Pass feeders to the `New()` method like the following example.

```go
f := feeder.Map{"name": "Pink Floyd"}

c, err := config.New(f)
name, err := c.Get("name") // "Pink Floyd"
```

Or you can pass multiple feeders like this example:

```go
f1 := feeder.Map{"name": "Divison Bell"}
f2 := feeder.Map{"year": 1994}

c, err := config.New(f1, f2)
name, err := c.Get("name") // "Pink Floyd"
name, err := c.Get("year") // 1994
```

#### Feeding using Map

You don't like config files!?
It's OK you can use the Map feeder like this example:

```go
c, err := config.New(feeder.Map{
    "name": "Hey You",
    "band": "Pink Floyd",
    "year": 1979,
    "rate": 4.6,
})

name, err := c.Get("name")
// OR
name, err := c.GetString("name")

year, err := c.Get("year")
// OR
year, err := c.GetInt("year")

year, err := c.Get("rate")
// OR
duration, err := c.GetFloat("rate")
```

#### Feeding using Json

Storing configuration data in a JSON file could be a brilliant idea. The example below shows how to use the Json feeder.

JSON file: https://github.com/golobby/config/blob/v2/feeder/test/config.json

Go code:

```go
c, err := config.New(feeder.Json{Path: "path/to/config.json"})

v, err := c.GetFloat("version")         // 3.14
v, err := c.GetInt("numbers.2")         // 3
v, err := c.Get("users.0.address.city") // Delfan
```

#### Feeding using JsonDirectory

You might have many configuration data and it doesn't fit in a single JSON file. In this case, you can use multiple JSON
files and feed them using JsonDirectory feeder like this example:

Sample project config directory:

```
- main.go
- config
- - app.json
- - db.json
```

JSON directory: https://github.com/golobby/config/tree/v2/feeder/test/json

Go code:

```go
c, err := config.New(feeder.JsonDirectory{Path: "path/to/config"})

v, err := c.GetFloat("app.version")     // 3.14
v, err := c.GetString("db.mysql.host")  // localhost
```

#### Feeding using Yaml

YAML files are a trend these days, so why not store configurations in them?

YAML file: https://github.com/golobby/config/blob/v2/feeder/test/config.yaml

Go code:

```go
c, err := config.New(feeder.Yaml{Path: "path/to/config.yaml"})

v, err := c.GetFloat("version")         // 3.14
v, err := c.GetInt("numbers.2")         // 3
v, err := c.Get("users.0.address.city") // Delfan
```

#### Feeding using YamlDirectory

You might have many configuration data and it doesn't fit in a single YAML file. In this case, you can use multiple YAML
files and feed them using YamlDirectory feeder like this example:

Sample project config directory:

```
- main.go
- config
- - app.yaml
- - db.yaml
```

Yaml directory: https://github.com/golobby/config/tree/v2/feeder/test/yaml

Go code:

```go
c, err := config.New(feeder.YamlDirectory{Path: "path/to/config"})

v, err := c.GetFloat("app.version")     // 3.14
v, err := c.GetString("db.mysql.host")  // localhost
```

#### Feeding using Env

An environment file could be a main or a secondary feeder to override other feeder values.

Because of different key names in env files, their names would be updated this way:

* `APP_NAME` => `app.name`
* `DB_MYSQL_HOST` => `db.mysql.host`

ENV file: https://github.com/golobby/config/blob/v2/feeder/test/.env

Go code:

```go
c, err := config.New(feeder.Env{Path: "path/to/.env"})

v, err := c.GetString("url")     // https://example.com (Original Key: URL)
v, err := c.GetString("db.host") // 127.0.0.1 (Original Key: DB_HOST)
```

Env feeder fetches operating system variables when the value is empty.

Go code:

```go
// Set an OS variable
_ = os.Setenv("APP_NAME", "MyAppUsingConfig")

c, err := config.New(feeder.Env{Path: "path/to/.env"})

v, err := c.Get("app.name") // MyAppUsingConfig (empty in .env ==> OS variable)
```

You can disable this feature this way:

```go
// Set an OS variable
_ = os.Setenv("APP_NAME", "MyAppUsingConfig")

c, err := config.New(feeder.Env{Path: "test/.env", DisableOSVariables: true})

v, err := c.Get("app.name") // "" (empty as in .env)
```

#### Feeding using OS

There is another feeder named OS that fetches OS variables and updates variable names like Env feeder. You should use OS
feeders to override other feeders like Env, Json and Yaml feeders.

```go
// Set an OS variable
_ = os.Setenv("APP_NAME", "MyAppUsingConfig")
_ = os.Setenv("APP_VERSION", "3.14")

c, err := config.New(OS{Keys: []string{"APP_NAME", "APP_VERSION", "APP_EMPTY"}})

v, err := c.Get("app.name")     // "MyAppUsingConfig"
v, err := c.Get("app.version")  // 3.14
v, err := c.Get("app.empty")    // ""
v, err := c.Get("app.new")      // ERROR!
```

### Multiple Feeders

One of the key features in the Config package is feeding using multiple feeders. Lately added feeders will override
early added ones.

```go
c, err := config.New(
    feeder.Map{
        "name": "MyAppUsingConfig"
        "url": "going to be overridden by the next feeders",
    },
    feeder.Map{
        "version": 3.14
        "url": "going to be overridden by the next feeder",
    },
    feeder.Map{
        "url": "https://github.com/golobby/config",
    },
)

v, err := c.Get("name")     // "MyAppUsingConfig"
v, err := c.Get("version")  // 3.14
v, err := c.Get("url")      // "https://github.com/golobby/config"
```

### A Sample using Json and Env feeders and OS variables

`path/to/config/app.json`:

```json
{
  "name": "MyAppUsingConfig",
  "version": 3.14,
  "url": "http://localhost"
}
```

`path/to/config/db.json`:

```json
{
  "mysql": {
    "host": "127.0.0.1",
    "user": "root",
    "pass": "secret"
  }
}
```

`.env`:

```
APP_URL=https://github.com/golobby/config
DB_MYSQL_HOST=192.168.0.1
DB_MYSQL_PASS=
```

Go code:

```go
_ = os.Setenv("APP_URL", "http://192.168.0.1")
_ = os.Setenv("DB_MYSQL_PASS", "password")

c, err := config.New(
    feeder.JsonDirectory{Path: "path/to/config"},
    feeder.Env{Path: ".env"}
)

v, err := c.Get("app.name")
// "MyAppUsingConfig" (from path/to/config/app.json)

v, err := c.Get("app.version")
// 3.14 (from path/to/config/app.json)

v, err := c.Get("app.url")
// "https://github.com/golobby/config" (from .env)
// APP_URL exists and it's not empty in .env so the Env feeder doesn't load it from OS variables.
// The value from .env overrides the value in the app.json file

v, err := c.Get("db.mysql.host")
// "192.168.0.1" (from .env)
// DB_MYSQL_HOST in .env overrides the value in the db.json file

v, err := c.Get("db.mysql.user")
// "root" (from path/to/config/app.json)

v, err := c.Get("db.mysql.pass")
// "password" (from OS variables)
// The key exists in .env but it's empty so the Env feeder loads it
// from OS variables (This feature is on in default) and override the value in the db.json file 
```

### Reload the feeders

One of the benefits of using config management tools is the ability to update the configurations without redeployment.
The Config package takes advantage of OS signals to handle this requirement.
It listens to the "SIGHUP" signal and reloads the config files (feeders) on receive.

You can send this signal to your application with the following shell command:

```shell script
KILL -SIGHUP [YOUR-APP-PROCESS-ID]
```

To get your application process id you can use the `ps` shell command.

### Strict Type Values

When you use the `Get()` method you have to cast the returned value type before using it.
You can use strict type methods like the following examples, instead.

List of strict type methods:
* `GetString()` returns only string values
* `GetInt()` returns only int and float values (pruning the floating part)
* `GetFloat()` returns only float and int values (in float64 type)
* `GetBool()` returns true, "true" and 1 as true, and false, "false" and 0 as false
* `GetStrictBool()` returns only boolean values (true and false)

## License

GoLobby Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
