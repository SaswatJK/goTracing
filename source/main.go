package main

import "fmt"

type vec3 struct {
	a float32
	b float32
	c float32
}

type Camera struct {
	position  vec3
	direction vec3
}

func main() {
	fmt.Println("Hello world")
}
