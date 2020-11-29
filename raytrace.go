package main

import (
    "math"
    "math/rand"
    "image"
    "image/color"
    "image/png"
    "io"
    "os"
)

const tolerance = 0.000000001

var world []sphere

func sq(x float64) float64 {
  return x*x
}

func minroot(a, b, c float64) float64 {
  if math.Abs(a) < tolerance {
    return b / -c
  }
  discrt := math.Sqrt(sq(b) - 4*a*c)
  return math.Min((-b + discrt)/(2*a), (-b - discrt)/(2*a))
}

func mag(v vector) float64 {
  return math.Sqrt(magsq(v))
}

func magsq(v vector) float64 {
  return sq(v.x) + sq(v.y) + sq(v.z)
}

type vector struct {
  x,y,z float64
}

var eye vector = vector{0,0,800}

type sphere struct {
  color color.Color
  radius float64
  center vector
}

func lambert(s sphere, intersection vector, v vector) float64 {
  return math.Max(0.0, math.Abs(dot(normal(s, intersection), v)))
}

func dot(a, b vector) float64 {
  return a.x*b.x + a.y*b.y + a.z*b.z
}

func normal(s sphere, intersection vector) vector {
  return unit(diff(s.center, intersection))
}

func scale(v vector, s float64) vector {
  return vector{v.x*s, v.y*s, v.z*s}
}

func defSphere(x, y, z, r float64, c color.Color) sphere {
  s := sphere{color: c, radius: r, center: vector{x: x, y: y, z: z}}
  world = append(world, s)
  return s
}

func scaleColor(c color.Color, scale float64) color.Color {
  r, g, b, a := c.RGBA()
  return color.RGBA64{uint16(float64(r) * scale), uint16(float64(g) * scale), uint16(float64(b) * scale), uint16(a)}
}

func unit(a vector) vector {
  d := mag(a)
  return vector{x: a.x/d, y: a.y/d, z: a.z/d}
}

func diff(from, to vector) vector {
  return vector{x: from.x - to.x, y: from.y - to.y, z: from.z - to.z}
}

func sendRay(from, dir vector) color.Color {
  if s, hit, ok := firstHit(from, dir); ok {
    return scaleColor(s.color, lambert(s, hit, dir))
  }
  return color.RGBA{0,0,0,0}
}

func firstHit(from, dir vector) (sphere, vector, bool) {
  var hit vector
  var dist float64 = -1.0
  var sp sphere
  found := false
  for _, s := range world {
    if h, ok := intersect(s, from, dir); ok {
      d := mag(diff(h, from))
      if dist < 0.0 || d < dist {
        sp, dist, hit, found = s, d, h, true
      }
    }
  }
  return sp, hit, found
}

func intersect(s sphere, pos, dir vector) (vector, bool) {
  n := minroot(magsq(dir),
      2*dot(diff(pos, s.center), dir),
      magsq(diff(pos, s.center))-sq(s.radius))
  if math.IsNaN(n) {
    return vector{}, false
  }
  return diff(pos, scale(dir, -n)), true
}

func colorAt(x, y int) color.Color {
  return sendRay(eye, unit(diff(vector{float64(x),float64(y),0}, eye)))
}

func tracer(out io.Writer) {
  const (
    width = 800
    height = 600
    res = 1
  )
  rect := image.Rect(0, 0, width, height)
  img := image.NewRGBA(rect)
  for i := 0; i < width; i++ {
    for j := 0; j < height; j++ {
      img.Set(i, j, colorAt(i-int(width/2/float64(res)),j-int(height/2/float64(res))))
    }
  }
  png.Encode(out, img)
}

func main() {
  defSphere(0, -300, -1200, 200, color.RGBA{0, 255, 0, 255})
  defSphere(-80,-150, -1200, 200, color.RGBA{0,0,255,255})
  defSphere(70,-100, -1200, 200, color.RGBA64{65535,0,0,65535})
  for i := -10; i < 11; i++ {
    for k := 2; k < 20; k++ {
      defSphere(float64(200*i), 700.0, float64(-400*k), 40 , color.RGBA{uint8(rand.Int31n(256)), uint8(rand.Int31n(256)), uint8(rand.Int31n(256)), 255})
    }
  }
  tracer(os.Stdout)
}
