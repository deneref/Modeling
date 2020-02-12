package main
import(
	"fmt"
	"math"
	"reflect"
	)

type term struct{
	coef float64
	pow float64
}
//возведение слогаемого в квадрат
func pow2(cur_term *term)*term{
	tmp := new(term)
	tmp.coef = math.Pow(cur_term.coef,2)
	tmp.pow = cur_term.pow*2
	return tmp
}
//двойное произведение двух множителей
func mult(t1, t2 *term) *term{
	res := new(term)
	res.coef = t1.coef*t2.coef
	res.pow = t1.pow+t2.pow
	return res
}
//возведение многочлена в квадрат
func poly_pow(poly[]*term) []*term{
	res := make([]*term, 0)
	print_info(poly)
	for _, i := range poly{
		for _, j := range poly{
			if reflect.DeepEqual(i,j){
				res = append(res, mult(i,j))
			}else{
				res = append(res, pow2(i))
			}
			fmt.Println(i,j,res)
		}
	}
	return res
}

func integrate(poly []*term, x0, x float64) ([]*term, float64){
	var answer float64
	for i, j := range poly{
		j := term{j.coef, j.pow+1}
		j.coef *= 1.0/j.pow
		answer += j.coef*math.Pow(x, j.pow) - j.coef*math.Pow(x0, j.pow)
		poly[i] = &j
	}
	return poly, answer
}

func picar(x float64)float64{
	//u0 := 0
	poly := make([]*term, 0)
	poly = append(poly, &term{1, 2})
	poly = append(poly, &term{1.0/3.0, 3})
	poly1 := poly_pow(poly)
	print_info(poly1)
	return 0
}
func print_info(p []*term){
	for i, j := range p{
		fmt.Println(i, ":", j.coef, j.pow)
	}
}
func main(){
	_ = picar(1.5)
}