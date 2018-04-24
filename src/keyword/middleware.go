package main

import (
	"fmt"
)

type middle func(int, int)

type Func func(middle) middle

func funItem1(f middle) middle {
	return middle(func(a, b int) {
		fmt.Println("funItem 1")
		f(a, b)
		fmt.Println("funItem 2")
	})
}

func funItem2(f middle) middle {
	return middle(func(a, b int) {
		fmt.Println("funItem 10")
		f(a, b)
		fmt.Println("funItem 20")
	})
}

func chain(m middle,f ...Func) middle {
	if len(f) == 0 {
		return m
	}
	return f[0](chain(m,f[1:]...))
}

func f1(a, b int) {
	fmt.Println("invole f1")
}

func main() {


	m := chain(f1,funItem1,funItem2)
	m(1,2)
}
