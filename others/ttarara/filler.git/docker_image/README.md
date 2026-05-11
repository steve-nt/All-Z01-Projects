# Filler docker image

- To build the image `docker build -t filler .`
- To run the container `docker run -v "$(pwd)/solution":/filler/solution -it filler`. This instruction will open a terminal in the container, the directory `solution` will be mounted in the container as well.
- Example of a command in the container `./linux_game_engine -f maps/map01 -p1 linux_robots/bender -p2 linux_robots/terminator`
- Your solution should be inside the `solution` directory so it will be mounted and compiled inside the container and it will be able to be run in the game engine.

## Required bundle files

The official 01-edu zip must provide these untracked files before the Docker image can be built:

- `linux_game_engine`
- `m1_game_engine`
- `linux_robots/`
- `m1_robots/`

This repo keeps the Rust solution, maps, Dockerfile, and scripts. The engine and robot binaries are intentionally gitignored because each student brings the official bundle locally.

## Notes

- `Terminator` is a very strong robot so it's optional to beat him.
- For M1 Macs use `m1_robots` and `m1_game_engine`.