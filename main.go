package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Print(quadraticEquation(setRootsOfEquation()))
}

func setRootsOfEquation() (float64, float64, float64) {
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
	return a, b, c
}

func quadraticEquation(a float64, b float64, c float64) string {
	delta := getDelta(a, b, c)

	var x1 float64
	var x2 float64
	var result = "S = {} = Ø\n"

	if delta > 0 {
		x1 = bhaskara(a, b, delta, "sum")
		x2 = bhaskara(a, b, delta, "minus")
		result = fmt.Sprintf("S = {%v, %v}\n", x1, x2)
	}

	if delta == 0 {
		x1 = bhaskara(a, b, delta, "")
		result = fmt.Sprintf("S = {%v}\n", x1)
	}

	return result
}

func getDelta(a float64, b float64, c float64) float64 {
	// delta = b² - 4ac
	delta := math.Pow(b, 2) - 4*a*c
	fmt.Printf("\nΔ = %v\n", delta)
	return delta
}

func bhaskara(a float64, b float64, delta float64, signal string) float64 {
	// x = ((b*-1) +/- raiz(delta)) / (2 * a)
	if signal == "sum" {
		return ((b * -1) + math.Sqrt(delta)) / (2 * a)
	}

	if signal == "minus" {
		return ((b * -1) - math.Sqrt(delta)) / (2 * a)
	}

	// x = (b*-1) / (2 * a)
	return (b * -1) / (2 * a)
}
