package main

import (
	"fmt"
)

type next func(a, b string)

func (n next) serveHTTP(a, b string) {
	n(a, b)
}

type handler interface {
	serveHTTP(a, b string)
}

type con func(handler) handler
type chain struct {
	cons []con
}

func _new(conss ...con) chain {
	return chain{append([]con(nil), conss...)}
}

func (c chain) then(h handler) handler {
	for i := range c.cons {
		h = c.cons[len(c.cons)-1-i](h)
	}
	return h
}

func m1(h handler) handler {
	return next(func(a, b string) {
		fmt.Println("be")
		h.serveHTTP(a, b)
		fmt.Println("af")
	})
}

func m2(h handler) handler {
	return next(func(a, b string) {
		fmt.Println("2 be")
		h.serveHTTP(a, b)
		fmt.Println("2 af")
	})
}

type name struct {
	info string
}

func (n name) serveHTTP(a, b string) {
	fmt.Println(n.info, a, b)
}

func main() {
	test := name{
		info: "test",
	}
	_new(m1, m2).then(test).serveHTTP("a", "b")
}
