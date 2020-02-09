package env

import (
	"os"
)

// Set sets variable. Set global to true to apply to the OS environment,
// otherwise, the value will only be effective in this Env object and its
// children
func (env *Set) Set(global bool, key string, value string) error {
	if env.override == nil {
		env.override = make(map[string]string)
	}
	env.override[key] = value
	if global {
		return os.Setenv(key, value)
	}

	return nil
}

// UnSet will remove the variable from environment. If it is overriden,
// removes the value from own memory. Otherwise, it will unset from OS environemt
func (env *Set) UnSet(key string) error {
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

// Clone Create snapshot of every loaded environment variable
// (NOT the full OS environment)
// Cloned instance can be used for manipulations in limited scope (Set/UnSet)
func (env *Set) Clone() *Set {
	n := &Set{}

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
