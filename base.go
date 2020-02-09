package env

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// ErrLoadFileAlready LoadFile() is designed to be used before any manipulations
	ErrLoadFileAlready = errors.New(".env file was loaded already, a variable was changed/added since then. Won't do load again")

	// ErrOverflow Primitive overflow check is triggered
	ErrOverflow = errors.New("Integer overflow")

	// ErrNotSet Envvar is unset
	ErrNotSet = errors.New("Envvar is unset")

	// ErrEmpty Envvar is set, but empty
	ErrEmpty = errors.New("Envvar is set, but empty")

	// FileName control .env file name (should not be path)
	FileName = ".env"
)

// Getter An interface of get functions for environment variables
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

// Setter An interface for modifier functions for environment variables
type Setter interface {
	Set(global bool, key string, value string) error
	UnSet(key string) error
	Clone() *Set
}

// DefaultGetter return string value or default, if key is not set
type DefaultGetter func(string) string

// DefaultGetterInt return int value or default, if key is not set or value invalid
type DefaultGetterInt func(string) int

// DefaultGetterBool  return boolean value or default, if key is not set or value invalid
type DefaultGetterBool func(string) bool

// DefaultGetterInt32 return int32 value or default, if key is not set or value invalid
type DefaultGetterInt32 func(string) int32

// DefaultGetterInt64 return int64 value or default, if key is not set or value invalid
type DefaultGetterInt64 func(string) int64

// Set Read, Copy or modify environemt variable set
type Set struct {
	fileData map[string]string
	override map[string]string
}

// Getter ensure compativility with interface Getter
func (env *Set) Getter() Getter {
	return env
}

// Setter ensure compativility with interface Setter
func (env *Set) Setter() Setter {
	return env
}

// LoadFile reads .env file from working directory. Loaded values
// will have priority on the OS environment.
// No error is returned if file does not exist
func (env *Set) LoadFile() error {
	if env.fileData != nil && len(env.fileData) != 0 {
		if env.override != nil {
			return ErrLoadFileAlready
		}
	}

	envFilename, err := filepath.Abs(FileName)
	if err != nil {
		message := fmt.Sprintf("filepath.Abs() error: %s", err.Error())
		return errors.New(message)
	}

	f, err := os.Open(envFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
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

		for i := range parts {
			parts[i] = strings.Trim(parts[i], " ,\t;#\"")
		}

		env.fileData[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ApplyLogLevel Dettermine log level from LOG_LEVEL variable and return it,
// also uses zerolog.SetGlobalLevel() to apply it
func (env Set) ApplyLogLevel(defaultLevel zerolog.Level) zerolog.Level {
	slevel, _ := env.String("LOG_LEVEL", "")

	level := defaultLevel

	switch strings.ToLower(slevel) {
	case "debug", "d":
		level = zerolog.DebugLevel
	case "info", "i":
		level = zerolog.InfoLevel
	case "warn", "warning", "w":
		level = zerolog.WarnLevel
	case "error", "err", "e":
		level = zerolog.ErrorLevel
	case "fatal", "f":
		level = zerolog.FatalLevel
	case "panic", "p":
		level = zerolog.PanicLevel
	case "off", "no", "none", "":
		level = zerolog.Disabled
	default:
		if slevel != "" {
			log.Warn().Msgf("Log level \"%s\" invalid.", slevel)
		}
	}

	zerolog.SetGlobalLevel(level)
	return level
}
