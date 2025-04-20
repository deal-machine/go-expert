package main

import (
	"fmt"
	"math"
)

func main() {
	var (
		a float64
		b float64
		c float64
	)
	fmt.Print("Digite o primeiro valor: ")
	fmt.Scanln(&a)
	fmt.Print("Digite o segundo valor: ")
	fmt.Scanln(&b)
	fmt.Print("Digite o terceiro valor: ")
	fmt.Scanln(&c)

	fmt.Printf("a = %v, b = %v, c = %v\n", a, b, c)

	fmt.Print(quadraticEquation(a, b, c))
}

func quadraticEquation(a float64, b float64, c float64) string {
	// delta = b² - 4ac
	delta := math.Pow(b, 2) - 4*a*c
	fmt.Printf("\nΔ = %v\n", delta)

	var x1 float64
	var x2 float64
	var result = "S = {} = Ø\n"

	if delta > 0 {
		// x = ((b*-1) +/- raiz(delta)) / (2 * a)
		x1 = ((b * -1) + math.Sqrt(delta)) / (2 * a)
		x2 = ((b * -1) - math.Sqrt(delta)) / (2 * a)
		result = fmt.Sprintf("S = {%v, %v}\n", x1, x2)
	}

	if delta == 0 {
		// x = (b*-1) / (2 * a)
		x1 = (b * -1) / (2 * a)
		result = fmt.Sprintf("S = {%v}\n", x1)
	}

	return result
}
