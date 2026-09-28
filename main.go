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
		SamplesPerPixel: 2000,
		MaxDepth: 20,
		Position: vec.Vec3{0.0, 0.9, 1.2},
		LookAt: vec.Vec3{0.0, 0.5, -1.0},
		Up: vec.Vec3{0.0, 1.0, 0.0},
		DefocusAngle: 0.0,
		FocusDist: 1.4,
		Background: vec.Color{ 0.003, 0.003, 0.003 },
	}

	// world := teapotScene(&camDesc)
	// world := boatScene(&camDesc)
	world := cornellBoxScene(&camDesc)

	c := camera.NewCam(camDesc)

	outputFile, err := os.Create("img.ppm")
	if err != nil {
		panic(err)
	}

	start_time := time.Now()
	c.Render(&world, outputFile)
	log.Println("Time to render:", time.Since(start_time))
}

func teapotScene(camDesc *camera.CamDescriptor) render_object.HittableList {
	teapotTexture := &render_object.FlatColorProvider{ Color: vec.Color{ 0.933, 0.945, 0.940 } }
	lightTexture1 := &render_object.FlatColorProvider{ Color: vec.Color{ 8.0, 4.0, 3.0 } }
	lightTexture2 := &render_object.FlatColorProvider{ Color: vec.Color{ 3.0, 4.0, 8.0 } }
	lightTexture3 := &render_object.FlatColorProvider{ Color: vec.Color{ 4.0, 4.0, 4.0 } }

	teapotMaterial := &render_object.Lambertian{Texture: teapotTexture}
	lightMaterial1 := &render_object.Emitter{Texture: lightTexture1}
	lightMaterial2 := &render_object.Emitter{Texture: lightTexture2}
	lightMaterial3 := &render_object.Emitter{Texture: lightTexture3}

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ -1.0,    1.2, -0.7 }, Radius:   0.5, Material: lightMaterial1 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ +1.0,    1.2, -0.8 }, Radius:  0.35, Material: lightMaterial2 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{  0.0,   -5.0, -0.8 }, Radius:   0.2, Material: lightMaterial3 })

	teapotMesh, _, err := data_loader.Load3DModel("models/teapot/teapot.obj")
	if err != nil {
		log.Println("Error loading teapot model:", err)
	}

	teapotPosition := vec.Vec3{0.0, 0.3, -1.0}
	teapotRotation := vec.Vec3{0.0, 3.0,  0.0}
	teapotScale    := vec.Vec3{0.5, 0.5, 0.5}

	log.Println("Making BVH...")
	teapotModel := render_object.NewModel(teapotMesh, teapotMaterial, teapotPosition, teapotRotation, teapotScale)
	log.Println("Made BVH successfully!")

	world.Add(&teapotModel)

	return world
}

func boatScene(camDesc *camera.CamDescriptor) render_object.HittableList {
	lightTexture1 := &render_object.FlatColorProvider{ Color: vec.Color{ 8.0, 4.0, 3.0 } }
	lightTexture2 := &render_object.FlatColorProvider{ Color: vec.Color{ 3.0, 4.0, 8.0 } }
	lightTexture3 := &render_object.FlatColorProvider{ Color: vec.Color{ 4.0, 4.0, 4.0 } }

	lightMaterial1 := &render_object.Emitter{Texture: lightTexture1}
	lightMaterial2 := &render_object.Emitter{Texture: lightTexture2}
	lightMaterial3 := &render_object.Emitter{Texture: lightTexture3}

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ -1.0,    1.2, -0.7 }, Radius:   0.5, Material: lightMaterial1 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{ +1.0,    1.2, -0.8 }, Radius:  0.35, Material: lightMaterial2 })
	world.Add(&render_object.Sphere{ Center: vec.Vec3{  0.0,   -3.0, -0.3 }, Radius:   0.2, Material: lightMaterial3 })

	boatMesh, boatMaterial, err := data_loader.Load3DModel("models/boat/boat.obj")
	if err != nil {
		log.Println("Error loading boat model:", err)
	}

	boatPosition := vec.Vec3{0.0, 0.3, -1.0}
	boatRotation := vec.Vec3{0.0, 3.0,  0.0}
	boatScale    := vec.Vec3{0.5, 0.5, 0.5}

	log.Println("Making BVH...")
	boatModel := render_object.NewModel(boatMesh, boatMaterial, boatPosition, boatRotation, boatScale)
	log.Println("Made BVH successfully!")

	world.Add(&boatModel)

	return world
}

// Box spans x [-1, 1], y [0, 2], z [-2, 0] and is open towards +z where the camera is.
// Triangles are only hit from the side their normal (AB x AC) points to, so every quad's normal points into the box.
func cornellBoxScene(camDesc *camera.CamDescriptor) render_object.HittableList {
	camDesc.AspectRatio = 1.0
	camDesc.Position = vec.Vec3{0.0, 1.0, 3.4}
	camDesc.LookAt = vec.Vec3{0.0, 1.0, -1.0}
	camDesc.FocusDist = 4.4

	whiteTexture  := &render_object.FlatColorProvider{ Color: vec.Color{ 0.73, 0.73, 0.73 } }
	greenTexture  := &render_object.FlatColorProvider{ Color: vec.Color{ 0.1, 0.7, 0.1 } }
	redTexture    := &render_object.FlatColorProvider{ Color: vec.Color{ 0.7, 0.1, 0.1 } }
	lightTexture  := &render_object.FlatColorProvider{ Color: vec.Color{ 8.0, 8.0, 8.0 } }
	teapotTexture := &render_object.FlatColorProvider{ Color: vec.Color{ 0.933, 0.945, 0.940 } }
	mirrorTexture := &render_object.FlatColorProvider{ Color: vec.Color{ 0.95, 0.95, 0.95 } }

	whiteMaterial  := &render_object.Lambertian{Texture: whiteTexture}
	greenMaterial  := &render_object.Lambertian{Texture: greenTexture}
	redMaterial    := &render_object.Lambertian{Texture: redTexture}
	lightMaterial  := &render_object.Emitter{Texture: lightTexture}
	teapotMaterial := &render_object.Lambertian{Texture: teapotTexture}
	mirrorMaterial := &render_object.Metal{Texture: mirrorTexture, Fuzz: 0.0}

	// WORLD
	world := render_object.HittableList{ Shapes: make([]render_object.Hittable, 0) }

	// quad with corners p, p+u, p+u+v, p+v and normal u x v
	addQuad := func(p, u, v vec.Vec3, m render_object.Material) {
		a := render_object.Vertex{Position: p}
		b := render_object.Vertex{Position: vec.Add(p, u)}
		c := render_object.Vertex{Position: vec.Add(p, u, v)}
		d := render_object.Vertex{Position: vec.Add(p, v)}
		world.Add(render_object.NewTriangle(a, b, c, m))
		world.Add(render_object.NewTriangle(a, c, d, m))
	}

	addQuad(vec.Vec3{-1.0, 0.0,  0.0}, vec.Vec3{2.0, 0.0,  0.0}, vec.Vec3{0.0, 0.0, -2.0}, whiteMaterial) // floor
	addQuad(vec.Vec3{-1.0, 2.0,  0.0}, vec.Vec3{0.0, 0.0, -2.0}, vec.Vec3{2.0, 0.0,  0.0}, whiteMaterial) // ceiling
	addQuad(vec.Vec3{-1.0, 0.0, -2.0}, vec.Vec3{2.0, 0.0,  0.0}, vec.Vec3{0.0, 2.0,  0.0}, whiteMaterial) // back wall
	addQuad(vec.Vec3{-1.0, 0.0,  0.0}, vec.Vec3{0.0, 0.0, -2.0}, vec.Vec3{0.0, 2.0,  0.0}, greenMaterial) // left wall
	addQuad(vec.Vec3{ 1.0, 0.0,  0.0}, vec.Vec3{0.0, 2.0,  0.0}, vec.Vec3{0.0, 0.0, -2.0}, redMaterial)   // right wall

	// light, slightly below the ceiling so it doesn't overlap with it
	addQuad(vec.Vec3{-0.4, 1.999, -0.6}, vec.Vec3{0.0, 0.0, -0.8}, vec.Vec3{0.8, 0.0, 0.0}, lightMaterial)

	world.Add(&render_object.Sphere{ Center: vec.Vec3{ 0.5, 0.3, -1.3 }, Radius: 0.3, Material: mirrorMaterial })

	teapotMesh, _, err := data_loader.Load3DModel("models/teapot/teapot.obj")
	if err != nil {
		log.Println("Error loading teapot model:", err)
	}

	teapotPosition := vec.Vec3{-0.35, 0.0, -0.9}
	teapotRotation := vec.Vec3{ 0.0,  0.5,  0.0}
	teapotScale    := vec.Vec3{ 0.5,  0.5,  0.5}

	log.Println("Making BVH...")
	teapotModel := render_object.NewModel(teapotMesh, teapotMaterial, teapotPosition, teapotRotation, teapotScale)
	log.Println("Made BVH successfully!")

	world.Add(&teapotModel)

	return world
}
