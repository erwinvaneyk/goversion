//go:generate goversion generate -o version.gen.go
package main

import (
	"fmt"

	"github.com/erwinvaneyk/goversion"
)

func main() {
	fmt.Println(goversion.Get())
}
