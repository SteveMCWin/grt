package render_object

import (
	"grt/ray"
	"grt/util"
)

type HittableList struct {
	Shapes []Hittable
}

func (list *HittableList) Add(shape Hittable) {
	list.Shapes = append(list.Shapes, shape)
}

func (list *HittableList) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	hit_anything := false
	closest_so_far := interval.Max

	for _, shape := range list.Shapes {
		if shape.Hit(r, util.Interval{Min: interval.Min, Max: closest_so_far}, hit_record) {
			hit_anything = true
			closest_so_far = hit_record.T
		}
	}

	return hit_anything
}
