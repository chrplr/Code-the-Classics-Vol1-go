package main

import (
	"embed"
	"math/rand"
	"path"
	"strconv"

	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
)

// audioFS embeds the sound effects and music into the binary.
//
//go:embed sounds music
var audioFS embed.FS

// Audio wraps SDL3_mixer. All operations are best-effort: if the mixer cannot be
// created (e.g. no audio device) every method becomes a no-op, mirroring the
// original game which silently ignores sound errors.
type Audio struct {
	mixer  *mixer.Mixer
	sounds map[string]*mixer.Audio
	music  *mixer.Track
}

// NewAudio initialises the mixer and preloads every sound in sounds/ and the
// looping theme in music/. It never returns an error; on failure it returns a
// no-op Audio.
func NewAudio() *Audio {
	a := &Audio{sounds: make(map[string]*mixer.Audio)}

	if err := mixer.Init(); err != nil {
		return a
	}
	m, err := mixer.CreateMixerDevice(sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK, nil)
	if err != nil {
		return a
	}
	a.mixer = m

	// Preload every embedded .ogg sound effect, keyed by filename without extension.
	entries, _ := audioFS.ReadDir("sounds")
	for _, e := range entries {
		fname := e.Name()
		if path.Ext(fname) != ".ogg" {
			continue
		}
		if snd := loadAudioFromFS(m, "sounds/"+fname); snd != nil {
			a.sounds[fname[:len(fname)-len(".ogg")]] = snd
		}
	}

	// Load and start the looping theme at a low volume.
	if themeAudio := loadAudioFromFS(m, "music/theme.ogg"); themeAudio != nil {
		if track, err := m.CreateTrack(); err == nil {
			track.SetAudio(themeAudio)
			track.SetLoops(-1)
			track.SetGain(0.3)
			track.Play(0)
			a.music = track
		}
	}

	return a
}

// loadAudioFromFS decodes an embedded audio file into an in-memory Audio via an
// SDL IOStream (predecoded, so no stream stays open afterwards).
func loadAudioFromFS(m *mixer.Mixer, p string) *mixer.Audio {
	data, err := audioFS.ReadFile(p)
	if err != nil {
		return nil
	}
	stream, err := sdl.IOFromConstMem(data)
	if err != nil {
		return nil
	}
	snd, err := m.LoadAudio_IO(stream, true, true) // predecode + closeio
	if err != nil {
		return nil
	}
	return snd
}

// PlaySound plays one of a family of sound variants. count > 1 chooses randomly
// among <name>0 .. <name>(count-1); count == 1 plays <name>0.
func (a *Audio) PlaySound(name string, count int) {
	if a.mixer == nil {
		return
	}
	variant := name + strconv.Itoa(rand.Intn(count))
	if snd, ok := a.sounds[variant]; ok {
		a.mixer.PlayAudio(snd)
	}
}

func (a *Audio) Destroy() {
	if a.mixer != nil {
		a.mixer.Destroy()
	}
}
