package main

import (
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
)

// Audio wraps SDL3_mixer. Every operation is best-effort: if the mixer cannot be
// created (e.g. no audio device) all methods become no-ops, mirroring the
// original game's silent handling of sound errors.
type Audio struct {
	mixer  *mixer.Mixer
	sounds map[string]*mixer.Audio
	music  *mixer.Track
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

	// Preload every .ogg sound effect, keyed by filename without extension.
	matches, _ := filepath.Glob(filepath.Join(assetDir, "sounds", "*.ogg"))
	for _, path := range matches {
		name := filepath.Base(path)
		name = name[:len(name)-len(filepath.Ext(name))]
		if snd, err := m.LoadAudio(path, true); err == nil {
			a.sounds[name] = snd
		}
	}

	// Load and start the looping theme at a low volume.
	if themeAudio, err := m.LoadAudio(filepath.Join(assetDir, "music", "theme.ogg"), false); err == nil {
		if track, err := m.CreateTrack(); err == nil {
			track.SetAudio(themeAudio)
			track.SetLoops(-1)
			track.SetGain(0.4)
			track.Play(0)
			a.music = track
		}
	}

	return a
}

// PlaySound plays one of a family of sound variants: <name>0 .. <name>(count-1).
func (a *Audio) PlaySound(name string, count int) {
	a.play(name + strconv.Itoa(rand.Intn(count)))
}

// Play plays a single sound by its exact name (e.g. "gameover").
func (a *Audio) Play(name string) {
	a.play(name)
}

func (a *Audio) play(key string) {
	if a.mixer == nil {
		return
	}
	if snd, ok := a.sounds[key]; ok {
		a.mixer.PlayAudio(snd)
	}
}

func (a *Audio) Destroy() {
	if a.mixer != nil {
		a.mixer.Destroy()
	}
}
