package server

import (
	"faraway/internal/proof_of_work"
	"faraway/internal/proof_of_work/hashcash"
	"time"
)

type Config struct {
	DDosProtectionStrategy string
	Port                   string
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
}

func getDDosProtection(v string) DDosProtector {
	switch v {
	case "OFF":
		return &proof_of_work.EmptyPow{}
	case "HASHCASH":
		return proof_of_work.New(hashcash.New("word_of_wisdom", 5))
	default:
		panic("unknown ddos protection")
	}
}
