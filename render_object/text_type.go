package render_object

import (
	"grt/vec"
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
	// *texture
}

// func (tcp *TextureColorProvider) GetColor(rec *HitRecord) vec.Color {
	// pX = tcp.Texture.Width  * rec.U
	// pY = tcp.Texture.Height * rec.V
	// return tcp.Texture.At(pX, pY)
// }
