package camera

import (
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"sync"
	"grt/ray"
	"grt/render_object"
	"grt/util"
	"grt/vec"
)

type Camera struct {
	AspectRatio  float64
	FocalLength  float64
	Vfov         float64
	DefocusAngle float64
	FocusDist    float64

	ImageWidth      int
	SamplesPerPixel int
	MaxDepth        int

	Center vec.Vec3
	LookAt vec.Vec3
	Up     vec.Vec3

	imageHeight    int
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
}

func NewCam(aspect_ratio, vfov, defocus_angle, focus_dist float64, image_width, samples_per_pixel, max_depth int, cam_center, look_at, up vec.Vec3) Camera {
	cam := Camera{
		AspectRatio:     aspect_ratio,
		Vfov:            vfov,
		DefocusAngle:    defocus_angle,
		FocusDist:       focus_dist,
		ImageWidth:      image_width,
		SamplesPerPixel: max(1, samples_per_pixel),
		MaxDepth:        max_depth,
		Center:          cam_center,
		LookAt:          look_at,
		Up:              up,
	}

	cam.imageHeight = max(1, int(float64(image_width)/cam.AspectRatio))

	cam.samplesScale = 1.0 / float64(samples_per_pixel)

	// focal_length := vec.Subtract(cam.Center, cam.LookAt).Len()
	theta := util.DegToRad(cam.Vfov)
	h := math.Tan(theta / 2.0)
	cam.viewportHeight = 2.0 * h * cam.FocusDist
	cam.viewportWidth = cam.viewportHeight * (float64(cam.ImageWidth) / float64(cam.imageHeight))

	cam.w = vec.Subtract(cam.Center, cam.LookAt).Normalized()
	cam.u = vec.Cross(cam.Up, cam.w).Normalized()
	cam.v = vec.Cross(cam.w, cam.u)

	cam.viewportU = cam.u.Scaled(cam.viewportWidth)
	cam.viewportV = cam.v.Neg().Scaled(cam.viewportHeight)

	cam.pixelDeltaU = cam.viewportU.Scaled(1.0 / float64(cam.ImageWidth))
	cam.pixelDeltaV = cam.viewportV.Scaled(1.0 / float64(cam.imageHeight))

	viewport_upper_left := vec.Subtract(cam.Center, cam.w.Scaled(cam.FocusDist), cam.viewportU.Scaled(0.5), cam.viewportV.Scaled(0.5))

	cam.startPixelLoc = vec.Add(viewport_upper_left, vec.Add(cam.pixelDeltaU, cam.pixelDeltaV).Scaled(0.5))

	defocus_radius := cam.FocusDist * math.Tan(util.DegToRad(cam.DefocusAngle/2.0))
	cam.defocusU = cam.u.Scaled(defocus_radius)
	cam.defocusV = cam.v.Scaled(defocus_radius)

	return cam
}

func (cam Camera) Render(stuff_to_render render_object.Hittable, w ...io.Writer) {
	buffer := make([]vec.Vec3, cam.ImageWidth*cam.imageHeight)
	var wg sync.WaitGroup

	for j := range cam.imageHeight {
		for i := range cam.ImageWidth {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				col := vec.Vec3{0.0, 0.0, 0.0}
				for range cam.SamplesPerPixel {
					r := cam.getOffsetRay(i, j)
					col = vec.Add(col, vec.Vec3(cam.RayColor(r, cam.MaxDepth, stuff_to_render)))
				}
				buffer[j*cam.ImageWidth+i] = col.Scaled(cam.samplesScale)
			}(i, j)
		}
	}

	wg.Wait()

	if len(w) < 1 {
		w[0] = os.Stdout
	}

	fmt.Fprintf(w[0], "P3\n%d %d\n255\n", cam.ImageWidth, cam.imageHeight)
	for _, col := range buffer {
		fmt.Fprint(w[0], vec.Color(col))
	}
}

func (cam Camera) defocusDiskSample() vec.Vec3 {
	p := vec.RandInUnitDisk()
	return vec.Add(cam.Center, cam.defocusU.Scaled(p.X()), cam.defocusV.Scaled(p.Y()))
}

func (cam Camera) getOffsetRay(i, j int) ray.Ray {
	offset := vec.Vec3{rand.Float64() - 0.5, rand.Float64() - 0.5, 0.0}
	pixel_sample := vec.Add(cam.startPixelLoc, cam.pixelDeltaU.Scaled(float64(i)+offset.X()), cam.pixelDeltaV.Scaled(float64(j)+offset.Y()))

	var ray_origin vec.Vec3
	if cam.DefocusAngle <= 0 {
		ray_origin = cam.Center
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
	if shape.Hit(&r, util.Interval{Min: 0.001, Max: util.PosInf}, &rec) {
		var scattered ray.Ray
		var attenuation vec.Color

		if rec.Material.Scatter(&r, &rec, &attenuation, &scattered) {
			return vec.Color(vec.Mul(vec.Vec3(attenuation), vec.Vec3(cam.RayColor(scattered, depth-1, shape))))
		}

		return vec.Color{0.0, 0.0, 0.0}
	}

	a := 0.5 * (r.Direction().Y() + 1.0)
	res := vec.Add(vec.Vec3{1.0, 1.0, 1.0}.Scaled(1.0-a), vec.Vec3{0.5, 0.7, 1.0}.Scaled(a))
	return vec.Color(res)
}
