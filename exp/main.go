package main

import (
	"fmt"

	"github.com/Users/patrickfurtak/desktop/go-gallery/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
