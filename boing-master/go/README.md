# Boing! — Go / SDL3 port

An idiomatic Go re-implementation of the Pygame Zero game in `../boing.py`, using
[`go-sdl3`](https://github.com/Zyko0/go-sdl3) (SDL3 bindings via `purego`, no CGo).

The game reuses the original assets in `../images`, `../sounds`, `../music`.

## Run

```sh
go run .            # from this directory; assets resolve to ..
go run . -assets ..  # explicit asset root
```

`go-sdl3` bundles the SDL3, SDL3_image and SDL3_mixer shared libraries and
extracts them to a temp directory at startup, so no system SDL3 install is needed.

## Controls

- **Player 1:** `A`/`Up` up, `Z`/`Down` down
- **Player 2** (two-player mode): `K` up, `M` down
- **Menu:** `Up`/`Down` choose 1 or 2 players, `Space` starts
- **Game over:** `Space` returns to the menu

## Layout

| File | Responsibility |
|------|----------------|
| `main.go` | SDL init, window/renderer, 60 FPS loop, menu/play/over state machine |
| `assets.go` / `sprite.go` | texture cache; centred (actor) and top-left (UI) blits |
| `audio.go` | SDL3_mixer wrapper: preloaded SFX + looping theme (best-effort) |
| `input.go` | per-frame keyboard snapshot and control functions |
| `ball.go`, `bat.go`, `impact.go`, `game.go` | game entities and rules |

An embedded `Sprite` struct replaces Pygame Zero's `Actor` (Go has no inheritance);
`Ball`, `Bat` and `Impact` embed it.
