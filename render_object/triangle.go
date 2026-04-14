package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
	// "log"
)

type Triangle struct {
	A vec.Vec3
	B vec.Vec3
	C vec.Vec3

	N vec.Vec3
	Material Material
}

func (t *Triangle) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {

	vAC := t.A.To(t.C)
	vAB := t.A.To(t.B)
	vBC := t.B.To(t.C)
	vCA := t.C.To(t.A)

	N := vec.Cross(vAB, vAC).Normalized()
	D := -vec.Dot(N, t.A)

	// l for length of ray, but usually labeled 't'
	l := -(vec.Dot(N, r.Origin()) + D)/vec.Dot(N, r.Direction())

	// Point where ray hits plane
	P := r.At(l)

	vAP := t.A.To(P)
	vBP := t.B.To(P)
	vCP := t.C.To(P)

	// If intersection point isn't inside the triangle, it didn't hit
	if (vec.Dot(N, vec.Cross(vAB, vAP)) <= 0 ||
		vec.Dot(N, vec.Cross(vBC, vBP)) <= 0 ||
		vec.Dot(N, vec.Cross(vCA, vCP)) <= 0) {
		return false
	}

	if !interval.Surrounds(l) {
		// log.Println("!interval.Surrounds(P.Len())")
		return false
	}

	hit_record.T = l
	hit_record.Point = P
	outward_normal := N
	hit_record.SetFaceNormal(r, outward_normal)
	hit_record.Material = t.Material

	return true
}
