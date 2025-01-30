# Fogborne Architecture Documentation

## Project Overview

Fogborne is an ASCII-based survival game with real-time multiplayer capabilities. The project is structured to support incremental development, starting with core game mechanics and gradually expanding to include multiplayer features, procedural map generation, and advanced rendering effects.

## Directory Structure and Components

### Command Layer (`/cmd`)

The command directory contains the entry points for our different executables. We separate these to maintain clear boundaries between different server types:

- `server/`: Main game server handling WebSocket connections and HTTP API endpoints
- `sshserver/`: Dedicated SSH server for terminal-based game access

### Internal Components (`/internal`)

The internal directory contains the core game logic and supporting systems. Each package is designed with a specific responsibility:

#### Game Core (`/internal/game`)

- `world/`: Manages game world and map systems

  - `ziglib/`: Contains our Zig-based procedural map generation
  - `map.go`: Core map data structures and operations
  - `generation.go`: Interface between Go and Zig for map generation
  - `zone.go`: Handles different game areas and instancing

- `entity/`: Defines game entities and their behaviors

  - `player.go`: Player-specific logic
  - `npc.go`: Non-player character behaviors
  - `state.go`: Shared entity state management

- `systems/`: Core game mechanics
  - `combat.go`: Combat system logic
  - `inventory.go`: Item and inventory management
  - `fov.go`: Field of view and fog of war calculations

#### Networking (`/internal/net`)

Handles all network communication:

- `websocket/`: Real-time WebSocket communication
- `ssh/`: Terminal-based game access
- `api/`: HTTP endpoints for non-real-time operations

#### Rendering (`/internal/render`)

Manages ASCII rendering and visual effects:

- `ascii/`: Core ASCII rendering logic
- `camera/`: Viewport management and player perspective
- `engine.go`: Coordinates different rendering components

#### Storage (`/internal/storage`)

Handles data persistence and caching:

- `db/`: Database operations and migrations
- `cache/`: In-memory and distributed caching
- `repository/`: Data access patterns for different game elements

#### Core Services (`/internal/core`)

Fundamental services used throughout the application:

- `config/`: Configuration management
- `events/`: Event bus for game state changes
- `logging/`: Centralized logging service

### Web Client (`/web`)

Contains the browser-based game client:

- `public/`: Static assets
- `src/`: Client-side application code and components

### Supporting Directories

- `scripts/`: Build and development automation
- `docs/`: Project documentation
- `configs/`: Environment-specific configurations
- `deployments/`: Deployment configurations
- `test/`: Integration and end-to-end tests

## Development Flow

1. Core game logic is implemented in Go
2. Performance-critical map generation is handled by Zig
3. Game state updates flow through the event system
4. Rendering and network updates are handled asynchronously

## Future Considerations

- Scaling WebSocket connections
- Adding more procedural generation features
- Implementing advanced particle effects
- Enhanced persistence and caching strategies

## Why This Structure?

This architecture is designed to:

1. Support incremental development
2. Maintain clear separation of concerns
3. Allow for easy testing and modification
4. Provide flexibility for future enhancements
5. Keep related code together while maintaining clear boundaries

The structure follows Go best practices while accommodating game-specific needs and the integration of Zig for performance-critical operations.
