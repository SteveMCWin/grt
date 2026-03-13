package render_object

import (
	"tracer/ray"
	"tracer/util"
)

type HittableList struct {
	Shapes []Hittable
}

func (list *HittableList) Add(shape Hittable) {
	list.Shapes = append(list.Shapes, shape)
}

func (list *HittableList) Hit(r *ray.Ray, interval util.Interval,hit_record *HitRecord) bool {
	var rec HitRecord
	hit_anything := false
	closest_so_far := interval.Max

	for _, shape := range list.Shapes {
		if shape.Hit(r, util.Interval{Min: interval.Min, Max: closest_so_far}, &rec) {
			hit_anything = true
			closest_so_far = rec.T
			*hit_record = rec
		}
	}

	return hit_anything
}
