# dockerised-tetris
a containerised tetris game engine communicates with client via protobuf

Work in progress

DT-1
basic project setup.  A modular golang project.  server: a docker container that can server new tetris board requests, client: a protobuf client to send new tetris board request
DONE

DT-2
build a UI that can be updated with board updates and register key presses for user moves.  Add dockerfile for server process
DONE

DT-3
enable protobuf streaming for stubbed game state updates to be displayed on UI

DT-4
implement matrix algebra for piece rotations, provide a way to define coordinates of pieces
