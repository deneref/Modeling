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
var lp = 12

var Tstart = 300
var I0 = 0.5
var U0 = 1500

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

func (spline *dataSplines) T0(I float64) float64 {
	return spline.ItoTo.At(I)
}

func (spline *dataSplines) m(I float64) float64 {
	return spline.Itom.At(I)
}

func (spline *dataSplines) sigma(T float64) float64 {
	return spline.TtoSigma.At(T)
}

func dUdt(I float64) float64 {
	return -I / float64(Ck)
}

func dIdT(splines *dataSplines, U, I float64) float64 {
	return (U - (Rk+Rp(splines, I))*I) / Lk
}

func Rp(splines *dataSplines, I float64) float64 {
	s := make([]float64, 0)
	n := 100.0
	for i := 0.0; i <= 1.0; i += 1.0 / n {
		s = append(s, i*splines.sigma(T(splines, i, I)))
	}

	return float64(lp) / (2 * 3.14 * math.Pow(Radius, 2) * integrate(s, n))
}

func integrate(s []float64, h float64) float64 {
	result := (s[0] + s[len(s)-1]) * 0.5
	for i := 1; i < len(s)-1; i++ {
		result += s[i]
	}

	return result * h
}

func T(splines *dataSplines, z, I float64) float64 {
	return splines.T0(I) + (float64(Tw)-splines.T0(I))*math.Pow(z, splines.m(I))
}

func runge_kutta4(xn, I, U float64, n int, splines dataSplines) float64 {
	h := xn / float64(n)

	for i := 0; i < n; i++ {
		k1 := dIdT(&splines, I, U)
		m1 := dUdt(I)

		k2 := dIdT(&splines, I+h/2, U+h*m1*0.5)
		m2 := dUdt(I + h*k1*0.5)

		k3 := dIdT(&splines, I+h*0.5, U+h*m2*0.5)
		m3 := dUdt(I + h*k2*0.5)

		k4 := dIdT(&splines, I+h*k3, U+h*m3)
		m4 := dUdt(I + h*k3)

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

/*
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

*/
