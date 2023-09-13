package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Message struct {
	Payload  io.Reader
	FromPeer net.Addr
}

type Server struct {
	ServerConfig
	handler  Handler
	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
	addPeer  chan *Peer
	delPeer  chan *Peer
	msgCh    chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		handler:      NewDefaultHandler(),
		peers:        make(map[net.Addr]Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}
}

func (s *Server) Start() error {
	go s.loop()

	if err := s.listen(); err != nil {
		return err
	}

	fmt.Printf("the server start on port :%s\n", s.ServerConfig.ListenAddr)
	s.acceptLoop()
	return nil
}

func (s *Server) listen() error {
	ls, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	s.listener = ls
	return nil
}

// TODO Cretae new network for new room
// maybe construct a new peer and handshake after registtration a plain connection
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.send([]byte(s.Version))
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer
		if err := peer.send([]byte(s.Version)); err != nil {
			panic(err)
		}

		go s.handleConn(peer)
	}
}

func (s *Server) handleConn(p *Peer) {
	buf := make([]byte, 1<<10)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		message := &Message{
			Payload:  bytes.NewReader(buf[:n]),
			FromPeer: p.conn.RemoteAddr(),
		}

		s.msgCh <- message

		fmt.Println(string(buf[:n]))
	}

	s.delPeer <- p
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			addr := peer.conn.RemoteAddr()
			s.mu.Lock()
			delete(s.peers, addr)
			s.mu.Unlock()
			fmt.Printf("peer deleted, tcp address %s\n", addr)
		case peer := <-s.addPeer:
			fmt.Printf("new peer connected, tcp address %s\n", peer.conn.RemoteAddr())
			s.mu.Lock()
			s.peers[peer.conn.RemoteAddr()] = *peer
			s.mu.Unlock()
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
