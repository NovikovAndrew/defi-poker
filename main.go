package main

import (
	"defi-poker/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "DEFI poker v0.1-alpha\n",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(cfg)

	go server.Start()

	remoteCfg := p2p.ServerConfig{
		Version:    "DEFI poker v0.1-alpha\n",
		ListenAddr: ":4000",
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	remoteServer.Connect(":3000")

	select {}
}
