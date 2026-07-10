package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
)

type TriangleMesh []Triangle

type Model struct {
	Mesh     TriangleMesh
	Material Material
	Position vec.Vec3
	Rotation vec.Vec3
	Scale    vec.Vec3

	inverseMatrix vec.Mat4
	modelMatrix   vec.Mat4
}

func NewModel(mesh TriangleMesh, material Material, position, rotation, scale vec.Vec3) Model {
	return Model{
		Mesh:          mesh,
		Material:      material,
		Position:      position,
		Rotation:      rotation,
		Scale:         scale,
		modelMatrix:   vec.ModelMatrix(position, rotation, scale),
		inverseMatrix: vec.InverseModelMatrix(position, rotation, scale),
	}
}

func (m *Model) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	var rec HitRecord
	hit_anything := false
	closest_so_far := interval.Max

	for _, triangle := range m.Mesh {
		// transform the ray into model space
		transformed_ray := ray.NewRay(m.inverseMatrix.MultiplyPoint(r.Origin()), m.inverseMatrix.MultiplyDir(r.Direction()))
		if triangle.Hit(&transformed_ray, util.Interval{Min: interval.Min, Max: closest_so_far}, &rec) {
			hit_anything = true
			closest_so_far = rec.T
			*hit_record = rec
			hit_record.Material = m.Material

			hit_record.Point = m.modelMatrix.MultiplyPoint(hit_record.Point)
			worldNormal := m.inverseMatrix.Transposed().MultiplyDir(hit_record.Normal)
			hit_record.SetFaceNormal(r, worldNormal.Normalized())
		}
	}

	return hit_anything
}
