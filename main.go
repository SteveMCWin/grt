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
	camDesc := camera.CamDescriptor {
		AspectRatio: 16.0/9.0,
		ImageWidth: 800,
		Vfov: 40,
		SamplesPerPixel: 200,
		MaxDepth: 20,
		Position: vec.Vec3{0.0, 0.9, 1.2},
		LookAt: vec.Vec3{0.0, 0.5, -1.0},
		Up: vec.Vec3{0.0, 1.0, 0.0},
		DefocusAngle: 0.0,
		FocusDist: 1.4,
		Background: vec.Color{ 0.01, 0.01, 0.01 },
	}

	c := camera.NewCam(camDesc)

	// teapotTexture := &render_object.FlatColorProvider{ Color: vec.Color{ 0.933, 0.945, 0.940 } }
	lightTexture1 := &render_object.FlatColorProvider{ Color: vec.Color{ 8.0, 4.0, 3.0 } }
	lightTexture2 := &render_object.FlatColorProvider{ Color: vec.Color{ 3.0, 4.0, 8.0 } }
	lightTexture3 := &render_object.FlatColorProvider{ Color: vec.Color{ 4.0, 4.0, 4.0 } }
	// mistTexture   := &render_object.FlatColorProvider{ Color: vec.Color{ 0.8, 0.8, 0.8 } }

	// teapotMaterial := &render_object.Lambertian{Texture: teapotTexture}
	lightMaterial1 := &render_object.Emitter{Texture: lightTexture1}
	lightMaterial2 := &render_object.Emitter{Texture: lightTexture2}
	lightMaterial3 := &render_object.Emitter{Texture: lightTexture3}
	// mistMaterial   := &render_object.Lambertian{Texture: mistTexture}

	// mistBoundary := render_object.Sphere{ Center: vec.Vec3{ 0.0, 0.0, 0.0 }, Radius: 50.0, Material: mistMaterial }

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ -1.0,    1.2, -0.7 }, Radius:   0.5, Material: lightMaterial1 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ +1.0,    1.2, -0.8 }, Radius:  0.35, Material: lightMaterial2 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{  0.0,   -3.0, -0.3 }, Radius:   0.2, Material: lightMaterial3 })
	// world.Add(&render_object.Mist{ Boundary: &mistBoundary, Material: mistMaterial, NegInvDensity: -100.0 })

	// teapotMesh, err := data_loader.Load3DModel("models/teapot/teapot.obj")
	// if err != nil {
	// 	log.Println("Error loading teapot model:", err)
	// }

	// teapotPosition := vec.Vec3{0.0, 0.3, -1.0}
	// teapotRotation := vec.Vec3{0.0, 3.0,  0.0}
	// teapotScale    := vec.Vec3{0.5, 0.5, 0.5}
	
	boatMesh, boatMaterial, err := data_loader.Load3DModel("models/boat/boat.obj")
	if err != nil {
		log.Println("Error loading boat model:", err)
	}

	boatPosition := vec.Vec3{0.0, 0.3, -1.0}
	boatRotation := vec.Vec3{0.0, 3.0,  0.0}
	boatScale    := vec.Vec3{0.5, 0.5, 0.5}

	log.Println("Making BVH...")
	// teapotModel := render_object.NewModel(teapotMesh, teapotMaterial, teapotPosition, teapotRotation, teapotScale)
	boatModel := render_object.NewModel(boatMesh, boatMaterial, boatPosition, boatRotation, boatScale)
	log.Println("Made BVH successfully!")

	world.Add(&boatModel)

	outputFile, err := os.Create("img.ppm")
	if err != nil {
		panic(err)
	}

	start_time := time.Now()
	c.Render(&world, outputFile)
	log.Println("Time to render:", time.Since(start_time))
}
