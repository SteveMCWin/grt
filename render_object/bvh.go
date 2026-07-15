package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
	"math"
)

type BVHTree struct {
	Nodes        []BVHNode // NOTE: should actually be an interface, but I'll figure out what the interface should be later on
	Triangles    []Triangle
	TriangleIdxs []int
	RootNodeIdx  int
	NodesUsed    int
}

func NewBVHTree(model Model) *BVHTree {
	tree := BVHTree{
		Triangles:    model.Mesh,
		TriangleIdxs: make([]int, len(model.Mesh)),
		Nodes:        make([]BVHNode, 2*len(model.Mesh)-1),
		RootNodeIdx:  0,
		NodesUsed:    1,
	}

	for i := range tree.TriangleIdxs {
		tree.TriangleIdxs[i] = int(i)
	}

	tree.Nodes[tree.RootNodeIdx].NumTris = len(model.Mesh)
	tree.Nodes[tree.RootNodeIdx].UpdateBounds(tree.Triangles, tree.TriangleIdxs)
	tree.Nodes[tree.RootNodeIdx].Subdivide(tree.Triangles, tree.TriangleIdxs, tree.Nodes, &tree.NodesUsed)

	return &tree
}

type BVHNode struct {
	AABBMin    vec.Vec3
	AABBMax    vec.Vec3
	LeftChild  int
	FirstPrim  int
	NumTris    int
}

func (node *BVHNode) IsLeaf() bool {
	return node.NumTris > 0
}

func (node *BVHNode) UpdateBounds(allTris []Triangle, triIdxs []int) {
	node.AABBMin = vec.LargestVec3()
	node.AABBMax = vec.SmallestVec3()
	firstTri := node.FirstPrim
	for  i := range node.NumTris {
		tri := allTris[triIdxs[firstTri+i]]
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.A.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.B.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.C.Position)

		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.A.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.B.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.C.Position)
	}
}

func (node *BVHNode) Subdivide(tris []Triangle, triIdxs []int, allNodes []BVHNode, nodesUsed *int) {
	if node.NumTris <= 2 {
		return
	}

	// determine spllit axis
	diff := vec.Subtract(node.AABBMax, node.AABBMin)
	axis := 0
	if diff.X() < diff.Y() {
		axis = 1
	}
	if diff[axis] < diff.Z() {
		axis = 2
	}

	splitPos := node.AABBMin[axis] + 0.5*diff[axis]

	i := node.FirstPrim
	j := i + node.NumTris-1

	for i <= j {
		if tris[triIdxs[i]].center[axis] < splitPos {
			i++
		} else {
			triIdxs[i], triIdxs[j] = triIdxs[j], triIdxs[i]
			j--
		}
	}

	leftCount := i - node.FirstPrim
	if leftCount <= 0 || leftCount >= node.NumTris {
		return
	}

	leftChild := *nodesUsed
	rightChild := *nodesUsed+1
	node.LeftChild = leftChild
	*nodesUsed += 2

	allNodes[leftChild].FirstPrim = node.FirstPrim
	allNodes[leftChild].NumTris = leftCount
	allNodes[rightChild].FirstPrim = i
	allNodes[rightChild].NumTris = node.NumTris - leftCount
	node.NumTris = 0

	allNodes[leftChild].UpdateBounds(tris, triIdxs)
	allNodes[rightChild].UpdateBounds(tris, triIdxs)

	allNodes[leftChild].Subdivide(tris, triIdxs, allNodes, nodesUsed)
	allNodes[rightChild].Subdivide(tris, triIdxs, allNodes, nodesUsed)
}

func (tree *BVHTree) Intersect(r *ray.Ray, interval util.Interval, hit_record *HitRecord, nodeIdx int) bool {

	node := tree.Nodes[nodeIdx]
	if !node.Hit(r, interval, hit_record) {
		return false
	}

	if node.IsLeaf() {

		hit_anything := false
		closest_so_far := interval.Max

		for i := range node.NumTris {
			if tree.Triangles[tree.TriangleIdxs[node.FirstPrim+i]].Hit(r, util.Interval{Min: interval.Min, Max: closest_so_far}, hit_record) {
				hit_anything = true
				closest_so_far = hit_record.T
			}
		}
		return hit_anything
	} 

	// didn't hit anything so we don't need to update the interval here
	hitLeft := tree.Intersect(r, interval, hit_record, node.LeftChild)

	rightInterval := interval
	if hitLeft {
		rightInterval.Max = hit_record.T
	}

	hitRight := tree.Intersect(r, rightInterval, hit_record, node.LeftChild+1)

	return hitLeft || hitRight
}

func (tree *BVHTree) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	nodeIdx := tree.RootNodeIdx

	return tree.Intersect(r, interval, hit_record, nodeIdx)
}

func (node *BVHNode) Hit(r *ray.Ray, interval util.Interval, hit_record *HitRecord) bool {
	tx1 := (node.AABBMin.X()-r.Origin().X()) / r.Direction().X()
	tx2 := (node.AABBMax.X()-r.Origin().X()) / r.Direction().X()
	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (node.AABBMin.Y()-r.Origin().Y()) / r.Direction().Y()
	ty2 := (node.AABBMax.Y()-r.Origin().Y()) / r.Direction().Y()
	tmin = math.Max(tmin, min(ty1, ty2))
	tmax = math.Min(tmax, max(ty1, ty2))

	tz1 := (node.AABBMin.Z()-r.Origin().Z()) / r.Direction().Z()
	tz2 := (node.AABBMax.Z()-r.Origin().Z()) / r.Direction().Z()
	tmin = math.Max(tmin, min(tz1, tz2))
	tmax = math.Min(tmax, max(tz1, tz2))

	return tmax >= tmin && tmin < interval.Max && tmax > 0
}
