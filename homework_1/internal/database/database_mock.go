package database

import (
	"errors"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_1/internal/database/compute"
)

// MockComputeLayer is mock of ComputeLayer interface
type MockComputeLayer struct {
}

// NewMockComputeLayer creates a new mock instance
func NewMockComputeLayer() *MockComputeLayer {
	mock := &MockComputeLayer{}
	return mock
}

// MockStorageLayer is mock of StorageLayer interface
type MockStorageLayer struct {
	Key   string
	Value string
}

// NewMockStorageLayer creates a new mock instance
func NewMockStorageLayer() *MockStorageLayer {
	mock := &MockStorageLayer{}
	return mock
}

// Parse mocks method
func (m *MockComputeLayer) Parse(cmd string) (compute.Query, error) {

	switch cmd {
	case compute.SetCommand:
		return compute.NewQuery(cmd, "key", "value"), nil
	case compute.GetCommand:
		return compute.NewQuery(cmd, "key", ""), nil
	case compute.DelCommand:
		return compute.NewQuery(cmd, "key", ""), nil
	}

	return compute.Query{}, errors.New("some error")
}

// Del mocks method
func (m *MockStorageLayer) Del(key string) {

}

// Get mocks method
func (m *MockStorageLayer) Get(key string) (string, error) {
	return "value", nil
}

// Set mocks method
func (m *MockStorageLayer) Set(key, value string) {

}
