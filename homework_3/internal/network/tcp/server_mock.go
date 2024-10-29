package tcp

import (
	"errors"
	"fmt"
	"strings"
)

// MockDatabase is mock of Database interface
type MockDatabase struct{}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

// HandleQuery mocks method
func (m *MockDatabase) HandleQuery(request string) (string, error) {

	if strings.Contains(request, "error") {
		return "", errors.New("error has been occurred during request handling")
	}

	return fmt.Sprintf("successful response to the [%s] request", request), nil
}
