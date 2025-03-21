package main

import "fmt"

type Todo struct {
	id        int
	libelle   string
}

var todo = Todo{id: 1, libelle: "coucou"}
var todo2 = &todo
var todo3 = Todo{}

func main() {
	eh := "coucou"
	fmt.Printf(eh)
	display()
	espace()
	todo.id = 2
	espace()
	display()
	espace()
}

func espace(){
	fmt.Println()
}

func display() {
	fmt.Printf("todo : %v", todo)
	fmt.Printf("todo 2 : %v", todo2)
	fmt.Printf("todo 3 : %v", todo3)
}
