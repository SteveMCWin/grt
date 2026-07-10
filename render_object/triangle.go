package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
)

type Vertex struct {
	Position vec.Vec3
	Normal   vec.Vec3
	TexCoord vec.Vec3
	Color    vec.Color
}

type Triangle struct {
	A Vertex
	B Vertex
	C Vertex

	Material    Material
	PlaneNormal vec.Vec3
}

func NewTriangle(a, b, c Vertex, m Material) *Triangle {
	vAC := a.Position.To(c.Position)
	vAB := a.Position.To(b.Position)
	planeNormal := vec.Cross(vAB, vAC).Normalized()

	if a.Normal.Len() == 0 {
		a.Normal = planeNormal
		b.Normal = planeNormal
		c.Normal = planeNormal
	}

	return &Triangle{
		A:           a,
		B:           b,
		C:           c,
		Material:    m,
		PlaneNormal: vec.Cross(vAB, vAC).Normalized(),
	}
}

// Möller-Trumbore algorithm
// https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-rendering-a-triangle//moller-trumbore-ray-triangle-intersection.html
func (tri *Triangle) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {

	vAB := tri.A.Position.To(tri.B.Position)
	vAC := tri.A.Position.To(tri.C.Position)
	vP := vec.Cross(r.Direction(), vAC)
	det := vec.Dot(vAB, vP)

	// if ray near parallel to plane, doesn't hit
	if det < util.Small {
		return false
	}

	invDet := 1.0 / det
	vT := vec.Subtract(r.Origin(), tri.A.Position)
	v := vec.Dot(vT, vP) * invDet

	// ray not in triangle (must be u, v, w >= 0 && u, v, w <= 1)
	if v < 0 || v > 1 {
		return false
	}

	vQ := vec.Cross(vT, vAB)
	w := vec.Dot(r.Direction(), vQ) * invDet
	// ray not in triangle (must be u + v + w = 1)
	if w < 0 || v+w > 1 {
		return false
	}

	u := 1.0 - v - w
	t := vec.Dot(vAC, vQ) * invDet

	if !interval.Surrounds(t) {
		return false
	}

	pointNormal := vec.Add(tri.A.Normal.Scaled(u), tri.B.Normal.Scaled(v), tri.C.Normal.Scaled(w))
	vertexColor := vec.Add(vec.Vec3(tri.A.Color).Scaled(u), vec.Vec3(tri.B.Color).Scaled(v), vec.Vec3(tri.C.Color).Scaled(w))

	hit_record.T = t
	hit_record.Point = r.At(t)
	hit_record.SetFaceNormal(r, pointNormal)
	hit_record.Material = tri.Material
	hit_record.VertexColor = vec.Color(vertexColor)

	return true
}
