package main

import "github.com/pauljeremyturner/dockerised-tetris/client"

func main() {

	ui = client.NewTetrisClientUi(listenKeyboard)

	tp = client.NewTetrisProto(updateBard)
	ui.NewGame()

}

var tp client.TetrisProto
var ui client.TetrisClientUi

func listenKeyboard(r rune) {
	tp.Move(r)
}

func updateBard(gs client.GameState) {
	ui.Update(gs)
}
