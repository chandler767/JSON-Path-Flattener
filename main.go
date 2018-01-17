package main

import (
	"os"
	"errors"
	"fmt"
)

func main() {
	if len(os.Args) > 1 {
		jsontoflatten := os.Args[1] // Expecting a path, URL, or string.
		if jsontoflatten != "" {
			return
	    }
	}
	fmt.Println(errors.New("A path, URL, or string is required."))
}
