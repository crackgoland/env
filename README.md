# env

Read, manipulate process/context environmental variables with ease in Go.

- [x] Get string, int32, int64 and bool variables and check if set
- [x] One-liner Gets() with default value
- [x] Auto-trim, basic integer validation
- [x] Load .env file, merge with process environment
- [x] Set value contextually (internal) or globally (process)
- [x] Clone loaded/changed environment vars for further isolated reads/writes

# Usage

## Getters

```go
type Getter interface {
	String(string, string) (string, bool)
	Int(string, int) (int, error)
	Bool(string, bool) (bool, error)
	Int32(string, int32) (int32, error)
	Int64(string, int64) (int64, error)

	Default(string) DefaultGetter
	DefaultInt(int) DefaultGetterInt
	DefaultBool(bool) DefaultGetterBool
	DefaultInt32(int32) DefaultGetterInt32
	DefaultInt64(int64) DefaultGetterInt64
}
```

```go

package main

import "github.com/crackgoland/env"

var env env.Set

func init() {

  // No error returned if file does not exist
  if err := env.LoadFile(); err != nil {
    panic(err)
  }

  val := env.Default("default")("STRING_VAR")
  valInt := env.DefaultInt(0)("INTEGER_VAR")

  // and so on ..
}

```

The syntax in `Default*(..)(..)` is not trivial, but clean enough to use!

## Setters

```
type Setter interface {
	Set(global bool, key string, value string) error
	UnSet(key string) error
  Clone() *Set
}
```

Set/change globally OR only for current `env,Set` value.

```go
package main

import "github.com/crackgoland/env"

var myEnv env.Set

func init() {
  myEnv.Set(true, "VAR2", "This value is applied to the process environment")
  myEnv.Set(false, "VAR2", "This value will only be visible to users of 'myEnv'")

  val, isset := myEnv.String("VAR2", "") // => "This value will only be visible to users of 'myEnv'", true

  (&env.Env{}).String("VAR2", "") // => "This value is applied to the process environment", true
}
```


# See also
Alternatively I use [kelseyhightower/envconfig](github.com/kelseyhightower/envconfig)
to load envvars as struct fields,
see example [here](https://github.com/nytimes/gizmo/blob/master/pubsub/kafka/config.go).
