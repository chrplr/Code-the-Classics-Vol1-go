package main

import (
	"embed"
	"path"
	"strconv"

	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
)

// audioFS embeds the sound effects and music into the binary.
//
//go:embed sounds music
var audioFS embed.FS

// Audio wraps SDL3_mixer. All operations are best-effort.
type Audio struct {
	mixer  *mixer.Mixer
	sounds map[string]*mixer.Audio

	music        *mixer.Track
	currentMusic *mixer.Track
}

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

	entries, _ := audioFS.ReadDir("sounds")
	for _, e := range entries {
		fname := e.Name()
		if path.Ext(fname) != ".ogg" {
			continue
		}
		name := fname[:len(fname)-len(".ogg")]
		if snd := loadAudioFromFS(m, "sounds/"+fname); snd != nil {
			a.sounds[name] = snd
		}
	}

	a.music = a.loopingTrack(m, "music/theme.ogg")
	return a
}

// loadAudioFromFS decodes an embedded audio file into an in-memory Audio via an
// SDL IOStream (predecoded, so no filesystem or stream stays open afterwards).
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

func (a *Audio) loopingTrack(m *mixer.Mixer, p string) *mixer.Track {
	audio := loadAudioFromFS(m, p)
	if audio == nil {
		return nil
	}
	t, err := m.CreateTrack()
	if err != nil {
		return nil
	}
	t.SetAudio(audio)
	t.SetLoops(-1)
	return t
}

// PlaySound plays one of a family of sound variants (name0 .. name(count-1)).
func (a *Audio) PlaySound(name string, count int) {
	if a.mixer == nil {
		return
	}
	variant := name + "0"
	if count > 1 {
		variant = name + strconv.Itoa(randIntn(count))
	}
	if snd, ok := a.sounds[variant]; ok && snd != nil {
		a.mixer.PlayAudio(snd)
	}
}

// PlayMusic starts the looping theme at the given gain.
func (a *Audio) PlayMusic(volume float32) {
	if a.music == nil {
		return
	}
	a.music.SetGain(volume)
	a.music.Play(0)
	a.currentMusic = a.music
}

func (a *Audio) Destroy() {
	if a.mixer != nil {
		a.mixer.Destroy()
	}
}
