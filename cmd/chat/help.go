// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"flag"
	"fmt"
)

func Help() {
	fmt.Println("Copyright:", "2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>")
	fmt.Println("Version:", Version)
	fmt.Printf("Usage: %s [options]\n", Name)
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println()
}
