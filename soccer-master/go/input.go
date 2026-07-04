package main

import "github.com/Zyko0/go-sdl3/sdl"

// keys holds a snapshot of the keyboard state for the current frame.
var keys []bool

// keyStatus records each tracked key's state at its last keyJustPressed check,
// mirroring the original's key_status dictionary.
var keyStatus = map[sdl.Scancode]bool{}

func refreshKeys() {
	keys = sdl.GetKeyboardState()
}

func keyDown(sc sdl.Scancode) bool {
	return keys != nil && int(sc) < len(keys) && keys[sc]
}

// keyJustPressed reports whether sc is down now but wasn't at its previous check.
// As in the original, it updates the stored state on every call.
func keyJustPressed(sc sdl.Scancode) bool {
	prev := keyStatus[sc]
	cur := keyDown(sc)
	keyStatus[sc] = cur
	return !prev && cur
}

// Controls maps one player's keys and reports movement/shoot intents.
type Controls struct {
	up, down, left, right, shootKey sdl.Scancode
}

func NewControls(playerNum int) *Controls {
	if playerNum == 0 {
		return &Controls{
			up: sdl.SCANCODE_UP, down: sdl.SCANCODE_DOWN,
			left: sdl.SCANCODE_LEFT, right: sdl.SCANCODE_RIGHT,
			shootKey: sdl.SCANCODE_SPACE,
		}
	}
	return &Controls{
		up: sdl.SCANCODE_W, down: sdl.SCANCODE_S,
		left: sdl.SCANCODE_A, right: sdl.SCANCODE_D,
		shootKey: sdl.SCANCODE_LSHIFT,
	}
}

// move returns the movement vector for the held direction keys, scaled by speed.
func (c *Controls) move(speed float64) Vec2 {
	dx, dy := 0.0, 0.0
	if keyDown(c.left) {
		dx = -1
	} else if keyDown(c.right) {
		dx = 1
	}
	if keyDown(c.up) {
		dy = -1
	} else if keyDown(c.down) {
		dy = 1
	}
	return Vec2{dx, dy}.Mul(speed)
}

func (c *Controls) shoot() bool {
	return keyJustPressed(c.shootKey)
}
