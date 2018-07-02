package main

import (
	"fmt"
)

type Point struct{
	X, Y int
}

func main() {
	a := []int{1,2,3,4,5}
	fmt.Println(a)

	a = a[:3]
	fmt.Println(a)

	a = a[:5]
	fmt.Println(a)

	a = a[3:]
	fmt.Println(a)


	a = append(a, 1, 2)
	fmt.Println(a)
}
