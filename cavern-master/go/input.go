package main

import "github.com/Zyko0/go-sdl3/sdl"

// keys holds a per-frame snapshot of the keyboard state.
var keys []bool

func refreshKeys() {
	current := sdl.GetKeyboardState()
	if current == nil {
		return
	}
	if len(keys) != len(current) {
		keys = make([]bool, len(current))
	}
	copy(keys, current)
}

func keyDown(sc sdl.Scancode) bool {
	return keys != nil && int(sc) < len(keys) && keys[sc]
}

// Held-key helpers mirroring Pygame Zero's keyboard.left / .right / .up / .space.
func keyLeft() bool  { return keyDown(sdl.SCANCODE_LEFT) }
func keyRight() bool { return keyDown(sdl.SCANCODE_RIGHT) }
func keyUp() bool    { return keyDown(sdl.SCANCODE_UP) }
func keySpace() bool { return keyDown(sdl.SCANCODE_SPACE) }

// spaceDown latches the previous space-bar state, exactly like the Python global
// of the same name. spacePressed must be called at most once per frame.
var spaceDown bool

// spacePressed reports whether space was just pressed this frame (edge trigger).
func spacePressed() bool {
	if keySpace() {
		if spaceDown {
			return false
		}
		spaceDown = true
		return true
	}
	spaceDown = false
	return false
}
