package main

import "github.com/Zyko0/go-sdl3/sdl"

// keys holds a snapshot of the keyboard state for the current frame, refreshed
// once per frame by refreshKeys. It is indexed by scancode.
var keys []bool

func refreshKeys() {
	keys = sdl.GetKeyboardState()
}

func keyDown(sc sdl.Scancode) bool {
	return keys != nil && int(sc) < len(keys) && keys[sc]
}

func keyDown_left() bool  { return keyDown(sdl.SCANCODE_LEFT) }
func keyDown_right() bool { return keyDown(sdl.SCANCODE_RIGHT) }
func keyDown_up() bool    { return keyDown(sdl.SCANCODE_UP) }
func keyDown_down() bool  { return keyDown(sdl.SCANCODE_DOWN) }
func keyDown_space() bool { return keyDown(sdl.SCANCODE_SPACE) }
