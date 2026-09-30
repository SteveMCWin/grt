package camera

import (
	"bufio"
	"fmt"
	"grt/ray"
	"grt/render_object"
	"grt/util"
	"grt/vec"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Camera struct {
	AspectRatio  float64
	FocalLength  float64
	Vfov         float64
	DefocusAngle float64
	FocusDist    float64

	ImageWidth int
	SqrtSPP    float64
	MaxDepth   int

	Position vec.Vec3
	LookAt   vec.Vec3
	Up       vec.Vec3

	imageHeight    int
	recipSqrtSPP   float64
	viewportWidth  float64
	viewportHeight float64
	samplesScale   float64
	startPixelLoc  vec.Vec3
	pixelDeltaU    vec.Vec3
	pixelDeltaV    vec.Vec3
	viewportU      vec.Vec3
	viewportV      vec.Vec3
	defocusU       vec.Vec3
	defocusV       vec.Vec3
	u              vec.Vec3
	v              vec.Vec3
	w              vec.Vec3

	Background vec.Color
}

type CamDescriptor struct {
	Position        vec.Vec3
	LookAt          vec.Vec3
	Up              vec.Vec3
	Background      vec.Color
	AspectRatio     float64
	Vfov            float64
	DefocusAngle    float64
	FocusDist       float64
	ImageWidth      int
	SamplesPerPixel int
	MaxDepth        int
}

func NewCam(desc CamDescriptor) Camera {
	cam := Camera{
		AspectRatio:  desc.AspectRatio,
		Vfov:         desc.Vfov,
		DefocusAngle: desc.DefocusAngle,
		FocusDist:    desc.FocusDist,
		ImageWidth:   desc.ImageWidth,
		SqrtSPP:      math.Sqrt(float64(max(1, desc.SamplesPerPixel))),
		MaxDepth:     max(2, desc.MaxDepth),
		Position:     desc.Position,
		LookAt:       desc.LookAt,
		Up:           desc.Up,
		Background:   desc.Background,
	}

	cam.imageHeight = max(1, int(float64(desc.ImageWidth)/cam.AspectRatio))

	cam.samplesScale = 1.0 / float64(cam.SqrtSPP*cam.SqrtSPP)
	cam.recipSqrtSPP = 1.0 / cam.SqrtSPP

	// focal_length := vec.Subtract(cam.Center, cam.LookAt).Len()
	theta := util.DegToRad(cam.Vfov)
	h := math.Tan(theta / 2.0)
	cam.viewportHeight = 2.0 * h * cam.FocusDist
	cam.viewportWidth = cam.viewportHeight * (float64(cam.ImageWidth) / float64(cam.imageHeight))

	cam.w = vec.Subtract(cam.Position, cam.LookAt).Normalized()
	cam.u = vec.Cross(cam.Up, cam.w).Normalized()
	cam.v = vec.Cross(cam.w, cam.u)

	cam.viewportU = cam.u.Scaled(cam.viewportWidth)
	cam.viewportV = cam.v.Neg().Scaled(cam.viewportHeight)

	cam.pixelDeltaU = cam.viewportU.Scaled(1.0 / float64(cam.ImageWidth))
	cam.pixelDeltaV = cam.viewportV.Scaled(1.0 / float64(cam.imageHeight))

	viewport_upper_left := vec.Subtract(cam.Position, cam.w.Scaled(cam.FocusDist), cam.viewportU.Scaled(0.5), cam.viewportV.Scaled(0.5))

	cam.startPixelLoc = vec.Add(viewport_upper_left, vec.Add(cam.pixelDeltaU, cam.pixelDeltaV).Scaled(0.5))

	defocus_radius := cam.FocusDist * math.Tan(util.DegToRad(cam.DefocusAngle/2.0))
	cam.defocusU = cam.u.Scaled(defocus_radius)
	cam.defocusV = cam.v.Scaled(defocus_radius)

	return cam
}

func (cam Camera) RenderScene(stuff_to_render render_object.Hittable, w ...io.Writer) {
	totalPixels := cam.ImageWidth * cam.imageHeight
	buffer := make([]vec.Vec3, totalPixels)

	var completedPixels atomic.Int64

	progressDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				done := completedPixels.Load()
				fmt.Fprintf(os.Stderr, "\rRendered %d/%d pixels (%.1f%%)", done, totalPixels, 100*float64(done)/float64(totalPixels))
			case <-progressDone:
				return
			}
		}
	}()

	maxPatchSize := 16
	numWorkers := runtime.NumCPU()

	horizontalPatches := int(math.Ceil(float64(cam.ImageWidth)/float64(maxPatchSize)))
	verticalPatches := int(math.Ceil(float64(cam.imageHeight)/float64(maxPatchSize)))
	numPatches :=  horizontalPatches * verticalPatches
	startPixels := make([]vec.Vec3, 0, numPatches)

	for j := range verticalPatches {
		for i := range horizontalPatches {
			startPixels = append(startPixels, vec.Vec3{float64(i*maxPatchSize), float64(j*maxPatchSize), 0})
		}
	}

	patchQueue := make(chan vec.Vec3, numWorkers)

	go func() {
		for _, vec := range startPixels {
			patchQueue <- vec
		}
		close(patchQueue)
	}()

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pixelPos := range patchQueue {
				patchWidth := maxPatchSize - max(0, int(pixelPos.X()) + maxPatchSize - cam.ImageWidth)
				patchHeight := maxPatchSize - max(0, int(pixelPos.Y()) + maxPatchSize - cam.imageHeight)


				for j := range patchHeight {
					for i := range patchWidth {
						pixel_x := int(pixelPos.X()) + i
						pixel_y := int(pixelPos.Y()) + j
						col := vec.Vec3{0, 0, 0}
						for s_j := 0; float64(s_j) < cam.SqrtSPP; s_j++ {
							for s_i := 0; float64(s_i) < cam.SqrtSPP; s_i++ {
								r := cam.getOffsetRay(pixel_x, pixel_y, s_i, s_j)
								col = vec.Add(col, vec.Vec3(cam.RayColor(r, cam.MaxDepth, stuff_to_render)))
							}
						}
						buffer[pixel_y*cam.ImageWidth+pixel_x] = col.Scaled(cam.samplesScale)
						completedPixels.Add(1)
					}
				}
			}
		}()
	}

	wg.Wait()
	close(progressDone)
	fmt.Fprintf(os.Stderr, "\rRendered %d/%d pixels (100.0%%)\n", totalPixels, totalPixels)

	if len(w) < 1 {
		w = make([]io.Writer, 1)
		w[0] = os.Stdout
	}

	fmt.Fprintf(w[0], "P6\n%d %d\n255\n", cam.ImageWidth, cam.imageHeight)
	writer := bufio.NewWriterSize(w[0], 1<<20)
	defer writer.Flush()

	for _, col := range buffer {
		r, g, b := vec.Color(col).Components()
		writer.Write([]byte{r, g, b})
	}
}

func (cam Camera) sampleSquareStratified(s_i, s_j int) vec.Vec3 {
	px := ((float64(s_i) + rand.Float64()) * cam.recipSqrtSPP) - 0.5
	py := ((float64(s_j) + rand.Float64()) * cam.recipSqrtSPP) - 0.5

	return vec.Vec3{px, py, 0.0}
}

func (cam Camera) defocusDiskSample() vec.Vec3 {
	p := vec.RandInUnitDisk()
	return vec.Add(cam.Position, cam.defocusU.Scaled(p.X()), cam.defocusV.Scaled(p.Y()))
}

func (cam Camera) getOffsetRay(i, j, s_i, s_j int) ray.Ray {
	offset := cam.sampleSquareStratified(s_i, s_j)
	pixel_sample := vec.Add(cam.startPixelLoc, cam.pixelDeltaU.Scaled(float64(i)+offset.X()), cam.pixelDeltaV.Scaled(float64(j)+offset.Y()))

	var ray_origin vec.Vec3
	if cam.DefocusAngle <= 0 {
		ray_origin = cam.Position
	} else {
		ray_origin = cam.defocusDiskSample()
	}

	ray_direction := vec.Subtract(pixel_sample, ray_origin)
	return ray.NewRay(ray_origin, ray_direction)
}

func (cam Camera) RayColor(r ray.Ray, depth int, shape render_object.Hittable) vec.Color {
	if depth <= 0 {
		return vec.Color{0.0, 0.0, 0.0}
	}

	var rec render_object.HitRecord
	if !shape.Hit(&r, util.Interval{Min: 0.001, Max: util.PosInf}, &rec) {
		return cam.Background
	}

	if rec.Material == nil {
		panic("No material!!!")
	}

	var attenuation vec.Color
	var scattered ray.Ray

	if !rec.Material.Scatter(&r, &rec, &attenuation, &scattered) {
		return rec.Emitted
	}

	colorFromScatter := vec.Color(vec.Mul(vec.Vec3(attenuation), vec.Vec3(cam.RayColor(scattered, depth-1, shape))))
	return vec.Color{rec.Emitted[0] + colorFromScatter[0], rec.Emitted[1] + colorFromScatter[1], rec.Emitted[2] + colorFromScatter[2]}
}
