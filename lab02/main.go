package main

import (
	"fmt"
	"math"
	_ "reflect"
	"sort"

	"gospline"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var Ck = 0.000268
var Lk = 0.000187
var Rk = 0.25
var Radius = 0.0035
var Uco = 1400
var Tw = 2000

//term structure
type term struct {
	coef float64
	pow  int
}

type dataSplines struct {
	ItoTo    gospline.Spline
	Itom     gospline.Spline
	TtoSigma gospline.Spline
}

//add term to polinomial
func add(poly map[int]float64, term *term) map[int]float64 {
	if _, ok := poly[term.pow]; ok {
		poly[term.pow] += term.coef
	} else {
		poly[term.pow] = term.coef
	}
	return poly
}

//multiply two terms and write to polinomial
func mult(poly map[int]float64, term1, term2 *term) map[int]float64 {
	to_add := multterm(term1, term2)
	add(poly, to_add)
	return poly
}

//multiply two terms and return result
func multterm(term1, term2 *term) *term {
	res := term{term1.coef * term2.coef, term2.pow + term1.pow}
	return &res
}

//squaring a polinomial
func poly_pow(poly map[int]float64) map[int]float64 {
	res := make(map[int]float64)
	for i, j := range poly {
		for k, z := range poly {
			mult(res, &term{j, i}, &term{z, k})
		}
	}
	return res
}

func integrate(poly map[int]float64, x0, x float64) (map[int]float64, float64) {
	var answer float64
	res := make(map[int]float64)
	for i, j := range poly {
		integr := term{j, i + 1}
		integr.coef *= 1.0 / float64(integr.pow)
		res[integr.pow] = integr.coef
		answer += integr.coef*math.Pow(x, float64(integr.pow)) - integr.coef*math.Pow(x0, float64(integr.pow))
	}
	return res, answer
}

//given function
func f(x, y float64) float64 {
	return x*x + y*y
}

func dUdt(I float64) float64 {
	return -I / float64(Ck)
}

func dIdT(U, I float64) float64 {
	return (U - (Rk+R(I))*I) / Lk
}

func R(I float64) float64 {
	return 1
}

func T(z, I, m, T0 float64) float64 {
	return T0 + (float64(Tw)-T0)*math.Pow(z, m)
}

func D(a, b, c float64) float64 {
	return b*b - 4*a*c
}

func print_info(poly map[int]float64) {
	for i, j := range poly {
		fmt.Println(i, j)
	}
}

func runge_kutta4(xn, I0, Uo float64, n int, splines dataSplines) float64 {
	h := xn / float64(n)
	y := 0.0
	x := 0.0

	polinom := make(map[int]float64)

	for i := 0; i < n; i++ {
		k1 := dIdT(x, y)
		k2 := dIdT(x+h/2, y+h/2*k1)
		k3 := dIdT(x+h/2, y+h/2*k2)
		k4 := dIdT(x+h, y+h*k3)

		//fmt.Printf("%f %f %f %f\n", k1, k2, k3, k4)

		y += h / 6 * (k1 + 2*k2 + 2*k3 + k4)
		x += h
	}
	_ = x
	return y
}

func main() {
	//x, n for euler, n for picard
	//print_res(2.2, 1000000, 3)
	n := 1000

	splines := new(dataSplines)

	splines.ItoTo = gospline.NewCubicSpline([]float64{0.5, 1, 5, 10, 50, 200, 400, 800, 1200}, []float64{6730, 6790, 7150, 7270, 8010, 9185, 10010, 11140, 12010})
	splines.Itom = gospline.NewCubicSpline([]float64{0.5, 1, 5, 10, 50, 200, 400, 800, 1200}, []float64{0.5, 0.55, 1.7, 3, 11, 32, 40, 41, 39})
	splines.TtoSigma = gospline.NewCubicSpline([]float64{4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000, 12000, 13000, 14000},
		[]float64{0.031, 0.27, 2.05, 6.06, 12.0, 19.9, 29.6, 41.1, 54.1, 67.7, 81.5})

	makeGraphFileFromSpline(splines.ItoTo, "ItoT0.png", "Зависимость T0 от I", "I", "T0", 0, 1200)
	makeGraphFileFromSpline(splines.Itom, "Itom.png", "Зависимость m от I", "I", "m", 0, 1200)
	makeGraphFileFromSpline(splines.TtoSigma, "ItoSigma.png", "Зависимость σ от T", "T", "σ", 4000, 10000)

	rungeKutta4Points := make(map[float64]float64)

	for i := 0.0; i < 2.0; i += 0.1 {
		rungeKutta4Points[i] = runge_kutta4(i, n, *splines)
	}

	makeGraphFileFromMap(rungeKutta4Points, "f.png", "Test x^2+y^2", "i", "i")

	//fmt.Print(ItoTo.At(4))
}

func makeGraphPointsFromSpline(spline gospline.Spline, n, start, end int) plotter.XYs {
	h := (end - start) / n
	startValue := start
	pts := make(plotter.XYs, n)
	for i := range pts {
		pts[i].X = float64(startValue)
		pts[i].Y = spline.At(float64(startValue))
		startValue += h
	}
	return pts
}

func makeGraphFileFromSpline(spline gospline.Spline, filename, title, Xlegend, Ylegend string, start, end int) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = Xlegend
	p.Y.Label.Text = Ylegend

	err := plotutil.AddLinePoints(p,
		"_", makeGraphPointsFromSpline(spline, 10, start, end))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func makeGraphPointsFromMap(spline map[float64]float64) plotter.XYs {
	pts := make(plotter.XYs, len(spline))

	keys := make([]float64, 0)
	for k, _ := range spline {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	i := 0
	for _, k := range keys {
		fmt.Println(k, spline[k])
		pts[i].X = k
		pts[i].Y = spline[k]
		i++
	}
	return pts
}

func makeGraphFileFromMap(spline map[float64]float64, filename, title, Xlegend, Ylegend string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = Xlegend
	p.Y.Label.Text = Ylegend

	err := plotutil.AddLinePoints(p,
		"_", makeGraphPointsFromMap(spline))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, filename); err != nil {
		panic(err)
	}
}
