package main

import (
	"flag"
	"fmt"

	"github.com/amenzhinsky/flagenv"
)

func main() {
	stringFlag := flag.String("string", "foo", "string flag")
	boolFlag := flag.Bool("bool", false, "bool flag")
	intFlag := flag.Int("int", 666, "int flag")
	flagenv.Parse()

	fmt.Printf("-string = %q\n", *stringFlag)
	fmt.Printf("-bool   = %t\n", *boolFlag)
	fmt.Printf("-int    = %d\n", *intFlag)
}

type stringSliceFlag []string
