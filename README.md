# proto-tetris
A terminal client communicates with a Tetris game engine via grpc.

## What is this?
An implementation of the classic Tetris game.  The game engine sends pixel information to the client.
The client is stupid.  It only paints pixels by co-ordinates and colour.

Client-server communication is via google protobuf streaming.

## Tooling

### Proto generation

## How do I make it work

First build the docker container:

`docker build -t tetris .`

Then run the container:

`docker run -p 50051:50051 tetris`

Then run the game client

`go run client-main.go`

## How do I play?

`s, d` Left & right move
`a, f` Left & right rotate
`x` Move down
`e` Drop

<img src="./doc/tetris-animated.gif" width="436" height="600">


