package utils

import (
	"sync"
	"time"
)

var mu sync.RWMutex

type ConcurrentMap struct {
	c map[string]interface{}
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	mu.Lock()
	if m.c == nil {
		m.c = make(map[string]interface{})
	}

	m.c[key] = value
	mu.Unlock()
}

func (m ConcurrentMap) Clear() {
	m.c = nil
}

func (m ConcurrentMap) Get(key string) (value interface{}, exists bool) {
	mu.RLock()
	value, exists = m.c[key]
	mu.RUnlock()
	return value, exists
}
func (m *ConcurrentMap) MustGet(key string) interface{} {
	if value, exists := m.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (m *ConcurrentMap) GetString(key string) (s string) {
	if val, ok := m.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (m *ConcurrentMap) GetBool(key string) (b bool) {
	if val, ok := m.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (m *ConcurrentMap) GetInt(key string) (i int) {
	if val, ok := m.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (m *ConcurrentMap) GetInt64(key string) (i64 int64) {
	if val, ok := m.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (m *ConcurrentMap) GetUint(key string) (ui uint) {
	if val, ok := m.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (m *ConcurrentMap) GetUint64(key string) (ui64 uint64) {
	if val, ok := m.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (m *ConcurrentMap) GetFloat64(key string) (f64 float64) {
	if val, ok := m.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (m *ConcurrentMap) GetTime(key string) (t time.Time) {
	if val, ok := m.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (m *ConcurrentMap) GetDuration(key string) (d time.Duration) {
	if val, ok := m.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (m *ConcurrentMap) GetStringSlice(key string) (ss []string) {
	if val, ok := m.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *ConcurrentMap) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := m.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *ConcurrentMap) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := m.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (m *ConcurrentMap) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := m.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
