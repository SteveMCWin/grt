package render_object

import (
	"grt/ray"
	"grt/util"
	"grt/vec"
	"math"
)

type BVHTree struct {
	Nodes        []BVHNode
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

type AABB struct {
	AABBMin vec.Vec3
	AABBMax vec.Vec3
}

func (aabb *AABB) GrowToInclude(p vec.Vec3) {
	aabb.AABBMin = vec.MinCoords(aabb.AABBMin, p)
	aabb.AABBMax = vec.MaxCoords(aabb.AABBMax, p)
}

func (aabb *AABB) Area() float64 {
	e := vec.Subtract(aabb.AABBMax, aabb.AABBMin)
	return e.X()*e.Y() + e.Y()*e.Z() + e.X()*e.Z()
}

type BVHNode struct {
	AABB
	LeftChild int
	FirstPrim int
	NumTris   int
}

func (node *BVHNode) IsLeaf() bool {
	return node.NumTris > 0
}

func (node *BVHNode) UpdateBounds(allTris []Triangle, triIdxs []int) {
	node.AABBMin = vec.LargestVec3()
	node.AABBMax = vec.SmallestVec3()
	firstTri := node.FirstPrim
	for i := range node.NumTris {
		tri := allTris[triIdxs[firstTri+i]]
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.A.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.B.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, tri.C.Position)

		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.A.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.B.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, tri.C.Position)
	}
}

func evaluateSAH(node *BVHNode, axis int, pos float64, tris []Triangle, triIdxs []int) float64 {
	leftBox := AABB{
		AABBMin: vec.LargestVec3(),
		AABBMax: vec.SmallestVec3(),
	}
	rightBox := leftBox
	leftCount := 0
	rightCount := 0

	for i := range node.NumTris {
		tri := tris[triIdxs[node.FirstPrim+i]]
		if tri.center[axis] < pos {
			leftCount++
			leftBox.GrowToInclude(tri.A.Position)
			leftBox.GrowToInclude(tri.B.Position)
			leftBox.GrowToInclude(tri.C.Position)
		} else {
			rightCount++
			rightBox.GrowToInclude(tri.A.Position)
			rightBox.GrowToInclude(tri.B.Position)
			rightBox.GrowToInclude(tri.C.Position)
		}
	}

	cost := float64(leftCount)*leftBox.Area() + float64(rightCount)*rightBox.Area()
	if cost > 0 {
		return cost
	}
	return util.PosInf
}

func (node *BVHNode) cost() float64 {
	return float64(node.NumTris) * node.AABB.Area()
}

func (node *BVHNode) bestSplitPlane(tris []Triangle, triIdxs []int) (bestAxis int, bestPos, bestCost float64) {
	bestAxis = -1
	bestPos = 0.0
	bestCost = util.PosInf

	for axis := range 3 {
		for i := range node.NumTris {
			tri := tris[triIdxs[node.FirstPrim+i]]
			candidatePos := tri.center[axis]
			cost := evaluateSAH(node, axis, candidatePos, tris, triIdxs)
			if cost < bestCost {
				bestPos = candidatePos
				bestAxis = axis
				bestCost = cost
			}
		}
	}

	return
}

func (node *BVHNode) bestSplitPlane2(tris []Triangle, triIdxs []int) (bestAxis int, bestPos, bestCost float64) {
	bestAxis = -1
	bestPos = 0.0
	bestCost = util.PosInf

	for axis := range 3 {
		boundsMin := node.AABBMin[axis]
		boundsMax := node.AABBMax[axis]
		if boundsMin >= boundsMax {
			continue
		}
		subdivisions := 15
		scale := (boundsMax - boundsMin) / float64(subdivisions)
		for i := range subdivisions {
			candidatePos := boundsMin + scale * float64(i)
			cost := evaluateSAH(node, axis, candidatePos, tris, triIdxs)
			if cost < bestCost {
				bestPos = candidatePos
				bestAxis = axis
				bestCost = cost
			}
		}
	}

	return
}

func (node *BVHNode) Subdivide(tris []Triangle, triIdxs []int, allNodes []BVHNode, nodesUsed *int) {
	node.LeftChild = -1
	noSplitCost := node.cost()
	bestAxis, bestPos, splitCost := node.bestSplitPlane2(tris, triIdxs)

	if splitCost >= noSplitCost {
		return
	}

	axis := bestAxis
	splitPos := bestPos

	i := node.FirstPrim
	j := i + node.NumTris - 1

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
	rightChild := *nodesUsed + 1
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
	stack := make([]BVHNode, 0)
	hit_anything := false

	for {
		if node.IsLeaf() {
			if node.Intersect(r, interval) >= util.PosInf {
				if len(stack) <= 0 {
					return hit_anything
				}
				node = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				continue
			}
			closest_so_far := interval.Max

			for i := range node.NumTris {
				if tree.Triangles[tree.TriangleIdxs[node.FirstPrim+i]].Hit(r, util.Interval{Min: interval.Min, Max: closest_so_far}, hit_record) {
					hit_anything = true
					closest_so_far = hit_record.T
					interval.Max = hit_record.T
				}
			}

			if len(stack) <= 0 {
				return hit_anything
			}
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			continue
		}
		child1 := tree.Nodes[node.LeftChild]
		child2 := tree.Nodes[node.LeftChild+1]
		dist1 := child1.Intersect(r, interval)
		dist2 := child2.Intersect(r, interval)
		if dist1 > dist2 {
			dist1, dist2 = dist2, dist1
			child1, child2 = child2, child1
		}

		if dist1 == util.PosInf {
			if len(stack) <= 0 {
				return hit_anything
			}
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			node = child1
			if dist2 != util.PosInf {
				stack = append(stack, child2)
			}
		}
	}
}

func (node *BVHNode) Intersect(r *ray.Ray, interval util.Interval) float64 {
	invertedRayDir := r.InvertedDir()
	tx1 := (node.AABBMin.X() - r.Origin().X()) * invertedRayDir.X()
	tx2 := (node.AABBMax.X() - r.Origin().X()) * invertedRayDir.X()
	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (node.AABBMin.Y() - r.Origin().Y()) * invertedRayDir.Y()
	ty2 := (node.AABBMax.Y() - r.Origin().Y()) * invertedRayDir.Y()
	tmin = math.Max(tmin, min(ty1, ty2))
	tmax = math.Min(tmax, max(ty1, ty2))

	tz1 := (node.AABBMin.Z() - r.Origin().Z()) * invertedRayDir.Z()
	tz2 := (node.AABBMax.Z() - r.Origin().Z()) * invertedRayDir.Z()
	tmin = math.Max(tmin, min(tz1, tz2))
	tmax = math.Min(tmax, max(tz1, tz2))

	if tmax >= tmin && tmin < interval.Max {
		return tmin
	}
	return util.PosInf
}
