package main
import(
	"fmt"
	"math"
	_"reflect"
	)
//term structure
type term struct{
	coef float64
	pow int
}
//add term to polinomial
func add(poly map[int]float64, term *term)map[int]float64{
	if _, ok := poly[term.pow]; ok{
		poly[term.pow] += term.coef
	}else{
		poly[term.pow] = term.coef
	}
	return poly
}
//multiply two terms and write to polinomial
func mult(poly map[int]float64, term1, term2 *term)map[int]float64{
	to_add := multterm(term1, term2)
	add(poly, to_add)
	return poly
}
//multiply two terms and return result
func multterm(term1, term2 *term)*term{
	res := term{term1.coef*term2.coef, term2.pow+term1.pow}
	return &res
}
//squaring a polinomial
func poly_pow(poly map[int]float64) map[int]float64{
	res := make(map[int]float64)
	for i, j := range poly{
		for k, z := range poly{
			mult(res, &term{j, i}, &term{z,k})
		}
	}
	return res
}

func integrate(poly map[int]float64, x0, x float64) (map[int]float64, float64){
	var answer float64
	res := make(map[int]float64)
	for i, j := range poly{
		integr := term{j, i+1}
		integr.coef *= 1.0/float64(integr.pow)
		res[integr.pow] = integr.coef
		answer += integr.coef*math.Pow(x, float64(integr.pow)) - integr.coef*math.Pow(x0, float64(integr.pow))
	}
	return res, answer
}

func picard(x float64, n int)float64{
	u0 := 0.0
	answer := 0.0
	poly := make(map[int]float64)
	curr := make(map[int]float64)
	poly[2] = 1.0
	var res float64
	for i:=0;i<n;i++{
		curr = poly_pow(curr)
		curr = add(curr, &term{1.0,2})
		curr, res = integrate(curr, 0.0, x)
		answer = u0 + res
	}
	return answer
}

//given function
func f(x, y float64)float64{
	return x*x+y*y
}

//Runge–Kutta method
func runge_kutta(xn float64, n int)float64{
	h := xn / float64(n)
	alpha := 0.5
	y := 0.0
	x := 0.0

	for i:=0; i<n;i++{
		y += h * ( (1- alpha) * f(x,y) + alpha * f(x+h/(2*alpha), y+(h/2*alpha)*f(x,y)) )
		x+=h
	}
	_ = x
	return y
}


//explicit euler method
func euler_explicit(xn float64, n int)float64{
	h := xn / float64(n)
	y:=0.0 
	x:=0.0
	for i:=0; i<n;i++{
		y += h*f(x, y)
		x+=h
	}
	_ = x
	return y
}

//implicit(backward) euler method
func euler_implicit(xn float64, n int)float64{
	//yn+1^2 - 1/h*yn+1 + 1/h*yn + xn+1^2
	//n *= 10
	h := xn / float64(n)
	y:=0.0 
	x:=0.0
	var a, b, c, dis, y1, y2 float64
	for i:=0; i<n;i++{
		a = 1; b = -1.0/h; c = 1.0/h*y+(x+h)*(x+h)
		//a = h; b = -1.0; c = y+math.Pow(x+h,2)*h
		dis = D(a, b, c)
		if dis>=0{
			y1 = (-b - math.Sqrt(dis))/2/a
			y2 = (-b + math.Sqrt(dis))/2/a
			_ = y2
		}
		//fmt.Println(x, y, y1, y2)
		y = y1
		x+=h
	}
	return y
}

func D(a, b, c float64)float64{
	return b*b - 4*a*c
}

func print_info(poly map[int]float64){
	for i, j := range poly{
		fmt.Println(i, j)
	}
}
func print_res(x float64, n, npicar int){
	fmt.Printf("%10s|%-32s|%10s|%10s|%10s|\n", "x", fmt.Sprintf("%32s","Picard's"), "Explicit", "Implicit", "Runge-Kutta")
	for i:=0; i<66;i++{ fmt.Print("-")}
	fmt.Print("\n")
	fmt.Printf("%10s|%d %-8s|%d %-8s|%d %-8s|%10s|%10s|%10s|\n", " ",npicar-2, "-e", npicar-1,"-e", npicar,"-e", " ", " ", " ")
	i:=0.0
	for ;i<=x;i+=0.01{
		fmt.Printf("%10.5f|%10.5f|%10.5f|%10.5f|%10.5f|%10.5f|%10.5f|\n", i, picard(i, npicar-2), picard(i, npicar-1),
		 picard(i, npicar), euler_explicit(i,n), euler_implicit(i,n), runge_kutta(i, n))
	}
	for i:=0; i<66;i++{ fmt.Print("-")}
}
func main(){
	//x, n for euler, n for picard
	print_res(2.2, 1000000, 3)
}

