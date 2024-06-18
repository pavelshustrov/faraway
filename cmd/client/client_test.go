package main

import (
	"faraway/internal/proof_of_work/hashcash"
	"testing"
)

func BenchmarkSolveChallenge(b *testing.B) {
	complexity := 5
	for i := 0; i < b.N; i++ {
		challenge := hashcash.New("word_of_wisdom", complexity).Puzzle()
		SolveChallenge(challenge, complexity)
	}
}
