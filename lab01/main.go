package main
import(
	"fmt"
	"math"
	_"reflect"
	)
//структура слогаемого
type term struct{
	coef float64
	pow int
}
//добавить слагаемое к полиному
func add(poly map[int]float64, term *term)map[int]float64{
	if _, ok := poly[term.pow]; ok{
		poly[term.pow] += term.coef
	}else{
		poly[term.pow] = term.coef
	}
	return poly
}
//умножить два слогаемых и записать в полином
func mult(poly map[int]float64, term1, term2 *term)map[int]float64{
	to_add := multterm(term1, term2)
	add(poly, to_add)
	return poly
}
//перемножает два слогаемых и возвращает результат
func multterm(term1, term2 *term)*term{
	res := term{term1.coef*term2.coef, term2.pow+term1.pow}
	return &res
}
//возведение многочлена в квадрат
func poly_pow(poly map[int]float64) map[int]float64{
	res := make(map[int]float64)
	for i, j := range poly{
		for k, z := range poly{
			//fmt.Println(i,j,k,z)
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
		//res[integr.pow] = integr.coef
	}
	//fmt.Println("res - ", res)
	return res, answer
}

func picar(x float64, n int)float64{
	u0 := 0.0
	answer := 0.0
	poly := make(map[int]float64)
	curr := make(map[int]float64)
	poly[2] = 1.0
	var res float64
	for i:=0;i<n;i++{
		curr = poly_pow(curr)
		curr[2] = 1.0
		//fmt.Println("-----", i, "-----")
		//fmt.Println("curr", curr)
		curr, res = integrate(curr, 0.0, x)
		answer = u0 + res
	}
	return answer
}
func f(x, y float64)float64{
	return x*x+y*y
}
func euler_implicit(xn float64, n int)float64{
	h := xn / float64(n)
	y:=0.0 
	x:=0.0
	for i:=0; i<=n;i++{
		y = y + h*f(x, y)
		x+=h
	}
	return y
}

func print_info(poly map[int]float64){
	for i, j := range poly{
		fmt.Println(i, j)
	}
}
func print_res(x float64, n, npicar int){
	fmt.Printf("%10s|%-32s|%10s|%10s|\n", "x", fmt.Sprintf("%32s","Пикар"), "Явн.", "Неявн.")
	for i:=0; i<66;i++{ fmt.Print("-")}
	fmt.Print("\n")
	fmt.Printf("%10s|%10s|%10s|%10s|%10s|%10s|\n", " ", "3-e", "4-e", "n-e", " ", " ")
	i:=0.0
	for ;i<=x;i+=0.01{
		fmt.Printf("%10.5f|%10.5f|%10.5f|%10.5f|%10.5f|%10.5f|\n", i, picar(i, 3), picar(i, 4), picar(i, npicar), euler_implicit(i,n), 0.1)
	}
	for i:=0; i<66;i++{ fmt.Print("-")}
}
func main(){
	print_res(2.1, 100, 6)
}