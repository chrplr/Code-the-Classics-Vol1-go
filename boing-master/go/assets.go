package main

import (
	"embed"

	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
)

// imagesFS embeds every PNG in images/ into the binary, so the game is fully
// self-contained and needs no asset files at run time.
//
//go:embed images
var imagesFS embed.FS

// Assets caches textures decoded from the embedded images and provides the two
// blit modes we need: top-left anchored (for backgrounds and UI) and centre
// anchored (for actors).
type Assets struct {
	renderer *sdl.Renderer
	textures map[string]*sdl.Texture
}

func NewAssets(renderer *sdl.Renderer) *Assets {
	return &Assets{
		renderer: renderer,
		textures: make(map[string]*sdl.Texture),
	}
}

// Texture lazily decodes the embedded images/<name>.png and caches it. A load
// failure returns nil, which the blit helpers treat as a no-op so a single
// missing sprite never crashes the game.
func (a *Assets) Texture(name string) *sdl.Texture {
	if tex, ok := a.textures[name]; ok {
		return tex
	}
	a.textures[name] = loadTextureFromFS(a.renderer, imagesFS, "images/"+name+".png")
	return a.textures[name]
}

// loadTextureFromFS decodes an embedded image file into a texture via an
// in-memory SDL IOStream (no filesystem access).
func loadTextureFromFS(renderer *sdl.Renderer, fs embed.FS, path string) *sdl.Texture {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil
	}
	stream, err := sdl.IOFromConstMem(data)
	if err != nil {
		return nil
	}
	tex, err := img.LoadTextureIO(renderer, stream, true) // closeio: frees the stream
	if err != nil {
		return nil
	}
	return tex
}

// Blit draws an image with its top-left corner at (x, y), like Pygame Zero's
// screen.blit. Used for the table, menu, game-over and digit sprites.
func (a *Assets) Blit(name string, x, y float64) {
	tex := a.Texture(name)
	if tex == nil {
		return
	}
	dst := sdl.FRect{X: float32(x), Y: float32(y), W: float32(tex.W), H: float32(tex.H)}
	a.renderer.RenderTexture(tex, nil, &dst)
}

// BlitCentred draws an image centred on (x, y), like an Actor.draw().
func (a *Assets) BlitCentred(name string, x, y float64) {
	tex := a.Texture(name)
	if tex == nil {
		return
	}
	a.Blit(name, x-float64(tex.W)/2, y-float64(tex.H)/2)
}

// Destroy frees every cached texture.
func (a *Assets) Destroy() {
	for _, tex := range a.textures {
		if tex != nil {
			tex.Destroy()
		}
	}
}
