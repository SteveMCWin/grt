package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
)

type HitRecord struct {
	Point    vec.Vec3
	Normal   vec.Vec3
	T        float64
	Material Material
	Front    bool
}

func (hr *HitRecord) SetFaceNormal(r *ray.Ray, outward_normal vec.Vec3) {
	hr.Front = vec.Dot(r.Direction(), outward_normal) < 0
	hr.Normal = outward_normal
	if !hr.Front {
		hr.Normal = hr.Normal.Scaled(-1.0)
	}
}

type Hittable interface {
	Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool
	// NOTE: perhaps add material here
}
