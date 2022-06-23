package textclient

import (
	"bufio"
	"log"
	"net"
)

// TextClient sends payload via tcp
type TextClient struct {
	address string
}

func New(address string) (*TextClient, error) {
	return &TextClient{address: address}, nil
}

func (t *TextClient) Send() {
	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		log.Println("could not created client connection to server: %w", err)

		return
	}

	text := "ball,ball,ball,eggs,pool,wild,daily"
	writer := bufio.NewWriter(conn)

	if _, err := writer.WriteString(text); err != nil {
		log.Println("failed writing to connection: %w", err)
	}

	if err := writer.Flush(); err != nil {
		log.Println("failed flushing: %w", err)
	}
}
