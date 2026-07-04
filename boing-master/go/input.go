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

// p1Controls: Z or Down moves down, A or Up moves up.
func p1Controls() float64 {
	switch {
	case keyDown(sdl.SCANCODE_Z) || keyDown(sdl.SCANCODE_DOWN):
		return PlayerSpeed
	case keyDown(sdl.SCANCODE_A) || keyDown(sdl.SCANCODE_UP):
		return -PlayerSpeed
	}
	return 0
}

// p2Controls: M moves down, K moves up.
func p2Controls() float64 {
	switch {
	case keyDown(sdl.SCANCODE_M):
		return PlayerSpeed
	case keyDown(sdl.SCANCODE_K):
		return -PlayerSpeed
	}
	return 0
}
