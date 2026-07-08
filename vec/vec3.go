package vec

import (
	"fmt"
	"math"
	"math/rand/v2"
	"grt/util"
)

type Vec3 [3]float64

func (v Vec3) SetCoords(x, y, z float64) {
	v[0] = x
	v[1] = y
	v[2] = z
}

func (v Vec3) X() float64 {
	return v[0]
}

func (v Vec3) Y() float64 {
	return v[1]
}

func (v Vec3) Z() float64 {
	return v[2]
}

func Subtract(v1 Vec3, rest ...Vec3) Vec3 {
	res := v1
	for _, v := range rest {
		res[0] -= v[0]
		res[1] -= v[1]
		res[2] -= v[2]
	}
	return res
}

func (v1 Vec3) To(v2 Vec3) Vec3 {
	return Subtract(v2, v1)
}

func Add(v1 Vec3, rest ...Vec3) Vec3 {
	res := v1
	for _, v := range rest {
		res[0] += v[0]
		res[1] += v[1]
		res[2] += v[2]
	}
	return res
}

func Mul(v1, v2 Vec3) Vec3 {
	return Vec3{v1[0] * v2[0], v1[1] * v2[1], v1[2] * v2[2]}
}

func (v Vec3) Scaled(f float64) Vec3 {
	v[0] *= f
	v[1] *= f
	v[2] *= f
	return v
}

func Scale(v *Vec3, f float64) {
	v[0] *= f
	v[1] *= f
	v[2] *= f
}

func (v Vec3) Normalized() Vec3 {
	return v.Scaled(1.0 / v.Len())
}

func Normalize(v *Vec3) {
	Scale(v, (1.0 / v.Len()))
}

func Negate(v *Vec3) {
	v[0] *= -1.0
	v[1] *= -1.0
	v[2] *= -1.0
}

func Reflect(v *Vec3, n Vec3) {
	Subtract(*v, n.Scaled(2*Dot(*v, n)))
}

func Refract(v *Vec3, n Vec3, etai_over_etat float64) {
	cos_theta := math.Min(Dot(v.Scaled(-1.0), n), 1.0)
	r_out_perp := Add(*v, n.Scaled(cos_theta)).Scaled(etai_over_etat)
	r_out_parallel := n.Scaled(-math.Sqrt(math.Abs(1.0 - r_out_perp.LenSquared())))
	*v = Add(r_out_perp, r_out_parallel)
}

func (v Vec3) LenSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vec3) NearZero() bool {
	s := 1e-8
	return (math.Abs(v[0]) < s) && (math.Abs(v[1]) < s) && (math.Abs(v[2]) < s)
}

func (v Vec3) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v[0], -v[1], -v[2]}
}

func (v Vec3) Reflected(n Vec3) Vec3 {
	return Subtract(v, n.Scaled(2*Dot(v, n)))
}

func (v Vec3) Refracted(n Vec3, etai_over_etat float64) Vec3 {
	cos_theta := math.Min(Dot(v.Scaled(-1.0), n), 1.0)
	r_out_perp := Add(v, n.Scaled(cos_theta)).Scaled(etai_over_etat)
	r_out_parallel := n.Scaled(-math.Sqrt(math.Abs(1.0 - r_out_perp.LenSquared())))
	return Add(r_out_perp, r_out_parallel)
}

func (v Vec3) String() string {
	return fmt.Sprintf("%f %f %f", v[0], v[1], v[2])
}

func Dot(v1, v2 Vec3) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func Cross(v1, v2 Vec3) Vec3 {
	var res Vec3
	res[0] = v1[1]*v2[2] - v1[2]*v2[1]
	res[1] = v1[2]*v2[0] - v1[0]*v2[2]
	res[2] = v1[0]*v2[1] - v1[1]*v2[0]

	return res
}

func Random() Vec3 {
	return Vec3{rand.Float64()*2.0-1.0, rand.Float64()*2.0-1.0, rand.Float64()*2.0-1.0}
}

func RandomRanged(min, max float64) Vec3 {
	return Vec3{
		util.RandomFloat64(min, max),
		util.RandomFloat64(min, max),
		util.RandomFloat64(min, max)}
}

func RandUnit() Vec3 {
	for {
		v := Random()
		if 1e-150 < v.LenSquared() && v.LenSquared() <= 1 {
			return v.Normalized()
		}
	}
}

func RandInUnitDisk() Vec3 {
	for {
		res := Vec3{
			rand.Float64(),
			rand.Float64(),
			0}
		if res.LenSquared() <= 1.0 {
			return res
		}
	}
}

func RandomOnHemisphere(n Vec3) Vec3 {
	unit := RandUnit()
	if Dot(unit, n) > 0.0 {
		return unit
	}

	return unit.Scaled(-1.0)
}
