package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var (
	ErrIO              = errors.New("IO error")
	ErrUnknownResource = errors.New("unknown resource")
)

type DDosProtector interface {
	DDosProtection(reader io.Reader, writer io.Writer) error
}
type QuoteDic interface {
	GetRandomQuote() string
}

type ConnectionHandler struct {
	readBuf []byte
	quotes  QuoteDic
	pow     DDosProtector
}

func (handler *ConnectionHandler) Handle(conn net.Conn) error {
	log.Println("handling new connection")

	reader := bufio.NewReader(conn)
	var msg string
	msg, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("%w: %w", ErrIO, err)
	}

	msg = strings.TrimSpace(msg)
	log.Println("received message:", msg)

	if msg != "word_of_wisdom" {
		return ErrUnknownResource
	}

	if err := handler.pow.DDosProtection(reader, conn); err != nil {
		return err
	}

	quote := handler.quotes.GetRandomQuote()
	log.Println("sending message:", quote)

	if _, err := conn.Write(append([]byte(quote), '\n')); err != nil {
		return fmt.Errorf("%w: %w", ErrIO, err)
	}

	return nil
}
