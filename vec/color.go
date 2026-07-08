package vec

import (
	"fmt"
	"math"
	"grt/util"
)

type Color Vec3

func LinearToGamma(col Color) Color {
	res := Color{}
	for i := range col {
		res[i] = math.Sqrt(col[i])
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
