package hashcash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type HashCash struct {
	version    int
	complexity int
	resource   string
	timeFn     func() time.Time
}

var (
	versionOne = 1
)

func New(resource string, complexity int) HashCash {
	return HashCash{
		version:    versionOne,
		resource:   resource,
		complexity: complexity,
		timeFn:     time.Now,
	}
}

func (p HashCash) Puzzle() string {
	nonce := strings.ReplaceAll(uuid.New().String(), "-", "")

	challenge := fmt.Sprintf("%d:%s:%s::%d:%s", p.version, p.timeFn().Format("20060102"), p.resource, p.complexity, nonce)
	return challenge
}

func (p HashCash) Verify(puzzle, response string) bool {
	hash := sha1.New()
	hash.Write([]byte(puzzle))
	hash.Write([]byte(":"))
	hash.Write([]byte(response))

	return strings.HasPrefix(hex.EncodeToString(hash.Sum(nil)[:]), strings.Repeat("0", p.complexity))
}
