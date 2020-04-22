package client

type Block struct {
	X     int
	Y     int
	Color int
}

type GameState struct {
	Blocks    []Block
	NextPiece []Block
	GameOver  bool
	Lines     int
	Duration  int64
}

type MoveType rune
