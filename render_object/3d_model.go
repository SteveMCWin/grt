package render_object

import (
	"grt/ray"
	"grt/util"
	// "grt/vec"
)

type Model struct {
	Triangles []Triangle
	Material  Material
}

func (m *Model) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	// TODO

	return false
}
