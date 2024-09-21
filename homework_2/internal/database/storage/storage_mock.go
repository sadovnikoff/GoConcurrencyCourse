package storage

// MockEngine is mock of Engine interface
type MockEngine struct {
	Key   string
	Value string
}

// NewMockEngine creates a new mock instance
func NewMockEngine() *MockEngine {
	mock := &MockEngine{}
	return mock
}

// Del mocks method
func (m *MockEngine) Del(key string) {
	if m.Key == key {
		m.Key = ""
		m.Value = ""
	}
}

// Get mocks method
func (m *MockEngine) Get(key string) (string, error) {
	if m.Key == key {
		return m.Value, nil
	}
	return "", ErrNotFound
}

// Set mocks method
func (m *MockEngine) Set(key, value string) {
	m.Key = key
	m.Value = value
}
