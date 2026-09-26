package data_loader

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"grt/render_object"
	"grt/vec"
)

const (
	OBJ = ".obj"
	OBJ_DATA_SEP = " "
	OBJ_VERT_SEP = "/"
)

func Load3DModel(filepath string) (render_object.TriangleMesh, render_object.Material, error) {

	extension := path.Ext(filepath)

	if extension == "" {
		return nil, nil, fmt.Errorf("Provided file doesn't have an extension!!")
	}

	switch extension {
	case OBJ:
		return loadOBJModel(filepath)
	default:
		return nil, nil, fmt.Errorf("Extension %s not supported", extension)
	}
}

func loadOBJModel(filepath string) (render_object.TriangleMesh, render_object.Material, error) {
	modelFile, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}

	fileScanner := bufio.NewScanner(modelFile)
	fileScanner.Split(bufio.ScanLines)

	mesh := render_object.TriangleMesh{}

	vertex_positions := make([]vec.Vec3, 0)
	vertex_normals  := make([]vec.Vec3, 0)
	tex_coords := make([]vec.Vec3, 0)
	

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

	add_vertex_tex := func(line string) {
		parts := strings.Split(line, OBJ_DATA_SEP)
		u, _ := strconv.ParseFloat(parts[1], 64)
		v, _ := strconv.ParseFloat(parts[2], 64)
		tex_coords = append(tex_coords, vec.Vec3{u, v, 0.0})
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
		material_line = "mtllib"
	)

	var mat render_object.Material

	for fileScanner.Scan() {
		line := fileScanner.Text()
		line_start := strings.SplitN(line, " ", 2)[0]

		switch line_start {
		case material_line:
			mat = LoadTextureMap(line)
		case vertex_pos_line:
			add_vertex_pos(line)
		case vertex_norm_line:
			add_vertex_norm(line)
		case vertex_tex_line:
			add_vertex_tex(line)
		case face_line:
			add_triangle(line)
		default:
			
		}
	}

	fmt.Printf("Loaded %d positions, %d texture UVs, %d normals and produced %d triangles\n", len(vertex_positions), len(tex_coords), len(vertex_normals), len(mesh))

	return mesh, mat, nil
}

func LoadTextureMap(line string) (mat render_object.Material) {
	parts := strings.Split(line, OBJ_DATA_SEP)
	if len(parts) < 2 {
		log.Println("ERROR: when splitting material line, had less than two elements")
		mat = &render_object.Lambertian{ Texture: &render_object.DebugColorProvider{ Type: render_object.VertexColor }}
		return
	}
	mtl_file_name := parts[1]

	mtlFile, err := os.Open(mtl_file_name)
	if err != nil {
		log.Println("ERROR: couldn't open mtl file from mtl_file_name:", mtl_file_name)
		mat = &render_object.Lambertian{ Texture: &render_object.DebugColorProvider{ Type: render_object.VertexColor }}
		return
	}

	mtlFileScanner := bufio.NewScanner(mtlFile)
	var mapKdLine string

	for mtlFileScanner.Scan() {
		line := strings.TrimSpace(mtlFileScanner.Text())
		if strings.HasPrefix(line, "map_Kd") {
			mapKdLine = line
			break // Found the first match, stop scanning
		}
	}

	if err := mtlFileScanner.Err(); err != nil {
		log.Println("ERROR: while scanning mtl file:", err)
		mat = &render_object.Lambertian{ Texture: &render_object.DebugColorProvider{ Type: render_object.VertexColor }}
		return
	}

	albedoTexture := strings.Split(mapKdLine, OBJ_DATA_SEP)[1]
	imageReader, err := os.Open(albedoTexture)
	if err != nil {
		log.Println("ERROR: while opening image: ", albedoTexture, " ERR_MSG: ", err)
		mat = &render_object.Lambertian{ Texture: &render_object.DebugColorProvider{ Type: render_object.VertexColor }}
		return
	}
	defer imageReader.Close()

	tex, format, err := image.Decode(imageReader)
	if err != nil {
		log.Println("ERROR: while decoding image: ", err)
		mat = &render_object.Lambertian{ Texture: &render_object.DebugColorProvider{ Type: render_object.VertexColor }}
		return
	}
	log.Println("Decoded texture map with format ", format)

	mat = &render_object.Lambertian{ Texture: &render_object.TextureColorProvider{ Texture: tex }}

	return
}
