package main

import (
	"flag"
	"strconv"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binmix"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

const (
	Width  = 480
	Height = 800

	targetFPS   = 60
	frameMillis = 1000 / targetFPS
)

// --- small numeric helpers ---

func abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// pmod is Python-style modulo: the result has the sign of the divisor.
func pmod(a, b int) int {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}

func floorDiv(a, b int) int {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

// b converts a bool to 0/1; b2s to "0"/"1".
func b(v bool) int {
	if v {
		return 1
	}
	return 0
}

func b2s(v bool) string { return strconv.Itoa(b(v)) }

// --- game state ---

type State int

const (
	StateMenu State = iota
	StatePlay
	StateGameOver
)

var (
	state     State
	game      *Game
	spaceDown bool

	assets *Assets
	audio  *Audio
)

// spacePressed reports a fresh press of the space bar (down now, up last frame).
func spacePressed() bool {
	if keyDown_space() {
		if spaceDown {
			return false
		}
		spaceDown = true
		return true
	}
	spaceDown = false
	return false
}

func update() {
	switch state {
	case StateMenu:
		if spacePressed() {
			state = StatePlay
			game = NewGame(NewPlayer(240, 768), assets, audio)
		}
		game.Update()

	case StatePlay:
		if game.player.lives == 0 && game.player.timer == 100 {
			audio.Play("gameover")
			state = StateGameOver
		} else {
			game.Update()
		}

	case StateGameOver:
		if spacePressed() {
			state = StateMenu
			game = NewGame(nil, assets, audio)
		}
	}
}

func draw() {
	game.Draw()

	switch state {
	case StateMenu:
		assets.Blit("title", 0, 0)
		assets.Blit("space"+strconv.Itoa((game.time/4)%14), 0, 420)

	case StatePlay:
		for i := 0; i < game.player.lives; i++ {
			assets.Blit("life", float64(i*40+8), 4)
		}
		score := strconv.Itoa(game.score)
		for i := 1; i <= len(score); i++ {
			digit := string(score[len(score)-i])
			assets.Blit("digit"+digit, float64(468-i*24), 5)
		}

	case StateGameOver:
		assets.Blit("over", 0, 0)
	}
}

func main() {
	assetDir := flag.String("assets", "..", "directory containing images/, sounds/, music/")
	flag.Parse()

	defer binsdl.Load().Unload()
	defer binimg.Load().Unload()
	defer binmix.Load().Unload()
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("Myriapod", Width, Height, 0)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	assets = NewAssets(renderer, *assetDir)
	defer assets.Destroy()
	audio = NewAudio(*assetDir)
	defer audio.Destroy()

	state = StateMenu
	// The initial game has no player: it runs as the attract-mode demo.
	game = NewGame(nil, assets, audio)

	sdl.RunLoop(func() error {
		frameStart := sdl.Ticks()

		var event sdl.Event
		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
		}

		refreshKeys()
		update()

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		draw()
		renderer.Present()

		if elapsed := sdl.Ticks() - frameStart; elapsed < frameMillis {
			sdl.Delay(uint32(frameMillis - elapsed))
		}
		return nil
	})
}
