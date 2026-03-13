package ray

import (
	// "tracer/shapes"
	"tracer/vec"
)

type Ray struct {
	origin    vec.Vec3
	direction vec.Vec3
}

func NewRay(orig, dir vec.Vec3) Ray {
	r := Ray{
		origin:    orig,
		direction: dir.Normalized(),
	}

	return r
}

func (r *Ray) Origin() vec.Vec3 {
	return r.origin
}

func (r *Ray) Direction() vec.Vec3 {
	return r.direction
}

func (r *Ray) At(t float64) vec.Vec3 {
	return vec.Add(r.origin, r.direction.Scaled(t))
}
