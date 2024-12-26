// mocks/mock_api_client.go

package mocks

import (
    "sync"
)

// MockAPIClient is a manual mock implementation of the APIClient interface
type MockAPIClient struct {
    mutex     sync.RWMutex
    calls     map[string][]MockAPICall
    responses map[string]chan []byte
}

type MockAPICall struct {
    Endpoint string
    Method   string
    Payload  interface{}
}

func NewMockAPIClient() *MockAPIClient {
    return &MockAPIClient{
        calls:     make(map[string][]MockAPICall),
        responses: make(map[string]chan []byte),
    }
}

// MakeAPIRequest implements the APIClient interface
func (m *MockAPIClient) MakeAPIRequest(endpoint, method string, payload interface{}) chan []byte {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    // Record the call regardless of method
    key := endpoint + ":" + method
    m.calls[key] = append(m.calls[key], MockAPICall{
        Endpoint: endpoint,
        Method:   method,
        Payload:  payload,
    })

    if ch, exists := m.responses[key]; exists {
        return ch
    }

    // Default response if none is set
    ch := make(chan []byte, 1)
    ch <- []byte(`{"status":"success"}`)
    return ch
}

// SetResponse sets up a mock response for a specific endpoint and method
func (m *MockAPIClient) SetResponse(endpoint, method string, response []byte) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    key := endpoint + ":" + method
    ch := make(chan []byte, 1)
    ch <- response
    m.responses[key] = ch
}

// GetCalls returns the calls made to a specific endpoint and method
func (m *MockAPIClient) GetCalls(endpoint, method string) []MockAPICall {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    key := endpoint + ":" + method
    if calls, exists := m.calls[key]; exists {
        return calls
    }
    return []MockAPICall{}
}

// ClearCalls clears all recorded calls
func (m *MockAPIClient) ClearCalls() {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.calls = make(map[string][]MockAPICall)
}