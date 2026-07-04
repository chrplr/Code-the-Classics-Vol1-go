package main

import (
	"path/filepath"

	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
)

// Assets caches textures loaded from the images directory.
type Assets struct {
	renderer *sdl.Renderer
	dir      string
	textures map[string]*sdl.Texture
}

func NewAssets(renderer *sdl.Renderer, assetDir string) *Assets {
	return &Assets{
		renderer: renderer,
		dir:      filepath.Join(assetDir, "images"),
		textures: make(map[string]*sdl.Texture),
	}
}

func (a *Assets) Texture(name string) *sdl.Texture {
	if tex, ok := a.textures[name]; ok {
		return tex
	}
	tex, err := img.LoadTexture(a.renderer, filepath.Join(a.dir, name+".png"))
	if err != nil {
		tex = nil
	}
	a.textures[name] = tex
	return tex
}

func (a *Assets) Size(name string) (float64, float64) {
	tex := a.Texture(name)
	if tex == nil {
		return 0, 0
	}
	return float64(tex.W), float64(tex.H)
}

// Blit draws an image with its top-left corner at (x, y), like screen.blit.
func (a *Assets) Blit(name string, x, y float64) {
	tex := a.Texture(name)
	if tex == nil {
		return
	}
	dst := sdl.FRect{X: float32(x), Y: float32(y), W: float32(tex.W), H: float32(tex.H)}
	a.renderer.RenderTexture(tex, nil, &dst)
}

func (a *Assets) Destroy() {
	for _, tex := range a.textures {
		if tex != nil {
			tex.Destroy()
		}
	}
}
