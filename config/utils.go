package config

func MapGetFunc[K comparable, V any](m map[K]any, key K, f func(v any, present bool) V) V {
	v, present := m[key]
	return f(v, present)
}

func MapGetDefault[K comparable, V any](m map[K]any, key K, def V) V {
	return MapGetFunc(m, key, func(v any, present bool) V {
		if !present {
			return def
		}
		value, ok := v.(V)
		if !ok {
			return def
		}
		return value
	})
}
