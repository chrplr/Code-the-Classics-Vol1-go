package main

import (
	"flag"
	"fmt"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binmix"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

type State int

const (
	StateMenu State = iota
	StatePlay
	StateGameOver
)

var (
	state  State
	game   *Game
	assets *Assets
	audio  *Audio
)

func update() {
	switch state {
	case StateMenu:
		if spacePressed() {
			state = StatePlay
			game = NewGame(NewPlayer(), assets, audio)
		} else {
			game.Update()
		}

	case StatePlay:
		if game.player.lives < 0 {
			game.playSound("over", 1)
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
		// "Press SPACE" animation: 10 frames, holding on frame 9 most of the time.
		animFrame := ((game.timer + 40) % 160) / 4
		if animFrame > 9 {
			animFrame = 9
		}
		assets.Blit("space"+itoa(animFrame), 130, 280)

	case StatePlay:
		drawStatus(assets, game)

	case StateGameOver:
		drawStatus(assets, game)
		assets.Blit("over", 0, 0)
	}
}

func main() {
	selftest := flag.Bool("selftest", false, "run several levels headlessly, then exit")
	flag.Parse()

	defer binsdl.Load().Unload()
	defer binimg.Load().Unload()
	defer binmix.Load().Unload()
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("Cavern", Width, Height, 0)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	assets = NewAssets(renderer)
	defer assets.Destroy()
	audio = NewAudio()
	defer audio.Destroy()

	if *selftest {
		g := NewGame(NewPlayer(), assets, audio)
		// Free-running phase: exercises robots, bolts, fruit, gravity, collisions.
		for step := 0; step < 1500; step++ {
			g.Update()
		}
		// Orb phase: blow orbs and let them trap the robots, then float and pop.
		for step := 0; step < 1500; step++ {
			if step%30 == 0 && len(g.orbs) < 5 {
				o := NewOrb(g.player.X, g.player.Y-35, 1)
				g.orbs = append(g.orbs, o)
			}
			g.Update()
		}
		// Level-cycle phase: load and step every level layout in turn.
		for lvl := 0; lvl < len(LEVELS)*2; lvl++ {
			g.nextLevel()
			for step := 0; step < 120; step++ {
				g.Update()
			}
			fmt.Printf("level %d: %d grid rows, %d enemies, %d pending, %d fruits\n",
				g.level, len(g.grid), len(g.enemies), len(g.pendingEnemies), len(g.fruits))
		}
		// Verify embedded assets actually decode into real textures/sounds.
		loaded, total := 0, 0
		for _, n := range []string{"title", "over", "life", "block0", "still", "orb0", "space0"} {
			total++
			if w, h := assets.Size(n); w > 0 && h > 0 {
				loaded++
			}
		}
		fmt.Printf("embedded textures decoded: %d/%d; embedded sounds: %d\n", loaded, total, len(audio.sounds))
		fmt.Printf("SELFTEST OK: score %d, lives %d, health %d\n",
			g.player.score, g.player.lives, g.player.health)
		return
	}

	audio.PlayMusic(0.3)

	state = StateMenu
	game = NewGame(nil, assets, audio)

	sdl.RunLoop(func() error {
		frameStart := sdl.Ticks()

		var event sdl.Event
		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
			if event.Type == sdl.EVENT_KEY_DOWN && event.KeyboardEvent().Scancode == sdl.SCANCODE_ESCAPE {
				return sdl.EndLoop
			}
		}

		refreshKeys()

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		update()
		draw()

		renderer.Present()

		if elapsed := sdl.Ticks() - frameStart; elapsed < frameMillis {
			sdl.Delay(uint32(frameMillis - elapsed))
		}
		return nil
	})
}
