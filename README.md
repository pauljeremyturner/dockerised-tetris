# dockerised-tetris
a containerised tetris game engine communicates with client via protobuf streaming

## Work in progress.  

Right now communication between client and server is working, game pieces move as expected.  missing logic is ensuring pieces can only move according to bounds of the 'board' a new piece is created when the active piece stops and determining end of game.

DT-1
~~basic project setup.  A modular golang project.  server: a docker container that can server new tetris board requests, client: a protobuf client to send new tetris board request~~


DT-2
~~build a UI that can be updated with board updates and register key presses for user moves.  Add dockerfile for server process~~


DT-3
~~enable protobuf streaming for stubbed game state updates to be displayed on UI~~


DT-4
implement matrix algebra for piece rotations, provide a way to define coordinates of pieces


DT-5
use google wire for dependency injection
