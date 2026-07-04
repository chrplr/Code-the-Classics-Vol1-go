# Myriapod — Go / SDL3 port

An idiomatic Go re-implementation of the Pygame Zero game in `../myriapod.py`
(a Centipede-style arcade game), using [`go-sdl3`](https://github.com/Zyko0/go-sdl3)
(SDL3 bindings via `purego`, no CGo).

The game reuses the original assets in `../images`, `../sounds`, `../music`.

## Run

```sh
go run .             # from this directory; assets resolve to ..
go run . -assets ..  # explicit asset root
```

`go-sdl3` bundles SDL3, SDL3_image and SDL3_mixer and extracts them to a temp
directory at startup, so no system SDL3 install is needed.

## Controls

- **Arrow keys:** move the ship (within the bottom zone)
- **Space:** fire / start / continue

## Layout

| File | Responsibility |
|------|----------------|
| `main.go` | SDL init, window/renderer, 60 FPS loop, menu/play/over state machine, numeric helpers |
| `assets.go` / `sprite.go` | texture cache; anchor-aware drawing and `collidepoint` |
| `audio.go` | SDL3_mixer wrapper: preloaded SFX + looping theme (best-effort) |
| `input.go` | per-frame keyboard snapshot |
| `game.go` | grid, wave logic, damage, movement rules, depth-sorted drawing |
| `player.go` | player movement, gradual turning, firing, collisions |
| `segment.go` | myriapod segment movement (16-phase cell cycle, edge ranking, sprite direction) |
| `rock.go`, `bullet.go`, `explosion.go`, `flyingenemy.go` | remaining entities |

### Notes on the port
- Go has no inheritance, so an embedded `Sprite` (position + image + anchor)
  replaces Pygame Zero's `Actor`. Tall "totem" rocks use a custom anchor.
- The `occupied` set (which mixes 2- and 3-tuples in Python) becomes a
  `map[[3]int]bool`, with `-1` in the third slot standing in for the 2-tuple form.
- Segment edge selection packs the seven ordering factors into one integer so the
  lowest value wins, reproducing Python's `min(range(4), key=...)` tuple compare.
