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
	// material_left := render_object.Dielectric{RefractionIndex: 1.5}
	// material_bubble := render_object.Dielectric{RefractionIndex: 1.0/1.5}
	// material_right := render_object.Metal{Albedo: vec.Color{0.8, 0.6, 0.2}, Fuzz:0.02}
	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ 0.0, -100.5, -1.0}, Radius: 100.0, Material: &material_ground })
	// world.Add(&render_object.Sphere{ Center: vec.Vec3{ 0.0,    0.0, -1.2}, Radius:   0.5, Material: &material_center })
	// world.Add(&render_object.Sphere{ Center: vec.Vec3{-1.0,    0.0, -1.0}, Radius:   0.5, Material: &material_left })
	// world.Add(&render_object.Sphere{ Center: vec.Vec3{-1.0,    0.0, -1.0}, Radius:   0.4, Material: &material_bubble })
	// world.Add(&render_object.Sphere{ Center: vec.Vec3{ 1.0,    0.0, -1.0}, Radius:   0.5, Material: &material_right })
	point_a := vec.Vec3{-0.5, 0.2, -1.2}
	point_b := vec.Vec3{ 0.5, 0.2, -1.2}
	point_c := vec.Vec3{-0.5, 1.2, -1.2}
	point_d := vec.Vec3{ 0.5, 1.2, -1.2}

	world.Add(&render_object.Triangle{A: point_a, B: point_b, C: point_c, Material: &material_center})
	world.Add(&render_object.Triangle{A: point_c, B: point_b, C: point_d, Material: &material_center})

	f, err := os.Create("img.ppm")
	if err != nil {
		panic(err)
	}

	start_time := time.Now()
	c.Render(&world, f)
	log.Println("Time to render:", time.Since(start_time))
}
