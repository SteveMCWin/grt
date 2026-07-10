package vec

import (
	"math"
)

const DIM = 4

// first coord is the column, the second is the row
// e.g. m[3][1] would be
// ⚬ ⚬ ⚬ ⚬
// ⚬ ⚬ ⚬ ●
// ⚬ ⚬ ⚬ ⚬
// ⚬ ⚬ ⚬ ⚬
type Mat4 [DIM][DIM]float64

func Identity() Mat4 {
	res := Mat4{}
	for i := range DIM {
		res[i][i] = 1
	}
	return res
}

func (m Mat4) IsZero() bool {

	for c := range DIM {
		for r := range DIM {
			if m[c][r] != 0 {
				return false
			}
		}
	}
	return true
}

func (m *Mat4) Translate(v Vec3) {
	for i := range len(v) {
		m[3][i] += v[i]
	}
}

func (m *Mat4) Scale(v Vec3) {
	for i := range len(v) {
		m[i][i] *= v[i]
	}
}

func Rotate(axis Vec3, radians float64) Mat4 {
	cos := math.Cos(radians)
	sin := math.Sin(radians)

	var elems [3][3]float64
	elems[0][0] = cos + axis.X()*axis.X()*(1-cos)
	elems[0][1] = axis.Y()*axis.X()*(1-cos) + axis.Z()*sin
	elems[0][2] = axis.Z()*axis.X()*(1-cos) - axis.Y()*sin
	elems[1][0] = axis.X()*axis.Y()*(1-cos) - axis.Z()*sin
	elems[1][1] = cos + axis.Y()*axis.Y()*(1-cos)
	elems[1][2] = axis.Z()*axis.Y()*(1-cos) + axis.X()*sin
	elems[2][0] = axis.X()*axis.Z()*(1-cos) + axis.Y()*sin
	elems[2][1] = axis.Y()*axis.Z()*(1-cos) - axis.X()*sin
	elems[2][2] = cos + axis.Z()*axis.Z()*(1-cos)

	res := Identity()

	for c := range len(elems) {
		for r := range len(elems[0]) {
			res[c][r] = elems[c][r]
		}
	}

	return res

}

func Translation(v Vec3) Mat4 {
	res := Identity()

	for i := range len(v) {
		res[3][i] += v[i]
	}

	return res
}

func Scaling(v Vec3) Mat4 {
	res := Identity()

	for i := range len(v) {
		res[i][i] *= v[i]
	}

	return res
}

func Multiply(multiplicators ...Mat4) Mat4 {
	if len(multiplicators) < 2 {
		return Mat4{}
	}

	res := multiplicators[0]
	multiplicators = multiplicators[1:]

	for _, mul := range multiplicators {
		next := Mat4{}
		for c := range DIM {
			for r := range DIM {
				sum := 0.0
				for k := range DIM {
					sum += res[k][r] * mul[c][k]
				}
				next[c][r] = sum
			}
		}
		res = next
	}

	return res
}

func (m *Mat4) Scaled(amount float64) {
	for c := range DIM {
		for r := range DIM {
			m[c][r] *= amount
		}
	}
}

func (m *Mat4) Diag() Vec3 {
	return Vec3{m[0][0], m[1][1], m[2][2]}
}

func Vec3ToTranslate(v Vec3) Mat4 {
	res := Identity()
	for i := range 3 {
		res[3][i] = v[i]
	}
	return res
}

func Vec3ToRotation(v Vec3) Mat4 {
	axis := v.Normalized()
	degreesRadian := v.Len() // NOTE: This may need to be subtracted by 1.0 or somethign

	return Rotate(axis, degreesRadian)
}

func Vec3ToScale(v Vec3) Mat4 {
	res := Mat4{}
	for i := range 3 {
		res[i][i] = v[i]
	}
	return res
}

func ModelMatrix(position, rotation, scale Vec3) Mat4 {
	matScale := Vec3ToScale(scale)
	matRotation := Vec3ToRotation(rotation)
	matTranslation := Vec3ToTranslate(position)

	return Multiply(matScale, matRotation, matTranslation)
}

func InverseModelMatrix(position, rotation, scale Vec3) Mat4 {
	// for inverting the position, scale and rotation matrices, we don't need a generalized invertion matrix
	invScale := Vec3ToScale(scale.Reciprocated())
	invRotation := Vec3ToRotation(rotation).Transposed()
	invTranslation := Vec3ToTranslate(position.Neg())

	return Multiply(invScale, invRotation, invTranslation)
}

func (m Mat4) Transposed() Mat4 {
	res := Mat4{}
	for c := range DIM {
		for r := range DIM {
			res[c][r] = m[r][c]
		}
	}
	return res
}

func Transpose(m *Mat4) {
	for c := range DIM {
		for r := range c {
			m[c][r], m[r][c] = m[r][c], m[c][r]
		}
	}
}

func (m Mat4) MultiplyPoint(v Vec3) Vec3 {
	res := [4]float64{}
	v4 := [4]float64{v[0], v[1], v[2], 1.0}

	for r := range DIM {
		sum := 0.0
		for c := range DIM {
			sum += m[c][r] * v4[c]
		}
		res[r] = sum
	}

	return Vec3{res[0], res[1], res[2]}
}

func (m Mat4) MultiplyDir(v Vec3) Vec3 {
	res := [4]float64{}
	v4 := [4]float64{v[0], v[1], v[2], 0.0}

	for r := range DIM {
		sum := 0.0
		for c := range DIM {
			sum += m[c][r] * v4[c]
		}
		res[r] = sum
	}

	return Vec3{res[0], res[1], res[2]}
}
