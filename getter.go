package env

import (
	"os"
	"strconv"
	"strings"
)

// TrimCharacters Control characters trimmed from left/right of every read value
var TrimCharacters = " \n\t"

// Get Returns the environment variable value from its name.
// If it wasn't set, second value is set to false. In both cases, may return
// spaces or empty
func (env Set) String(key string, defaultValue string) (val string, ok bool) {
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

	val = strings.Trim(val, TrimCharacters)

	return val, ok
}

// Int Parse int value from env value with specified key
func (env *Set) Int(key string, defaultValue int) (int, error) {
	v, err := env.Int64(key, int64(defaultValue))
	if err != nil {
		return defaultValue, err
	}
	i := int(v)
	if (v <= 0) == (i > 0) {
		return defaultValue, ErrOverflow // overflow
	}

	return i, nil
}

// Bool Parse bool value from env value with specified key
func (env *Set) Bool(key string, defaultValue bool) (bool, error) {
	s, ok := env.String(key, "")
	if !ok {
		return defaultValue, ErrNotSet
	} else if len(s) == 0 {
		return defaultValue, ErrEmpty
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue, err
	}
	return v, nil
}

// Int32 Parse int32 value from env value with specified key
func (env *Set) Int32(key string, defaultValue int32) (int32, error) {
	s, ok := env.String(key, "")
	if !ok {
		return defaultValue, ErrNotSet
	} else if len(s) == 0 {
		return defaultValue, ErrEmpty
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue, err
	}
	return int32(v), nil
}

// Int64 Parse int64 value from env value with specified key
func (env *Set) Int64(key string, defaultValue int64) (int64, error) {
	s, ok := env.String(key, "")
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
