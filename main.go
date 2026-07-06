package main

import (
	"log"
	"os"
	"time"
	"grt/camera"
	"grt/render_object"
	"grt/vec"
)


func main() {

	aspect_ratio := 16.0/9.0
	img_width := 600
	vfov := 40.0
	samples_per_pixel := 25
	max_depth := 30
	camera_pos := vec.Vec3{0.0, 0.9, 1.2}
	camera_look_at := vec.Vec3{0.0, 0.5, -1.0}
	camera_up := vec.Vec3{0.0, 1.0, 0.0}
	defocus_angle := 0.0
	focus_dist := 3.4

	c := camera.NewCam(aspect_ratio, vfov, defocus_angle, focus_dist, img_width, samples_per_pixel, max_depth, camera_pos, camera_look_at, camera_up)

	material_ground := render_object.Lambertian{Albedo: vec.Color{0.8, 0.8, 0.0}}
	material_center := render_object.Lambertian{Albedo: vec.Color{0.1, 0.2, 0.5}}

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ 0.0, -100.5, -1.0}, Radius: 100.0, Material: &material_ground })

	point_a := render_object.Vertex{ Position: vec.Vec3{-0.5, 0.2, -1.2}, Color: vec.Color{1.0, 0.0, 0.0} }
	point_b := render_object.Vertex{ Position: vec.Vec3{ 0.5, 0.2, -1.2}, Color: vec.Color{0.0, 1.0, 0.0} }
	point_c := render_object.Vertex{ Position: vec.Vec3{-0.5, 1.2, -1.2}, Color: vec.Color{0.0, 0.0, 1.0} }
	point_d := render_object.Vertex{ Position: vec.Vec3{ 0.5, 1.2, -1.2}, Color: vec.Color{1.0, 1.0, 0.0} }

	t1 := render_object.NewTriangle(point_a, point_b, point_c, &material_center)
	t2 := render_object.NewTriangle(point_c, point_b, point_d, &material_center)

	world.Add(t1)
	world.Add(t2)

	f, err := os.Create("img.ppm")
	if err != nil {
		panic(err)
	}

	start_time := time.Now()
	c.Render(&world, f)
	log.Println("Time to render:", time.Since(start_time))
}
