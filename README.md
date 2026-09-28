# grt

Go Ray Tracer (actually a path tracer hehe).
Runs on CPU, but I plan on adding WebGPU compute shader calculation soon.

![teapot lit by three emitters](.github/pics/teapot_misty_lit.png)
![textured boat lit by three emitters](.github/pics/boat_textured.png)
![cornell box with teapot and mirror sphere](.github/pics/cornell_box.png)

Monte Carlo path tracing, BVH-accelerated triangle meshes, `.obj` loading.
Writes a binary PPM.

```sh
go build -o grt . && ./grt
```
