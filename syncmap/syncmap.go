// VERY slim wrapper around sync.Map with Go generics for that sweet sweet type safety
package syncmap

import "sync"

type Map[K any, V any] struct {
	m            sync.Map
	defaultValue V
}

func NewMap[K any, V any]() (Map[K, V], error) {
	return Map[K, V]{}, nil
}

func NewMapWithDefaultValue[K any, V any](def V) (Map[K, V], error) {
	return Map[K, V]{
		defaultValue: def,
	}, nil
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, loaded bool) {
	v, loaded := m.m.Load(key)
	value, ok := v.(V)
	if !ok {
		return m.defaultValue, loaded
	}
	return
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)
	value, ok := v.(V)
	if !ok {
		return m.defaultValue, loaded
	}
	return
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.m.LoadOrStore(key, value)
	actual, ok := v.(V)
	if !ok {
		return m.defaultValue, loaded
	}
	return
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
