package render_object

import (
	"grt/ray"
	"grt/util"
	"math"
	"math/rand/v2"
)

type Mist struct {
	Boundary      Hittable
	Material      Material
	NegInvDensity float64
}

func (m *Mist) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	rec1 := &HitRecord{}
	rec2 := &HitRecord{}
	if !m.Boundary.Hit(r, util.EverythingInterval, rec1) {
		return false
	}

	if !m.Boundary.Hit(r, util.Interval{Min: rec1.T + 0.0001, Max: util.PosInf}, rec2) {
		return false
	}

	rec1.T = max(rec1.T, interval.Min)
	rec2.T = min(rec2.T, interval.Max)
	
	if rec1.T >= rec2.T {
		return false
	}

	rec1.T = max(rec1.T, 0.0)

	rayLen := r.Direction().Len()
	distInsideBoundary := (rec2.T - rec1.T) * rayLen
	hitDistance := m.NegInvDensity * math.Log(rand.Float64())

	if hitDistance > distInsideBoundary {
		return false
	}

	hit_record.T = rec1.T + (hitDistance / rayLen)
	hit_record.Point = r.At(hit_record.T)
	hit_record.Material = m.Material

	return true
}
