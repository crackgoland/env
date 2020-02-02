package env

import (
        "strings"
	"os"
        "strconv"
	"path/filepath"
	"bufio"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrCannotLoadFileData = errors.New(".env file was loaded previously, Set() was used since then and file cannot be safely reloaded since environment is overriden. Use LoadFile() once only, at start.")
  ErrOverflow = errors.New("Integer overflow")
  ErrNotSet = errors.New("Integer value is unset")
  ErrEmpty = errors.New("Integer value is empty")
)

type Env struct {
	fileData        map[string]string
  override        map[string]string
  DefaultValue    string
}

func (env Env) Get(key string, defaultValue string) (val string, ok bool) {

  if env.override != nil {
    if val, ok := env.override[key]; ok {
      return val, true
    }
  }

	if val, ok = os.LookupEnv(key); !ok {
		val, ok = env.fileData[key]
		if !ok {
			val = defaultValue
		}
	}

	val = strings.Trim(val, " \n\t;.=")

	return val, ok
}

func (env * Env) GetInt(key string, defaultValue int) (int, error) {
  v, err := env.GetInt64(key, int64(defaultValue))
  if err != nil {
    return defaultValue, err
  }
  i := int(v)
  if (v <= 0) == (i > 0) {
    return defaultValue, ErrOverflow // overflow
  }
  return i, nil
}

func (env * Env) GetInt64(key string, defaultValue int64) (int64, error) {
  s, ok := env.Get(key, "")
  if !ok {
    return defaultValue, ErrNotSet
  } else if len(s) == 0 {
    return defaultValue, ErrEmpty
  }
  v, err := strconv.ParseInt(s, 10, 64)
  if err != nil {
    return defaultValue, err
  }
  return v, nil
}

func (env * Env) LoadFile() error {
	if env.fileData != nil && len(env.fileData) != 0 {
    if env.override != nil {
  		return ErrCannotLoadFileData
    }
	}

	envFilename, err := filepath.Abs(".env")
	if err != nil {
		message := fmt.Sprintf("filepath.Abs() error: %s", err.Error())
		return errors.New(message)
	}
	log.Debug().Msg("Reading env file: " + envFilename)

	f, err := os.Open(envFilename)
	if err != nil {
		return err
	}

	defer f.Close()

	env.fileData = make(map[string]string, 0)

	scanner := bufio.NewScanner(f)
  for scanner.Scan() {
      line := scanner.Text()

			if strings.HasPrefix(strings.Trim(line, " \t"), "#") {
				continue
			}

			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}

			for i, _ := range parts {
				parts[i] = strings.Trim(parts[i], " ,\t;#\"")
			}

			env.fileData[parts[0]] = parts[1]
  }

  if err := scanner.Err(); err != nil {
      return err
  }

	return nil
}

func (env * Env) Set(global bool, key string, value string) error {
  if env.override == nil {
    env.override = make(map[string]string)
  }
  env.override[key] = value
  if global {
    return os.Setenv(key, value)
  }

  return nil
}

func (env * Env) UnSet(key string) error {
  global := true

  if env.override != nil {
    if _, ok := env.override[key]; ok {
      delete(env.override, key)
      global = false
    }
  }

  if env.fileData != nil {
    if _, ok := env.fileData[key]; ok {
      delete(env.fileData, key)
      global = false
    }
  }

  if !global {
    return nil
  }

  return os.Unsetenv(key)
}

func (env * Env) DefaultGet(key string) string {
  val, _ := env.Get(key, env.DefaultValue)
  return val
}

func (env * Env) Clone() * Env {
  n := &Env{nil, nil, env.DefaultValue}

  if env.fileData != nil {
    n.fileData = make(map[string]string)
    for key, val := range env.fileData {
      n.fileData[key] = val
    }
  }

  if env.override != nil {
    n.override = make(map[string]string)
    for key, val := range env.override {
      n.override[key] = val
    }
  }

  return n
}


func (env Env) SetLogLevel(defaultLevel zerolog.Level) {
  	LOG_LEVEL, _ := env.Get("LOG_LEVEL", "")
  	switch strings.ToLower(LOG_LEVEL) {
  	case "debug", "d":
  		zerolog.SetGlobalLevel(zerolog.DebugLevel)
  	case "info", "i":
  		zerolog.SetGlobalLevel(zerolog.InfoLevel)
  	case "warn", "warning", "w":
  		zerolog.SetGlobalLevel(zerolog.WarnLevel)
  	case "error", "err", "e":
  			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
  	case "fatal", "f":
  			zerolog.SetGlobalLevel(zerolog.FatalLevel)
  	case "panic", "p":
  			zerolog.SetGlobalLevel(zerolog.PanicLevel)
  	case "off", "no", "none", "":
  			zerolog.SetGlobalLevel(zerolog.Disabled)
    default:
        if LOG_LEVEL != "" {
          log.Warn().Msgf("Log level \"%s\" invalid.", LOG_LEVEL)
        }
  			zerolog.SetGlobalLevel(defaultLevel)
  	}
}
