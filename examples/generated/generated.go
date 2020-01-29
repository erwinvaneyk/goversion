//go:generate goversion fields -o version.gen.go
package main

import (
	"fmt"

	goversion "github.com/erwinvaneyk/go-version"
)

func main() {
	fmt.Println(goversion.Get())
}
