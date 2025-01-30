# Fogborne

## Project Overview

Fogborne is an ASCII-based survival game with real-time multiplayer capabilities. The project is structured to support incremental development, starting with core game mechanics and gradually expanding to include multiplayer features, procedural map generation, and advanced rendering effects.

## Initial Concept Overview

### ASCII Survival Board

- 2D grid (like a classic Rogue-like or MUD) with walls, empty spaces, items, and players.
- Each cell is represented by an ASCII character—e.g., '#' for walls, '.' for floor, '@' for player, etc.

### Real-time Multiplayer

- Multiple users connect to the same “world” or “zone.”
- Each user can see other users in adjacent squares, maybe interact with them or pick up items.

### Minimal Survival Mechanics (Optional)

- **Hunger/thirst**: a simple counter that decreases over time, forcing players to “hunt” for food or water cells.
- **Simple items**: a healing herb ('H') or something that replenishes health or hunger.
- **Combat**: to keep it simple initially, will skip this initially and focus on movement + environment.
