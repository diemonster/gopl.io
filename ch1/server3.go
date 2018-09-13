// Server3 is an "echo" server that displays request parameters.
package main

import (
	// "fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!+handler
// handler echoes the HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	// for k, v := range r.Header {
	// 	fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	// }
	// fmt.Fprintf(w, "Host = %q\n", r.Host)
	// fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	// if err := r.ParseForm(); err != nil {
	// 	log.Print(err)
	// }
	// for k, v := range r.Form {
	// 	fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	// }
	if strings.Contains(r.URL.String(), "?cycles=") {
		c := strings.SplitAfter(r.URL.String(), "=")
		cycles, err := strconv.Atoi(c[len(c)-1])
		if err != nil {
			log.Print(err)
		}
		lissajous(cycles, w)
	}
}

func lissajous(cycles int, out io.Writer) {
	var palette = []color.Color{color.White, color.Pink}
	const (
		res        = 0.001 // angular resolution
		size       = 100   // image canvas covers [-size..+size]
		nframes    = 64    // number of animation frames
		delay      = 8     // delay between frames in 10ms units
		blackIndex = 1
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
