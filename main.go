package main

import (
	"fmt"

	"github.com/gumper23/sqlrs"
)

func main() {
	fmt.Println("Hello")
	rs := sqlrs.New()
	fmt.Printf("%+v\n", rs)
}
