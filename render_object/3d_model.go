package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
	"log"
)

type TriangleMesh []Triangle

type Model struct {
	Mesh     TriangleMesh
	Material Material
	Position vec.Vec3
	Rotation vec.Vec3
	Scale    vec.Vec3

	BVHTree *BVHTree

	inverseMatrix vec.Mat4
	modelMatrix   vec.Mat4
}

func NewModel(mesh TriangleMesh, material Material, position, rotation, scale vec.Vec3) Model {

	m := Model{
		Mesh:          mesh,
		Material:      material,
		Position:      position,
		Rotation:      rotation,
		Scale:         scale,
		modelMatrix:   vec.ModelMatrix(position, rotation, scale),
		inverseMatrix: vec.InverseModelMatrix(position, rotation, scale),
	}

	m.BVHTree = NewBVHTree(m)
	log.Println("Teapot bvh tree has this many nodes:", len(m.BVHTree.Nodes))

	return m
}

func (m *Model) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	transformed_ray := ray.NewRay(m.inverseMatrix.MultiplyPoint(r.Origin()), m.inverseMatrix.MultiplyDir(r.Direction()))

	hit_anything := m.BVHTree.Intersect(&transformed_ray, interval, hit_record, m.BVHTree.RootNodeIdx)
	hit_record.Material = m.Material
	hit_record.Point = m.modelMatrix.MultiplyPoint(hit_record.Point)

	outwardObjNormal := hit_record.Normal
	if !hit_record.Front {
		outwardObjNormal = outwardObjNormal.Neg()
	}

	worldNormal := m.inverseMatrix.Transposed().MultiplyDir(outwardObjNormal)
	hit_record.SetFaceNormal(r, worldNormal.Normalized())

	return hit_anything
}
