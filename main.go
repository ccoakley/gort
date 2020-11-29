package main

import (
  "image/color"
  "math/rand"
  "os"
  "raytrace"
)

func main() {
  raytrace.DefSphere(0, -300, -1200, 200, color.RGBA{0, 255, 0, 255})
  raytrace.DefSphere(-80,-150, -1200, 200, color.RGBA{0,0,255,255})
  raytrace.DefSphere(70,-100, -1200, 200, color.RGBA64{65535,0,0,65535})
  for i := -10; i < 11; i++ {
    for k := 2; k < 20; k++ {
      raytrace.DefSphere(float64(200*i), 700.0, float64(-400*k), 40 , color.RGBA{uint8(rand.Int31n(256)), uint8(rand.Int31n(256)), uint8(rand.Int31n(256)), 255})
    }
  }
  raytrace.Tracer(os.Stdout)
}
