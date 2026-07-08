package render_object

import (
	"math"
	"math/rand/v2"
	"grt/ray"
	"grt/vec"
)

type Material interface {
	Scatter(r_in *ray.Ray, rec *HitRecord, attenuation *vec.Color, r_out *ray.Ray) bool
}



type Lambertian struct {
	Texture TexType
}

func (l *Lambertian) Scatter(r_in *ray.Ray, rec *HitRecord, attenuation *vec.Color, r_out *ray.Ray) bool {
	scatter_direction := vec.Add(rec.Normal, vec.RandUnit())

	if (scatter_direction.NearZero()) {
		scatter_direction = rec.Normal
	}

	*r_out = ray.NewRay(rec.Point, scatter_direction)
	*attenuation = l.Texture.GetColor(rec)
	return true
}



type Metal struct {
	Texture TexType
	Fuzz float64
}

func (m *Metal) Scatter(r_in *ray.Ray, rec *HitRecord, attenuation *vec.Color, r_out *ray.Ray) bool {
	reflected := r_in.Direction().Reflected(rec.Normal)
	reflected = vec.Add(reflected.Normalized(), vec.RandUnit().Scaled(m.Fuzz))
	*r_out = ray.NewRay(rec.Point, reflected)
	*attenuation = m.Texture.GetColor(rec)

	return vec.Dot(r_out.Direction(), rec.Normal) > 0
}

type Dielectric struct {
	RefractionIndex float64
}

func (d *Dielectric) Scatter(r_in *ray.Ray, rec *HitRecord, attenuation *vec.Color, r_out *ray.Ray) bool {
	// TODO: i'd first have to experiment with this?
	*attenuation = vec.Color{1.0, 1.0, 1.0}

	var reflection_index float64

	if rec.Front {
		reflection_index = 1.0/d.RefractionIndex
	} else {
		reflection_index = d.RefractionIndex
	}

	cos_theta := math.Min(vec.Dot(r_in.Direction().Neg(), rec.Normal), 1.0)
	sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

	cannot_refract := reflection_index * sin_theta > 1.0
	var direction vec.Vec3

	if cannot_refract || reflectance(cos_theta, reflection_index) > rand.Float64() {
		direction = r_in.Direction().Reflected(rec.Normal)
	} else {
		direction = r_in.Direction().Refracted(rec.Normal, reflection_index)
	}

	*r_out = ray.NewRay(rec.Point, direction)

	return true
}

func reflectance(cosine, refraction_index float64) float64 {
	r0 := (1.0 - refraction_index) / (1 + refraction_index)
	r0 = r0*r0
	return r0 + (1.0-r0)*math.Pow((1-cosine), 5)
}
