package p2p

type Round uint32

const (
	Dealing Round = iota
	PreFlop
	Flop
	Turn
	River
)

type GameState struct {
	Round    uint32 //  atomic access
	isDealer bool   // atomic access
}

func NewGameState() *GameState {
	return &GameState{}
}

func (gs *GameState) loop() {
	for {
		select {}
	}
}
