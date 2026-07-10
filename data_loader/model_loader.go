package data_loader

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"grt/render_object"
	"grt/vec"
)

const (
	OBJ = ".obj"
)

// TODO: Add material support from .mtl files
func Load3DModel(filepath string) (render_object.TriangleMesh, error) {

	extension := path.Ext(filepath)

	if extension == "" {
		return nil, fmt.Errorf("Provided file doesn't have an extension!!")
	}

	switch extension {
	case OBJ:
		return loadOBJModel(filepath)
	default:
		return nil, fmt.Errorf("Extension %s not supported", extension)
	}
}

func loadOBJModel(filepath string) (render_object.TriangleMesh, error) {

	modelFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(modelFile)
	fileScanner.Split(bufio.ScanLines)

	mesh := render_object.TriangleMesh{}

	vertex_positions := make([]vec.Vec3, 0)
	vertex_normals  := make([]vec.Vec3, 0)
	tex_coords := make([]vec.Vec3, 0)
	
	OBJ_DATA_SEP := " "
	OBJ_VERT_SEP := "/"

	add_vertex_pos := func(line string) {
		parts := strings.Split(line, OBJ_DATA_SEP)
		x, _ := strconv.ParseFloat(parts[1], 64)
		y, _ := strconv.ParseFloat(parts[2], 64)
		z, _ := strconv.ParseFloat(parts[3], 64)
		if len(parts) > 4 {
			w, _ := strconv.ParseFloat(parts[4], 64)
			x = x/w
			y = y/w
			z = z/w
		}
		vertex_positions = append(vertex_positions, vec.Vec3{x, y, z})
	}

	add_vertex_norm := func(line string) {
		parts := strings.Split(line, OBJ_DATA_SEP)
		x, _ := strconv.ParseFloat(parts[1], 64)
		y, _ := strconv.ParseFloat(parts[2], 64)
		z, _ := strconv.ParseFloat(parts[3], 64)
		vertex_normals = append(vertex_normals, vec.Vec3{x, y, z})
	}

	
	add_triangle := func(line string) {
		parts := strings.Split(line, OBJ_DATA_SEP)
		num_vertices := 3
		vertices := make([]render_object.Vertex, num_vertices)

		for i := range num_vertices {
			vertex_data := strings.Split(parts[i+1], OBJ_VERT_SEP)
			posIdx, _ := strconv.Atoi(vertex_data[0])
			vertices[i].Position = vertex_positions[posIdx-1]
			if len(vertex_data) == 1 {
				return
			}

			texIdx, err := strconv.Atoi(vertex_data[1])
			if err == nil {
				vertices[i].TexCoord = tex_coords[texIdx-1]
			}

			if len(vertex_data) == 2 {
				return
			}

			normIdx, _ := strconv.Atoi(vertex_data[2])
			vertices[i].Normal = vertex_normals[normIdx-1]
		}

		mesh = append(mesh, *render_object.NewTriangle(vertices[0], vertices[1], vertices[2], nil))
	}

	const (
		vertex_pos_line = "v"
		vertex_norm_line = "vn"
		vertex_tex_line = "vt"
		face_line = "f"
	)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		line_start := strings.SplitN(line, " ", 2)[0]

		switch line_start {
		case vertex_pos_line:
			add_vertex_pos(line)
		case vertex_norm_line:
			add_vertex_norm(line)
		// case vertex_tex_line:
		// 	add_vertex_tex(line)
		case face_line:
			add_triangle(line)
		default:
			
		}
	}

	fmt.Printf("Loaded %d positions, %d texture UVs, %d normals and produced %d triangles\n", len(vertex_positions), 0, len(vertex_normals), len(mesh))

	return mesh, nil
}
