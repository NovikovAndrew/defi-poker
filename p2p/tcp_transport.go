package p2p

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) ReadLoop(msgCh chan *Message) {
	buf := make([]byte, 1<<10)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}
		msgCh <- &Message{
			Payload:  bytes.NewReader(buf[:n]),
			FromPeer: p.conn.RemoteAddr(),
		}
	}

	// TODO unregister this peer
	if err := p.conn.Close(); err != nil {
		fmt.Printf("tcp connection err: %s\n", err.Error())
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeerCh  chan *Peer
	DelPeerCh  chan *Peer
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer{
			conn: conn,
		}

		t.AddPeerCh <- peer

	}

	return fmt.Errorf("TCP transport stopped")
}
