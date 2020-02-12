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

func picar(x float64, n int)map[int]float64{
	u0 := 0.0
	answer := make(map[int]float64, 0)
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
		answer[i] = u0 + res
	}
	return answer
}
func print_info(poly map[int]float64){
	for i, j := range poly{
		fmt.Println(i, j)
	}
}
func print_res(res map[int]float64, n int){
	for i:= 0; i<n;i++{
		fmt.Println(i+1, res[i])
	}
}
func main(){
	res := picar(1.5, 5)
	print_res(res, 5)
}