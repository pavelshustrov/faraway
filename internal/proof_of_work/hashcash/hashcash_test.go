package hashcash

import (
	"bytes"
	"crypto/sha1"
	"github.com/google/uuid"
	"regexp"
	"testing"
	"time"
)

func TestPuzzle(t *testing.T) {
	fixedTime := time.Date(2024, time.June, 19, 0, 0, 0, 0, time.UTC)
	timeFn := func() time.Time { return fixedTime }

	hashCash := HashCash{
		version:    1,
		timeFn:     timeFn,
		resource:   "example@example.com",
		complexity: 20,
	}

	puzzle := hashCash.Puzzle()

	regex := `^X-Hashcash: \d+:\d{8}:[^:]+::\d+:[A-Za-z0-9_-]{32}\n$`
	matched, err := regexp.MatchString(regex, puzzle)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !matched {
		t.Errorf("Puzzle did not match expected format. Got: %s", puzzle)
	}
}

func generateValidResponse(puzzle string, complexity int) string {
	var response string
	for {
		response = uuid.New().String()
		hash := sha1.New()
		hash.Write([]byte(puzzle))
		hash.Write([]byte(":"))
		hash.Write([]byte(response))
		if bytes.HasSuffix(hash.Sum(nil), bytes.Repeat([]byte{'0'}, complexity)) {
			break
		}
	}
	return response
}

func TestVerify(t *testing.T) {
	fixedTime := time.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)
	timeFn := func() time.Time { return fixedTime }

	hashCash := HashCash{
		version:    1,
		timeFn:     timeFn,
		resource:   "example@example.com",
		complexity: 1, // Using a low complexity for testing
	}

	puzzle := hashCash.Puzzle()

	validResponse := generateValidResponse(puzzle, hashCash.complexity)

	if !hashCash.Verify(puzzle, validResponse) {
		t.Errorf("Valid response was not verified. Puzzle: %s, Response: %s", puzzle, validResponse)
	}

	invalidResponse := "invalid_response"
	if hashCash.Verify(puzzle, invalidResponse) {
		t.Errorf("Invalid response was incorrectly verified. Puzzle: %s, Response: %s", puzzle, invalidResponse)
	}
}
