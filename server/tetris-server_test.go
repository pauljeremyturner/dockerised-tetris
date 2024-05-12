package server

import (
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"sync"
	"testing"
	"time"
)

func TestShouldMoveDownPeriodically(t *testing.T) {
	u, _ := uuid.NewRandom()
	mq := make(chan shared.MoveType, 5)
	ss := &serverSession{
		player:      Player{u, "player"},
		moveQueue:   mq,
		gameQueue:   nil,
		activePiece: NewI(),
		lines:       Lines{},
		nextPiece:   NewI(),
		gameOver:    false,
	}

	go ss.Tick()

	time.Sleep(3 * time.Second)

	ss.gameOver = true

	time.Sleep(1 * time.Second)
	assert.True(t, len(mq) >= 2)

	close(mq)

	for got := range mq {
		assert.True(t, got == shared.DOWN)
	}

}

func Test_tetris_GetGame(t *testing.T) {
	type fields struct {
		activeGames *sync.Map
		sugar       *zap.SugaredLogger
	}
	type args struct {
		u  uuid.UUID
		ss ServerSession
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ServerSession
		wantErr bool
	}{
		{
			name: "success: retrieve ongoing game",
			fields: fields{
				activeGames: &sync.Map{},
				sugar:       zaptest.NewLogger(t).Sugar(),
			},
			args: args{
				u:  uuid.MustParse("3858a3a8-4fc6-4536-bae9-f1797deb57c7"),
				ss: NewServerSession(Player{}, zaptest.NewLogger(t).Sugar()),
			},
			want:    NewServerSession(Player{}, zaptest.NewLogger(t).Sugar()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tetris{
				activeGames: tt.fields.activeGames,
				sugar:       tt.fields.sugar,
			}
			got, err := r.GetGame(tt.args.u)
			assert.True(t, tt.wantErr == (err != nil))

			assert.Equal(t, tt.want.GetPlayer(), got.GetPlayer())
		})
	}
}
