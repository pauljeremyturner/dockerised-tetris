package server

import (
	"github.com/google/uuid"
)

type tetrisServiceState struct {
}

type TetrisService interface {
	NewGame() uuid.UUID
}

func NewTetrisService() TetrisService {
	return &tetrisServiceState{}
}

func (r *tetrisServiceState) NewGame() uuid.UUID {
	return uuid.New()
}
