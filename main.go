package main

import (
	"defi-poker/p2p"
	"time"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "DEFI poker v0.1-alpha\n",
		ListenAddr: ":3000",
		GameType:   p2p.TexasHolden,
	}
	server := p2p.NewServer(cfg)

	go server.Start()
	time.Sleep(time.Millisecond * 300)

	remoteCfg := p2p.ServerConfig{
		Version:    "DEFI poker v0.1-alpha\n",
		ListenAddr: ":4000",
		GameType:   p2p.TexasHolden,
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	remoteServer.Connect(":3000")

	select {}
}
