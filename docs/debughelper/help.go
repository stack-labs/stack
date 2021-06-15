package main

import "fmt"

const (
	TypeEnd      = 1
	TypeRegistry = iota * 2
	TypeRelaying
)

func main() {
	fmt.Println(TypeEnd)
	fmt.Println(TypeRegistry)
	fmt.Println(TypeRelaying)
}
