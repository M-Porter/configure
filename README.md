# configure

Dead simple configuration library.

# Install

```
go get -u github.com/m-porter/configure/v3
```

# Usage

`configure` uses a struct-based approach. Define your config shape as a Go struct
with `mapstructure` tags, then use a `Configure` instance to load values into it
from a file, environment variables, or defaults.

## Basic setup (from file)

```go
type AppConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
}

cfg := &AppConfig{}

conf := configure.New()
conf.SetConfigName("config")   // file name without extension
conf.SetConfigType("yaml")     // e.g. "yaml", "toml", "json"
conf.SetConfigDir("/etc/myapp")

if err := conf.Get(cfg); err != nil {
    log.Fatal(err)
}
```

Field names in your struct must have `mapstructure` tags that match the keys in
your config file or environment variables.

## Default values

Use `SetDefaults` to supply fallback values. Defaults are applied before the config
file or environment variables are read, so they can always be overridden.

```go
conf := configure.New()
conf.SetDefaults(AppConfig{
    Host: "localhost",
    Port: 5432,
})

if err := conf.Get(cfg); err != nil {
    log.Fatal(err)
}
```

## Environment variables

Environment variables are read automatically. The variable name is the uppercase
version of the `mapstructure` key (e.g. the `host` field maps to `HOST`).

Use `SetEnvPrefix` to namespace your variables and avoid collisions with other
programs:

```go
conf := configure.New()
conf.SetEnvPrefix("myapp") // HOST becomes MYAPP_HOST, PORT becomes MYAPP_PORT

if err := conf.Get(cfg); err != nil {
    log.Fatal(err)
}
```

## Write config if not exists

Set `SetWriteIfNotExists(true)` to have `configure` write a config file on the
first run if one does not already exist. Defaults will be used as the initial
values.

```go
conf := configure.New()
conf.SetConfigName("config")
conf.SetConfigType("yaml")
conf.SetConfigDir("/etc/myapp")
conf.SetDefaults(AppConfig{
    Host: "localhost",
    Port: 5432,
})
conf.SetWriteIfNotExists(true)

if err := conf.Get(cfg); err != nil {
    log.Fatal(err)
}
```

## Saving config

Use `Save` to persist a modified config struct back to disk.

```go
cfg.Host = "db.prod.example.com"
cfg.Port = 5433

if err := conf.Save(cfg); err != nil {
    log.Fatal(err)
}
```

Note: `Configure` is frozen after the first call to `Get` or `Save`. Calling any
setter after that point will return a `ConfigurationFrozenError`.

# Further examples

See [`configure_test.go`](./configure_test.go) for more detailed usage examples.
