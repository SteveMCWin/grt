package render_object

import (
	"grt/vec"
	"image"
	"math"
)

type TexType interface {
	GetColor(rec *HitRecord) vec.Color
}

type DebugColorType int

const (
	VertexNormal = iota
	VertexColor
	VertexUV
	WireMesh
)

type DebugColorProvider struct {
	Type DebugColorType
}

func (dcp *DebugColorProvider) GetColor(rec *HitRecord) vec.Color {
	res := vec.Color{1.0, 0.0, 1.0}

	switch dcp.Type {
	case VertexNormal:
		res = vec.Color(rec.Normal)
	case VertexColor:
		res = rec.VertexColor
	case VertexUV:
		res = vec.Color{rec.U, rec.V, 0.0}
	}

	return res
}

type FlatColorProvider struct {
	Color vec.Color
}

func (fcp *FlatColorProvider) GetColor(rec *HitRecord) vec.Color {
	return fcp.Color
}

type TextureColorProvider struct {
	Texture image.Image
}

func (tcp *TextureColorProvider) GetColor(rec *HitRecord) vec.Color {
	pX := float64(tcp.Texture.Bounds().Max.X)  * rec.U
	pY := float64(tcp.Texture.Bounds().Max.Y) * (1.0 - rec.V) // y is inverted
	r, g, b, _ := tcp.Texture.At(int(pX), int(pY)).RGBA()
	// Note: r, g and b are all uint 32 but values range from 0 to 0xffff per the docs
	fr := float64(r)/float64(math.MaxUint16)
	fg := float64(g)/float64(math.MaxUint16)
	fb := float64(b)/float64(math.MaxUint16)
	return vec.Color{fr, fg, fb}
}
