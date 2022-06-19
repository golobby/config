package config_test

import (
    "errors"
    "github.com/golobby/config/v3"
    "github.com/golobby/config/v3/pkg/feeder"
    "github.com/stretchr/testify/assert"
    "os"
    "syscall"
    "testing"
    "time"
)

type Sex int

const (
    Male Sex = iota
    Female
)

type FullConfig struct {
    App struct {
        Name string `env:"APP_NAME"`
        Port int    `env:"APP_PORT"`
    }
    Debug      bool     `env:"DEBUG"`
    Production bool     `env:"PRODUCTION"`
    Pi         float64  `env:"PI"`
    IPs        []string `env:"IPS"`
    IDs        []int16  `env:"IDS"`
    SexRaw     int      `env:"SEX"`
    Sex        Sex
}

func (fc *FullConfig) Setup() error {
    if fc.SexRaw == 0 {
        fc.Sex = Male
    } else if fc.SexRaw == 1 {
        fc.Sex = Female
    } else {
        return errors.New("app: invalid sex")
    }

    return nil
}

func TestConfig_Feed_With_No_Data(t *testing.T) {
    c := &struct{}{}
    err := config.New().AddFeeder(feeder.Env{}).AddStruct(c).Feed()
    assert.NoError(t, err)
}

func TestConfig_Feed_With_Invalid_File_It_Should_Fail(t *testing.T) {
    s := struct{}{}
    c := config.New().AddFeeder(feeder.Json{}).AddStruct(&s)
    err := c.Feed()
    assert.EqualError(t, err, "config: feeder error: json: read .: is a directory")
}

func TestConfig_Feed(t *testing.T) {
    _ = os.Setenv("PRODUCTION", "1")
    _ = os.Setenv("APP_PORT", "6969")

    c := &FullConfig{}

    f1 := feeder.Json{Path: "assets/sample1.json"}
    f2 := feeder.DotEnv{Path: "assets/.env.sample2"}
    f3 := feeder.Env{}

    err := config.New().AddFeeder(f1, f2, f3).AddStruct(c).Feed()
    assert.NoError(t, err)

    assert.Equal(t, "Blog", c.App.Name)
    assert.Equal(t, 6969, c.App.Port)
    assert.Equal(t, false, c.Debug)
    assert.Equal(t, true, c.Production)
    assert.Equal(t, 3.14, c.Pi)
    assert.Equal(t, []string{"192.168.0.1", "192.168.0.2"}, c.IPs)
    assert.Equal(t, []int16{10, 11, 12, 13}, c.IDs)

    assert.Equal(t, Male, c.Sex)
}

func TestConfig_Feed_With_Setup_Returning_Error(t *testing.T) {
    _ = os.Setenv("SEX", "3")

    c := &FullConfig{}

    f1 := feeder.Json{Path: "assets/sample1.json"}
    f2 := feeder.Env{}

    err := config.New().AddFeeder(f1, f2).AddStruct(c).Feed()
    assert.Error(t, err, "app: invalid sex")
}

func TestConfig_ReFeeding(t *testing.T) {
    _ = os.Setenv("NAME", "One")

    s := &struct {
        Name string `env:"NAME"`
    }{}

    c := config.New().AddFeeder(feeder.Env{}).AddStruct(s)
    err := c.Feed()
    assert.NoError(t, err)

    assert.Equal(t, "One", s.Name)

    _ = os.Setenv("NAME", "Two")

    err = c.Feed()
    assert.NoError(t, err)

    assert.Equal(t, "Two", s.Name)
}

func TestConfig_SetupListener(t *testing.T) {
    _ = os.Setenv("PI", "3.14")

    s := &struct {
        Pi float64 `env:"PI"`
    }{}

    fallbackTested := false
    c := config.New().AddFeeder(feeder.Env{}).AddStruct(s).SetupListener(func(err error) {
        fallbackTested = true
    })

    err := c.Feed()
    assert.NoError(t, err)

    assert.Equal(t, 3.14, s.Pi)

    _ = os.Setenv("PI", "3.666")

    err = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
    assert.NoError(t, err)

    time.Sleep(10 * time.Millisecond)

    assert.Equal(t, 3.666, s.Pi)

    _ = os.Setenv("PI", "INVALID!")

    err = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
    assert.NoError(t, err)

    time.Sleep(10 * time.Millisecond)

    assert.Equal(t, true, fallbackTested)
    assert.Equal(t, 3.666, s.Pi)
}
