package util

import (
	"math"
	"math/rand/v2"
)

var NegInf = math.Inf(-1)
var PosInf = math.Inf(1)
var Small  = 1e-5

type Interval struct {
	Min float64
	Max float64
}

var EmtpyInterval = Interval{Min: PosInf, Max: NegInf}
var EverythingInterval = Interval{Min: NegInf, Max: PosInf}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Contains(f float64) bool {
	return i.Min <= f && f <= i.Max
}

func (i Interval) Surrounds(f float64) bool {
	return i.Min < f && f < i.Max
}

func (i Interval) Clamp(f float64) float64 {
	return max(i.Min, min(i.Max, f))
}

func (i Interval) Expand(delta float64) Interval {
	padding := delta/2.0
	return Interval{Min: i.Min-padding, Max: i.Max+padding}
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func RandomFloat64(min, max float64) float64 {
	return min + (max-min)*rand.Float64() 
}
