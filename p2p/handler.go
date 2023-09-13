package p2p

import (
	"fmt"
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct {
}

func NewDefaultHandler() Handler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) HandleMessage(message *Message) error {
	b, err := io.ReadAll(message.Payload)
	if err != nil {
		return err
	}

	fmt.Printf("handle message from %s: %s\n", message.FromPeer, string(b))
	return nil
}
