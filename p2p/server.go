package p2p

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
)

type GameKing uint8

const (
	TexasHolden GameKing = iota
	Other
)

func (gk GameKing) String() string {
	switch gk {
	case TexasHolden:
		return "TEXAS HOLDEN"
	case Other:
		return "other"
	default:
		return "unknown"
	}
}

type ServerConfig struct {
	Version    string
	ListenAddr string
	GameType   GameKing
}

type Message struct {
	Payload  io.Reader
	FromPeer net.Addr
}

type Server struct {
	ServerConfig
	listener  net.Listener
	transport *TCPTransport
	peers     map[net.Addr]Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		gameState:    NewGameState(),
	}

	tp := NewTCPTransport(cfg.ListenAddr)
	s.transport = tp
	tp.AddPeerCh = s.addPeer
	tp.DelPeerCh = s.delPeer

	return s
}

func (s *Server) Start() error {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port": s.ServerConfig.ListenAddr,
		"type": s.GameType,
	}).Info("started new game server")

	return s.transport.ListenAndAccept()
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

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("player disconnected")

			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)
		case peer := <-s.addPeer:
			// if new player connect to the server we send our handshake and
			// wait his response
			if err := s.SendHandshake(peer); err != nil {
				logrus.Errorf("failed to send handshake, err: %s", err.Error())
			}

			if err := s.handshake(peer); err != nil {
				logrus.Errorf("failed to recive handshake, err: %s", err.Error())
			}

			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("handshake successful: new player connected")

			s.peers[peer.conn.RemoteAddr()] = *peer
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) SendHandshake(peer *Peer) error {
	hs := &Handshake{
		Version:  s.Version,
		GameKing: s.GameType,
	}

	buf := new(bytes.Buffer)
	if err := hs.Encode(buf); err != nil {
		return err
	}
	//if err := gob.NewEncoder(buf).Encode(hs); err != nil {
	//	return err
	//}

	return peer.send(buf.Bytes())
}

type Handshake struct {
	Version  string
	GameKing GameKing
}

func (hs *Handshake) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, []byte(hs.Version)); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, hs.GameKing)
}

func (hs *Handshake) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, []byte(hs.Version)); err != nil {
		return err
	}

	return binary.Read(r, binary.LittleEndian, hs.GameKing)
}

func (s *Server) handshake(peer *Peer) error {
	hs := &Handshake{}
	//if err := gob.NewDecoder(peer.conn).Decode(hs); err != nil {
	//	return err
	//}

	if err := hs.Decode(peer.conn); err != nil {
		return err
	}

	if s.GameType != hs.GameKing {
		return fmt.Errorf("mismatch game kind [%s] vs [%s]", s.GameType, hs.GameKing)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("mismatch version [%s] vs [%s]", s.Version, hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":      peer.conn.RemoteAddr(),
		"version":   hs.Version,
		"game kind": hs.GameKing,
	}).Info("received handshake")

	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
