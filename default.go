package env

// Default get string value for key with fallback value
func (env *Set) Default(defaultValue string) DefaultGetter {
	return func(key string) string {
		val, _ := env.String(key, defaultValue)
		return val
	}
}

// DefaultInt get int value for key with fallback value
func (env *Set) DefaultInt(defaultValue int) DefaultGetterInt {
	return func(key string) int {
		val, _ := env.Int(key, defaultValue)
		return val
	}
}

// DefaultBool get bool value for key with fallback value
func (env *Set) DefaultBool(defaultValue bool) DefaultGetterBool {
	return func(key string) bool {
		val, _ := env.Bool(key, defaultValue)
		return val
	}
}

// DefaultInt32 get int32 value for key with fallback value
func (env *Set) DefaultInt32(defaultValue int32) DefaultGetterInt32 {
	return func(key string) int32 {
		val, _ := env.Int32(key, defaultValue)
		return val
	}
}

// DefaultInt64 get int64 value for key with fallback value
func (env *Set) DefaultInt64(defaultValue int64) DefaultGetterInt64 {
	return func(key string) int64 {
		val, _ := env.Int64(key, defaultValue)
		return val
	}
}
