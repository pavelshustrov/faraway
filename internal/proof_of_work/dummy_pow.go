package proof_of_work

import "io"

type EmptyPow struct{}

func Off() *EmptyPow {
	return &EmptyPow{}
}

func (ep *EmptyPow) DDosProtection(_ io.Reader, _ io.Writer) error {
	return nil
}
