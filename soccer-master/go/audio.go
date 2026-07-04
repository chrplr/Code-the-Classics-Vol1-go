package main

import (
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
)

// Audio wraps SDL3_mixer. Every operation is best-effort: if the mixer cannot be
// created (no audio device) all methods become no-ops, matching the original
// game's silent handling of sound errors.
type Audio struct {
	mixer  *mixer.Mixer
	sounds map[string]*mixer.Audio
	music  *mixer.Track // looping theme (menu)
	crowd  *mixer.Track // looping crowd (in-game)
}

func NewAudio(assetDir string) *Audio {
	a := &Audio{sounds: make(map[string]*mixer.Audio)}

	if err := mixer.Init(); err != nil {
		return a
	}
	m, err := mixer.CreateMixerDevice(sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK, nil)
	if err != nil {
		return a
	}
	a.mixer = m

	matches, _ := filepath.Glob(filepath.Join(assetDir, "sounds", "*.ogg"))
	for _, path := range matches {
		name := filepath.Base(path)
		name = name[:len(name)-len(filepath.Ext(name))]
		if snd, err := m.LoadAudio(path, true); err == nil {
			a.sounds[name] = snd
		}
	}

	a.music = a.loopingTrack(m, filepath.Join(assetDir, "music", "theme.ogg"), 0.5)
	// The crowd loop is one of the sound effects rather than music.
	if snd, ok := a.sounds["crowd"]; ok {
		if t, err := m.CreateTrack(); err == nil {
			t.SetAudio(snd)
			t.SetLoops(-1)
			a.crowd = t
		}
	}
	return a
}

func (a *Audio) loopingTrack(m *mixer.Mixer, path string, gain float32) *mixer.Track {
	audio, err := m.LoadAudio(path, false)
	if err != nil {
		return nil
	}
	t, err := m.CreateTrack()
	if err != nil {
		return nil
	}
	t.SetAudio(audio)
	t.SetLoops(-1)
	t.SetGain(gain)
	return t
}

// PlaySound plays one of a family of variants: <name>0 .. <name>(count-1).
func (a *Audio) PlaySound(name string, count int) {
	a.play(name + strconv.Itoa(rand.Intn(count)))
}

// Play plays a single sound by exact name (e.g. "start", "move").
func (a *Audio) Play(name string) { a.play(name) }

func (a *Audio) play(key string) {
	if a.mixer == nil {
		return
	}
	if snd, ok := a.sounds[key]; ok {
		a.mixer.PlayAudio(snd)
	}
}

// StartMenuMusic plays the looping theme and stops the crowd.
func (a *Audio) StartMenuMusic() {
	if a.music != nil {
		a.music.Play(0)
	}
	if a.crowd != nil {
		a.crowd.Stop(0)
	}
}

// StartMatchAudio fades out the theme, starts the crowd loop and plays the whistle.
func (a *Audio) StartMatchAudio() {
	if a.music != nil {
		a.music.Stop(0)
	}
	if a.crowd != nil {
		a.crowd.Play(0)
	}
	a.Play("start")
}

func (a *Audio) Destroy() {
	if a.mixer != nil {
		a.mixer.Destroy()
	}
}
