package render_object

import (
	"grt/vec"
)

type BVHTree struct {
	Nodes        []BVHNode // NOTE: should actually be an interface, but I'll figure out what the interface should be later on
	Triangles    []Triangle
	TriangleIdxs []int
	RootNodeIdx  int
	NodesUsed    int
	MaxDepth     int
}

func NewBVHTree(model Model, maxDepth int) *BVHTree {
	tree := BVHTree{
		Triangles:    model.Mesh,
		TriangleIdxs: make([]int, len(model.Mesh)),
		Nodes:        make([]BVHNode, 2*len(model.Mesh)-1),
		RootNodeIdx:  0,
		NodesUsed:    1,
		MaxDepth:     maxDepth,
	}

	for i := range tree.TriangleIdxs {
		tree.TriangleIdxs[i] = int(i)
	}

	tree.Nodes[tree.RootNodeIdx].UpdateBounds(tree.Triangles)
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

func (node *BVHNode) UpdateBounds(allTris []Triangle) {
	node.AABBMin = vec.LargestVec3()
	node.AABBMax = vec.SmallestVec3()
	firstTri := node.FirstPrim
	for  i := range node.NumTris {
		node.AABBMin = vec.MinCoords(node.AABBMin, allTris[firstTri+i].A.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, allTris[firstTri+i].B.Position)
		node.AABBMin = vec.MinCoords(node.AABBMin, allTris[firstTri+i].C.Position)

		node.AABBMax = vec.MaxCoords(node.AABBMax, allTris[firstTri+i].A.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, allTris[firstTri+i].B.Position)
		node.AABBMax = vec.MaxCoords(node.AABBMax, allTris[firstTri+i].C.Position)
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

	allNodes[leftChild].UpdateBounds(tris)
	allNodes[rightChild].UpdateBounds(tris)

	allNodes[leftChild].Subdivide(tris, triIdxs, allNodes, nodesUsed)
	allNodes[rightChild].Subdivide(tris, triIdxs, allNodes, nodesUsed)
}
