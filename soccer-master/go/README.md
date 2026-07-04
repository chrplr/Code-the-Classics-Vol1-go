# Substitute Soccer — Go / SDL3 port

An idiomatic Go re-implementation of the Pygame Zero game in `../soccer.py`, using
[`go-sdl3`](https://github.com/Zyko0/go-sdl3) (SDL3 bindings via `purego`, no CGo).
This is the most complex of the Code the Classics ports: a scrolling 2D pitch with
full team AI.

The game reuses the original assets in `../images`, `../sounds`, `../music`.

## Run

```sh
go run .             # from this directory; assets resolve to ..
go run . -assets ..  # explicit asset root
```

`go-sdl3` bundles SDL3, SDL3_image and SDL3_mixer and extracts them at startup, so
no system SDL3 install is needed.

## Controls

- **Player 1:** arrow keys to move, **Space** to shoot / switch player
- **Player 2** (two-player mode): **W/A/S/D** to move, **Left Shift** to shoot / switch
- **Menu:** Up/Down to change option, **Space** to confirm

First team to 9 goals wins.

## Layout

| File | Responsibility |
|------|----------------|
| `main.go` | SDL init, 60 FPS loop, menu/play/over state machine, `pmod`/`floorDiv` |
| `vec.go` | `Vec2` and `safeNormalise` (replacing pygame's `Vector2`) |
| `geom.go` | level/pitch/goal geometry, difficulty table, angle math, `costAt`, ball physics |
| `actor.go` / `assets.go` | camera-offset drawing with per-sprite anchors; texture cache |
| `audio.go` | SDL3_mixer wrapper: theme music, looping crowd, one-shot SFX |
| `input.go` | keyboard snapshot, `keyJustPressed`, per-player `Controls` |
| `ball.go` | ball physics, dribbling, ownership, kick/pass target selection |
| `player.go` | the team AI: run targets, marking, interception, `targetable` |
| `game.go` | kickoff setup, marking/lead assignment, camera, depth-sorted drawing |
| `goal.go` | `Goal`, `Team`, and the `posTeam` / `Marker` interfaces |

### Notes on the port
- Go has no inheritance, so an embedded `Actor` (world `vpos` + image + anchor)
  replaces the game's `MyActor`; players use a feet anchor of (25, 37).
- Pass targets and marking assignments (which may be a `Player` **or** a `Goal`)
  are handled with the small `posTeam` and `Marker` interfaces plus type switches,
  standing in for Python's `isinstance` checks.
- The pursuer-selection `zip(a+[None,None], b+[None,None])` idiom is reproduced by
  `zipInterleave`, and `key_just_pressed`'s shared per-key state is preserved so
  that pressing shoot kicks (when you own the ball) or switches player (when you
  don't), exactly as in the original.
