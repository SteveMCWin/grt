package main

import (
	"grt/camera"
	"grt/data_loader"
	"grt/render_object"
	"grt/vec"
	"log"
	"os"
	"time"
)


func main() {
	aspect_ratio := 16.0/9.0
	img_width := 800
	vfov := 40.0
	samples_per_pixel := 20
	max_depth := 10
	camera_pos := vec.Vec3{0.0, 0.9, 1.2}
	camera_look_at := vec.Vec3{0.0, 0.5, -1.0}
	camera_up := vec.Vec3{0.0, 1.0, 0.0}
	defocus_angle := 0.0
	focus_dist := 1.4

	c := camera.NewCam(aspect_ratio, vfov, defocus_angle, focus_dist, img_width, samples_per_pixel, max_depth, camera_pos, camera_look_at, camera_up)

	groundTexture := &render_object.FlatColorProvider{ Color: vec.Color{0.8, 0.8, 0.0} }
	teapotTexture := &render_object.FlatColorProvider{ Color: vec.Color{ 0.933, 0.945, 0.940 } }

	material_ground := &render_object.Lambertian{Texture: groundTexture}
	teapotMaterial := &render_object.Lambertian{Texture: teapotTexture}

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ 0.0, -100.5, -1.0}, Radius: 100.0, Material: material_ground })

	teapotMesh, err := data_loader.Load3DModel("models/teapot.obj")
	if err != nil {
		log.Println("Error loading file:", err)
	}

	teapotPosition := vec.Vec3{0.0, 0.2, -1.0}
	teapotRotation := vec.Vec3{0.0, 3.0,  0.0}
	teapotScale    := vec.Vec3{1.0, 1.0, 1.0}

	teapotModel := render_object.NewModel(teapotMesh, teapotMaterial, teapotPosition, teapotRotation, teapotScale)

	world.Add(&teapotModel)

	outputFile, err := os.Create("img.ppm")
	if err != nil {
		panic(err)
	}

	start_time := time.Now()
	c.Render(&world, outputFile)
	log.Println("Time to render:", time.Since(start_time))
}
