package render_object

import (
	"math"
	"grt/ray"
	"grt/util"
	"grt/vec"
)

type Sphere struct {
	Center vec.Vec3
	Radius float64
	Material Material
}

func (s *Sphere) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	oc := vec.Subtract(s.Center, r.Origin())
	a := r.Direction().LenSquared()
	h := vec.Dot(r.Direction(), oc)
	c := oc.LenSquared() - s.Radius*s.Radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return false
	}

	sqrt_disc := math.Sqrt(discriminant)
	root := (h - sqrt_disc) / a

	if !interval.Surrounds(root) {
		root = (h + sqrt_disc) / a
		if !interval.Surrounds(root) {
			return false
		}
	}

	hit_record.T = root
	hit_record.Point = r.At(hit_record.T)
	outward_normal := vec.Subtract(hit_record.Point, s.Center).Scaled(1.0/s.Radius)
	hit_record.SetFaceNormal(r, outward_normal)
	hit_record.Material = s.Material

	return true
}
