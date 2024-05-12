package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/pauljeremyturner/dockerised-tetris/protogen"
	"github.com/pauljeremyturner/dockerised-tetris/shared"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type protoServer struct {
	tetris *tetris
	protogen.UnimplementedMoveServer
	protogen.UnimplementedGameServer
}

func NewProtoServer(tetris *tetris) *protoServer {
	return &protoServer{
		tetris: tetris,
	}
}

type ProtoServer interface {
	CreateNewGame() uuid.UUID
	Move(uuid uuid.UUID)
}

type StreamCompletedMovePublisher struct {
	in     *protogen.NewGameRequest
	stream protogen.Game_StartGameServer
}

func (r *StreamCompletedMovePublisher) PublishCompletedMove(ctx context.Context, gs GameState) error {
	gameUpdateResponse := &protogen.GameUpdateResponse{
		Uuid:       r.in.Uuid,
		PlayerName: r.in.PlayerName,
		Lines:      uint32(gs.LineCount),
		Pieces:     uint32(gs.PieceCount),
		Duration:   gs.Duration,
		Squares:    pixelsToSquares(gs.Pixels),
		NextPiece:  pixelsToSquares(gs.NextPiece),
	}

	if err := r.stream.Send(gameUpdateResponse); err != nil {
		return err
	}

	return nil
}

func (s *protoServer) StartGame(in *protogen.NewGameRequest, stream protogen.Game_StartGameServer) error {

	u, _ := uuid.Parse(in.Uuid)
	ss := s.tetris.StartNewGame(Player{UUID: u, Name: in.PlayerName})
	s.tetris.activeGames.Store(u, ss)

	pub := &StreamCompletedMovePublisher{
		in:     in,
		stream: stream,
	}
	err := ss.PublishCompletedMoves(stream.Context(), pub)
	if err != nil {
		return err
	}

	s.tetris.sugar.Infof("Game update channel closed, game over")

	return nil
}

func (s *protoServer) Move(ctx context.Context, in *protogen.MoveRequest) (*protogen.MoveResponse, error) {
	u, err := uuid.Parse(in.Uuid)
	if err != nil {
		return nil, NewErrorInvalidRequest("invalid game UUID", err)
	}
	m, err := protoMoveTypeToModel(in.GetMove())
	if err != nil {
		return nil, NewErrorInvalidRequest("invalid game move", err)
	}
	s.tetris.EnqueueMove(u, m)

	return &protogen.MoveResponse{}, nil
}

func protoMoveTypeToModel(me protogen.MoveRequest_MoveEnum) (shared.MoveType, error) {

	switch me {
	case protogen.MoveRequest_ROTATELEFT:
		return shared.ROTATELEFT, nil
	case protogen.MoveRequest_ROTATERIGHT:
		return shared.ROTATERIGHT, nil
	case protogen.MoveRequest_MOVELEFT:
		return shared.MOVELEFT, nil
	case protogen.MoveRequest_MOVERIGHT:
		return shared.MOVERIGHT, nil
	case protogen.MoveRequest_DROP:
		return shared.DROP, nil
	case protogen.MoveRequest_DOWN:
		return shared.DOWN, nil
	default:
		return shared.UNKNOWN, errors.New("")
	}
}

func pixelsToSquares(pixels []Pixel) []*protogen.Square {
	squares := make([]*protogen.Square, 0)

	for _, p := range pixels {
		squares = append(squares, &protogen.Square{
			X:     uint32(p.X),
			Y:     uint32(p.Y),
			Color: colorToProto(p.Color),
		})
	}

	return squares
}

type UnaryErrorInterceptor struct {
	fn func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error)
	sugar *zap.SugaredLogger
}

func NewUnaryErrorHandler(sugar *zap.SugaredLogger) *UnaryErrorInterceptor {
	return &UnaryErrorInterceptor{
		sugar: sugar,
	}
}

func (ei *UnaryErrorInterceptor) errorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)

	if err != nil {

		if errors.Is(err, ErrorInvalidRequest{}) {
			ei.sugar.Infof("invalid request, cause: %v", err)
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
	}

	return resp, nil
}
