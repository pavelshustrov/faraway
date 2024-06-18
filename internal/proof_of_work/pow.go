package proof_of_work

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

var ErrVerificationFailed = errors.New("verification failed")

type Algo interface {
	Puzzle() string
	Verify(puzzle, solution string) bool
}

type ProofOfWork struct {
	algo       Algo
	readBuffer []byte
}

func New(algo Algo) *ProofOfWork {
	return &ProofOfWork{
		algo:       algo,
		readBuffer: make([]byte, 1024),
	}
}

func (p *ProofOfWork) DDosProtection(reader io.Reader, writer io.Writer) error {
	puzzle := p.algo.Puzzle()

	log.Println("generated puzzle", puzzle)

	if _, err := writer.Write([]byte(fmt.Sprintf("X-Hashcash: %s\n", puzzle))); err != nil {
		return err
	}

	solution, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		return err
	}

	solution = strings.TrimSpace(solution)
	fmt.Println("received solution", solution)
	if !p.algo.Verify(puzzle, solution) {
		return ErrVerificationFailed
	}

	fmt.Println("verified")
	return nil
}
