package proof_of_work

import (
	"bytes"
	"errors"
	"testing"
)

type MockAlgo struct {
	puzzle string
	valid  bool
}

func (m *MockAlgo) Puzzle() string {
	return m.puzzle
}

func (m *MockAlgo) Verify(puzzle, solution string) bool {
	return m.valid
}

type MockConn struct {
	readData  []byte
	writeData bytes.Buffer
	readErr   error
	writeErr  error
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	copy(b, m.readData)
	return len(m.readData), nil
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return m.writeData.Write(b)
}

func TestProofOfWork(t *testing.T) {
	tests := []struct {
		name          string
		readData      []byte
		readErr       error
		writeErr      error
		algoValid     bool
		expectedError error
	}{
		{
			name:          "WriteError",
			readData:      []byte("solution"),
			writeErr:      errors.New("write error"),
			algoValid:     true,
			expectedError: errors.New("write error"),
		},
		{
			name:          "ReadError",
			readErr:       errors.New("read error"),
			algoValid:     true,
			expectedError: errors.New("read error"),
		},
		{
			name:          "VerifyPassed",
			readData:      []byte("solution"),
			algoValid:     true,
			expectedError: nil,
		},
		{
			name:          "VerifyFailed",
			readData:      []byte("failed_solution"),
			algoValid:     false,
			expectedError: errors.New("invalid puzzle"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &MockConn{
				readData: tt.readData,
				readErr:  tt.readErr,
				writeErr: tt.writeErr,
			}

			mockAlgo := &MockAlgo{
				puzzle: "puzzle",
				valid:  tt.algoValid,
			}

			pow := New(mockAlgo)
			err := pow.DDosProtection(mockConn)

			if err != nil && tt.expectedError == nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error %v, but got no error", tt.expectedError)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, but got %v", tt.expectedError, err)
			}
		})
	}
}
