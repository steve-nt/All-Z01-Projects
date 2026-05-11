# Road Intersection Simulator

A traffic simulation built with Rust and Macroquad. Vehicles spawn from four directions, navigate through an intersection, and obey traffic lights.

## Run

```bash
cargo run
```
## Controls
### Key	Action

↑	Spawn from south (northbound)

↓	Spawn from north (southbound)

→	Spawn from west (eastbound)

←	Spawn from east (westbound)

R	Random spawn

Esc Exit

## Features
Traffic lights: 10-second phases (N/S green ↔ E/W green)

Vehicle routes: Straight (green), Left turn (red), Right turn (blue)

Collision avoidance: Vehicles maintain safe following distance

Spawn protection: Vehicles only spawn when lane is clear

