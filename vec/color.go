package vec

import (
	"fmt"
	"math"
	"grt/util"
)

const (
	GAMMA_POWER = 2.4
	INV_GAMMA_POWER = 1.0/2.4
	LINEAR_THRESHOLD = 0.0031308
	GAMMA_THRESHOLD = 0.04045
)

type Color Vec3

func LinearToGamma(col Color) Color {
	res := Color{}
	for i := range col {
		if col[i] < LINEAR_THRESHOLD {
			res[i] = col[i] * 12.92
		} else {
			res[i] = 1.055 * math.Pow(col[i], INV_GAMMA_POWER) - 0.055
		}
	}
	
	return res
}

func GammaToLinear(col Color) Color {
	res := Color{}
	for i := range col {
		if col[i] < GAMMA_THRESHOLD {
			res[i] = col[i] / 12.92
		} else {
			c := (col[i] + 0.055) / 1.055
			res[i] = math.Pow(c, GAMMA_POWER)
		}
	}
	
	return res
}

func (c Color) String() string {
	gamma_corrected := LinearToGamma(c)
	color_intensity := util.Interval{Min: 0.0, Max: 0.999}
	r := int(256 * color_intensity.Clamp(gamma_corrected[0]))
	g := int(256 * color_intensity.Clamp(gamma_corrected[1]))
	b := int(256 * color_intensity.Clamp(gamma_corrected[2]))

	return fmt.Sprintf("%d %d %d\n", r, g, b)
}

func (c Color) Components() (r, g, b byte) {
	gamma_corrected := LinearToGamma(c)
	color_intensity := util.Interval{Min: 0.0, Max: 0.999}
	r = byte(256 * color_intensity.Clamp(gamma_corrected[0]))
	g = byte(256 * color_intensity.Clamp(gamma_corrected[1]))
	b = byte(256 * color_intensity.Clamp(gamma_corrected[2]))
	return
}

func (c1 Color) Equals(c2 Color) bool {
	return c1[0] == c2[0] && c1[1] == c2[1] && c1[2] == c2[2]
}

// Normalizes the color only if it's intensity is greater than 1
func (c Color) Normalized() Color {
	vc := Vec3(c)
	if vc.LenSquared() > 1.0 {
		vc = vc.Normalized()
	}
	return Color(vc)
}
